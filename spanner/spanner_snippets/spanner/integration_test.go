// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/grpc/codes"
)

type sampleFunc func(w io.Writer, dbName string) error
type instanceSampleFunc func(w io.Writer, projectID, instanceID string) error
type backupSampleFunc func(w io.Writer, dbName, backupID string) error

var (
	validInstancePattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)$")
)

func initTest(t *testing.T, id string) (dbName string, cleanup func()) {
	instance := getInstance(t)
	dbID := validLength(fmt.Sprintf("smpl-%s", id), t)
	dbName = fmt.Sprintf("%s/databases/%s", instance, dbID)

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		t.Fatalf("failed to create DB admin client: %v", err)
	}

	// Check for database existance prior to test start and delete, as resources
	// may not have been cleaned up from previous invocations.
	if db, err := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbName}); err == nil {
		t.Logf("database %s exists in state %s. delete result: %v", db.GetName(), db.GetState().String(),
			adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName}))
	}
	cleanup = func() {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			err := adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName})
			if err != nil {
				r.Errorf("DropDatabase(%q): %v", dbName, err)
			}
		})
		adminClient.Close()
	}
	return
}

func initBackupTest(t *testing.T, id, dbName string) (restoreDBName, backupID, cancelledBackupID string, cleanup func()) {
	instance := getInstance(t)
	restoreDatabaseID := validLength(fmt.Sprintf("restore-%s", id), t)
	restoreDBName = fmt.Sprintf("%s/databases/%s", instance, restoreDatabaseID)
	backupID = validLength(fmt.Sprintf("backup-%s", id), t)
	cancelledBackupID = validLength(fmt.Sprintf("cancel-%s", id), t)

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		t.Fatalf("failed to create admin client: %v", err)
	}
	if db, err := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: restoreDBName}); err == nil {
		t.Logf("database %s exists in state %s. delete result: %v", db.GetName(), db.GetState().String(),
			adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: restoreDBName}))
	}

	// Check for any backups that were created from that database and delete those as well
	iter := adminClient.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instance,
		Filter: "database:" + dbName,
	})
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Errorf("Failed to list backups for database %s: %v", dbName, err)
		}
		t.Logf("backup %s exists. delete result: %v", resp.Name,
			adminClient.DeleteBackup(ctx, &adminpb.DeleteBackupRequest{Name: resp.Name}))
	}

	cleanup = func() {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			err := adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: restoreDBName})
			if err != nil {
				r.Errorf("DropDatabase(%q): %v", restoreDBName, err)
			}
		})
	}
	return
}

func TestCreateInstance(t *testing.T) {
	_ = testutil.SystemTest(t)

	projectID, _, err := parseInstanceName(getInstance(t))
	if err != nil {
		t.Fatalf("failed to parse instance name: %v", err)
	}

	instanceID := fmt.Sprintf("go-sample-test-%s", uuid.New().String()[:8])
	out := runInstanceSample(t, createInstance, projectID, instanceID, "failed to create an instance")
	if err := cleanupInstance(projectID, instanceID); err != nil {
		t.Logf("cleanupInstance error: %s", err)
	}
	assertContains(t, out, fmt.Sprintf("Created instance [%s]", instanceID))
}

func TestSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	var out string
	mustRunSample(t, createDatabase, dbName, "failed to create a database")
	runSample(t, createClients, dbName, "failed to create clients")
	runSample(t, write, dbName, "failed to insert data")
	runSample(t, addNewColumn, dbName, "failed to add new column")
	runSample(t, delete, dbName, "failed to delete data")
	runSample(t, write, dbName, "failed to insert data")
	runSample(t, update, dbName, "failed to update data")
	out = runSample(t, writeWithTransactionUsingDML, dbName, "failed to write with transaction using DML")
	assertContains(t, out, "Moved 200000 from Album2's MarketingBudget to Album1")
	out = runSample(t, queryNewColumn, dbName, "failed to query new column")
	assertContains(t, out, "1 1 300000")
	assertContains(t, out, "2 2 300000")

	runSample(t, delete, dbName, "failed to delete data")
	runSample(t, write, dbName, "failed to insert data")
	runSample(t, update, dbName, "failed to update data")
	out = runSample(t, writeWithTransaction, dbName, "failed to write with transaction")
	assertContains(t, out, "Moved 200000 from Album2's MarketingBudget to Album1")
	out = runSample(t, queryNewColumn, dbName, "failed to query new column")
	assertContains(t, out, "1 1 300000")
	assertContains(t, out, "2 2 300000")

	runSample(t, delete, dbName, "failed to delete data")
	runSample(t, write, dbName, "failed to insert data")
	writeTime := time.Now()

	out = runSample(t, read, dbName, "failed to read data")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, query, dbName, "failed to query data")
	assertContains(t, out, "1 1 Total Junk")

	runSample(t, addIndex, dbName, "failed to add index")
	out = runSample(t, queryUsingIndex, dbName, "failed to query using index")
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	if strings.Contains(out, "Green") {
		t.Errorf("got output %q; should not contain Green", out)
	}

	out = runSample(t, readUsingIndex, dbName, "failed to read using index")
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	assertContains(t, out, "Green")

	runSample(t, delete, dbName, "failed to delete data")
	runSample(t, write, dbName, "failed to insert data")
	runSample(t, update, dbName, "failed to update data")

	runSample(t, addStoringIndex, dbName, "failed to add storing index")

	out = runSample(t, readStoringIndex, dbName, "failed to read storing index")
	assertContains(t, out, "500000")
	out = runSample(t, readOnlyTransaction, dbName, "failed to read with ReadOnlyTransaction")
	if strings.Count(out, "Total Junk") != 2 {
		t.Errorf("got output %q; wanted it to contain 2 occurrences of Total Junk", out)
	}

	// Wait at least 15 seconds since the write.
	time.Sleep(time.Until(writeTime.Add(16 * time.Second)))
	out = runSample(t, readStaleData, dbName, "failed to read stale data")
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	assertContains(t, out, "Green")

	out = runSample(t, readBatchData, dbName, "failed to read batch data")
	assertContains(t, out, "1 Marc Richards")

	runSample(t, addCommitTimestamp, dbName, "failed to add commit timestamp")
	runSample(t, updateWithTimestamp, dbName, "failed to update with timestamp")
	out = runSample(t, queryWithTimestamp, dbName, "failed to query with timestamp")
	assertContains(t, out, "1000000")

	runSample(t, writeStructData, dbName, "failed to write struct data")
	out = runSample(t, queryWithStruct, dbName, "failed to query with struct")
	assertContains(t, out, "6")
	out = runSample(t, queryWithArrayOfStruct, dbName, "failed to query with array of struct")
	assertContains(t, out, "6")
	assertContains(t, out, "7")
	assertContains(t, out, "8")
	out = runSample(t, queryWithStructField, dbName, "failed to query with struct field")
	assertContains(t, out, "6")
	out = runSample(t, queryWithNestedStructField, dbName, "failed to query with nested struct field")
	assertContains(t, out, "6 Imagination")
	assertContains(t, out, "9 Imagination")

	runSample(t, createTableDocumentsWithTimestamp, dbName, "failed to create documents table with timestamp")
	runSample(t, writeToDocumentsTable, dbName, "failed to write to documents table")
	runSample(t, updateDocumentsTable, dbName, "failed to update documents table")

	out = runSample(t, queryDocumentsTable, dbName, "failed to query documents table")
	assertContains(t, out, "Hello World 1 Updated")

	runSample(t, createTableDocumentsWithHistoryTable, dbName, "failed to create documents table with history table")
	runSample(t, writeWithHistory, dbName, "failed to write with history")
	runSample(t, updateWithHistory, dbName, "failed to update with history")

	out = runSample(t, queryWithHistory, dbName, "failed to query with history")
	assertContains(t, out, "1 1 Hello World 1 Updated")

	out = runSample(t, insertUsingDML, dbName, "failed to insert using DML")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, setCustomTimeoutAndRetry, dbName, "failed to insert using DML with custom timeout and retry")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, updateUsingDML, dbName, "failed to update using DML")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, deleteUsingDML, dbName, "failed to delete using DML")
	assertContains(t, out, "record(s) deleted")

	out = runSample(t, updateUsingDMLWithTimestamp, dbName, "failed to update using DML with timestamp")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, writeAndReadUsingDML, dbName, "failed to write and read using DML")
	assertContains(t, out, "Found record name with ")

	out = runSample(t, updateUsingDMLStruct, dbName, "failed to update using DML with struct")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, writeUsingDML, dbName, "failed to write using DML")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, commitStats, dbName, "failed to request commit stats")
	assertContains(t, out, "3 mutations in transaction")

	out = runSample(t, queryWithParameter, dbName, "failed to query with parameter")
	assertContains(t, out, "12 Melissa Garcia")

	out = runSample(t, updateUsingPartitionedDML, dbName, "failed to update using partitioned DML")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, deleteUsingPartitionedDML, dbName, "failed to delete using partitioned DML")
	assertContains(t, out, "record(s) deleted")

	out = runSample(t, updateUsingBatchDML, dbName, "failed to update using batch DML")
	assertContains(t, out, "Executed 2 SQL statements using Batch DML.")

	out = runSample(t, createTableWithDatatypes, dbName, "failed to create table with data types")
	assertContains(t, out, "Created Venues table")

	runSample(t, writeDatatypesData, dbName, "failed to write data with different data types")
	out = runSample(t, queryWithArray, dbName, "failed to query with array")
	assertContains(t, out, "19 Venue 19 2020-11-01")
	assertContains(t, out, "42 Venue 42 2020-10-01")

	out = runSample(t, queryWithBool, dbName, "failed to query with bool")
	assertContains(t, out, "19 Venue 19 true")

	out = runSample(t, queryWithBytes, dbName, "failed to query with bytes")
	assertContains(t, out, "4 Venue 4")

	out = runSample(t, queryWithDate, dbName, "failed to query with date")
	assertContains(t, out, "4 Venue 4 2018-09-02")
	assertContains(t, out, "42 Venue 42 2018-10-01")

	out = runSample(t, queryWithFloat, dbName, "failed to query with float")
	assertContains(t, out, "4 Venue 4 0.8")
	assertContains(t, out, "19 Venue 19 0.9")

	out = runSample(t, queryWithInt, dbName, "failed to query with int")
	assertContains(t, out, "19 Venue 19 6300")
	assertContains(t, out, "42 Venue 42 3000")

	out = runSample(t, queryWithString, dbName, "failed to query with string")
	assertContains(t, out, "42 Venue 42")

	// Wait 5 seconds to avoid a time drift issue for the next query:
	// https://github.com/GoogleCloudPlatform/golang-samples/issues/1146.
	time.Sleep(time.Second * 5)
	out = runSample(t, queryWithTimestampParameter, dbName, "failed to query with timestamp parameter")
	assertContains(t, out, "4 Venue 4")
	assertContains(t, out, "19 Venue 19")
	assertContains(t, out, "42 Venue 42")
	out = runSample(t, queryWithQueryOptions, dbName, "failed to query with query options")
	assertContains(t, out, "4 Venue 4")
	assertContains(t, out, "19 Venue 19")
	assertContains(t, out, "42 Venue 42")
	out = runSample(t, createClientWithQueryOptions, dbName, "failed to create a client with query options")
	assertContains(t, out, "4 Venue 4")
	assertContains(t, out, "19 Venue 19")
	assertContains(t, out, "42 Venue 42")

	runSample(t, dropColumn, dbName, "failed to drop column")
	runSample(t, addNumericColumn, dbName, "failed to add numeric column")
	runSample(t, updateDataWithNumericColumn, dbName, "failed to update data with numeric")
	out = runSample(t, queryWithNumericParameter, dbName, "failed to query with numeric parameter")
	assertContains(t, out, "4 ")
	assertContains(t, out, "35000")
}

func TestBackupSample(t *testing.T) {
	_ = testutil.EndToEndTest(t)

	id := randomID()
	dbName, cleanup := initTest(t, id)
	defer cleanup()
	restoreDBName, backupID, cancelledBackupID, cleanupBackup := initBackupTest(t, id, dbName)

	var out string
	// Set up the database for testing backup operations.
	mustRunSample(t, createDatabase, dbName, "failed to create a database")
	runSample(t, write, dbName, "failed to insert data")

	// Start testing backup operations.
	out = runBackupSample(t, createBackup, dbName, backupID, "failed to create a backup")
	assertContains(t, out, fmt.Sprintf("backups/%s", backupID))

	out = runBackupSample(t, cancelBackup, dbName, cancelledBackupID, "failed to cancel a backup")
	assertContains(t, out, "Backup cancelled.")

	out = runBackupSample(t, listBackups, dbName, backupID, "failed to list backups")
	assertContains(t, out, fmt.Sprintf("/backups/%s", backupID))
	assertContains(t, out, "Backups listed.")

	out = runSample(t, listBackupOperations, dbName, "failed to list backup operations")
	assertContains(t, out, fmt.Sprintf("on database %s", dbName))

	out = runBackupSample(t, updateBackup, dbName, backupID, "failed to update a backup")
	assertContains(t, out, fmt.Sprintf("Updated backup %s", backupID))

	out = runBackupSampleWithRetry(t, restoreBackup, restoreDBName, backupID, "failed to restore a backup", 10)
	assertContains(t, out, fmt.Sprintf("Source database %s restored from backup", dbName))

	// This sample should run after a restore operation.
	out = runSample(t, listDatabaseOperations, restoreDBName, "failed to list database operations")
	assertContains(t, out, fmt.Sprintf("Database %s restored from backup", restoreDBName))

	// Delete the restore DB.
	cleanupBackup()

	out = runBackupSample(t, deleteBackup, dbName, backupID, "failed to delete a backup")
	assertContains(t, out, fmt.Sprintf("Deleted backup %s", backupID))
}

func TestCreateDatabaseWithRetentionPeriodSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	wantRetentionPeriod := "7d"
	out := runSample(t, createDatabaseWithRetentionPeriod, dbName, "failed to create a database with a retention period")
	assertContains(t, out, fmt.Sprintf("Created database [%s] with version retention period %q", dbName, wantRetentionPeriod))
}

func runSample(t *testing.T, f sampleFunc, dbName, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, dbName); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runBackupSample(t *testing.T, f backupSampleFunc, dbName, backupID, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, dbName, backupID); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runBackupSampleWithRetry(t *testing.T, f backupSampleFunc, dbName, backupID, errMsg string, maxAttempts int) string {
	var b bytes.Buffer
	testutil.Retry(t, maxAttempts, time.Minute, func(r *testutil.R) {
		b.Reset()
		if err := f(&b, dbName, backupID); err != nil {
			if spanner.ErrCode(err) == codes.InvalidArgument && strings.Contains(err.Error(), "Please retry the operation once the pending restores complete") {
				r.Errorf("%s: %v", errMsg, err)
			} else {
				t.Fatalf("%s: %v", errMsg, err)
			}
		}
	})
	return b.String()
}

func runInstanceSample(t *testing.T, f instanceSampleFunc, projectID, instanceID, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, projectID, instanceID); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func mustRunSample(t *testing.T, f sampleFunc, dbName, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, dbName); err != nil {
		t.Fatalf("%s: %v", errMsg, err)
	}
	return b.String()
}

func getInstance(t *testing.T) string {
	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("Skipping spanner integration test. Set GOLANG_SAMPLES_SPANNER.")
	}
	if !strings.HasPrefix(instance, "projects/") {
		t.Fatal("Spanner instance ref must be in the form of 'projects/PROJECT_ID/instances/INSTANCE_ID'")
	}
	return instance
}

func assertContains(t *testing.T, out string, sub string) {
	t.Helper()
	if !strings.Contains(out, sub) {
		t.Errorf("got output %q; want it to contain %q", out, sub)
	}
}

// Maximum length of database name is 30 characters, so trim if the generated name is too long
func validLength(databaseName string, t *testing.T) (trimmedName string) {
	if len(databaseName) > 30 {
		trimmedName := databaseName[:30]
		t.Logf("Name too long, '%s' trimmed to '%s'", databaseName, trimmedName)
		return trimmedName
	}

	return databaseName
}

func cleanupInstance(projectID, instanceID string) error {
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot create instance databaseAdmin client: %v", err)
	}
	defer instanceAdmin.Close()

	instanceName := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	if err := instanceAdmin.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: instanceName}); err != nil {
		return fmt.Errorf("failed to delete instance %s (error %v), might need a manual removal",
			instanceName, err)
	}
	return nil
}

func randomID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%s-%s", strconv.FormatInt(now.Unix(), 10), uuid.New().String()[:8])
}

func parseInstanceName(instanceName string) (project, instance string, err error) {
	matches := validInstancePattern.FindStringSubmatch(instanceName)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("failed to parse database name from %q according to pattern %q",
			instanceName, validInstancePattern.String())
	}
	return matches[1], matches[2], nil
}
