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
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sampleFunc func(w io.Writer, dbName string) error
type sampleFuncWithContext func(ctx context.Context, w io.Writer, dbName string) error
type instanceSampleFunc func(w io.Writer, projectID, instanceID string) error
type backupSampleFunc func(ctx context.Context, w io.Writer, dbName, backupID string) error
type backupSampleFuncWithoutContext func(w io.Writer, dbName, backupID string) error
type backupScheduleSampleFunc func(w io.Writer, dbName string, scheduleId string) error
type createBackupSampleFunc func(ctx context.Context, w io.Writer, dbName, backupID string, versionTime time.Time) error
type instancePartitionSampleFunc func(w io.Writer, projectID, instanceID, instancePartitionID string) error

var (
	validInstancePattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)$")
)

func initTest(t *testing.T, id string) (instName, dbName string, cleanup func()) {
	projectID := getSampleProjectId(t)
	configName := getSamplesInstanceConfig()
	if configName == "" {
		configName = "regional-us-central1"
	}
	log.Printf("Running test by using the instance config: %s\n", configName)
	instanceID, cleanup := createTestInstance(t, projectID, configName)
	instName = fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	dbID := validLength(fmt.Sprintf("smpl-%s", id), t)
	dbName = fmt.Sprintf("%s/databases/%s", instName, dbID)

	return
}

func initTestWithConfig(t *testing.T, id string, instanceConfigName string) (instName, dbName string, cleanup func()) {
	projectID := getSampleProjectId(t)
	instanceID, cleanup := createTestInstance(t, projectID, instanceConfigName)
	instName = fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	dbID := validLength(fmt.Sprintf("smpl-%s", id), t)
	dbName = fmt.Sprintf("%s/databases/%s", instName, dbID)

	return
}

func getVersionTime(t *testing.T, dbName string) (versionTime time.Time) {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, dbName)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT CURRENT_TIMESTAMP()`,
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err != nil {
		t.Fatalf("failed to get current time: %v", err)
	}
	if err := row.Columns(&versionTime); err != nil {
		t.Fatalf("failed to get version time: %v", err)
	}

	return versionTime
}

func initBackupTest(t *testing.T, id, instName string) (restoreDBName, backupID, cancelledBackupID string) {
	restoreDatabaseID := validLength(fmt.Sprintf("restore-%s", id), t)
	restoreDBName = fmt.Sprintf("%s/databases/%s", instName, restoreDatabaseID)
	backupID = validLength(fmt.Sprintf("backup-%s", id), t)
	cancelledBackupID = validLength(fmt.Sprintf("cancel-%s", id), t)

	return
}

func initInstancePartitionTest(t *testing.T, id string) (string, string, string, func()) {
	projectID := getSampleProjectId(t)
	instancePartitionID := fmt.Sprintf("instance-partition-%s", id)
	instanceID, cleanup := createTestInstance(t, projectID, "regional-us-central1")

	return projectID, instanceID, instancePartitionID, cleanup
}

func TestCreateInstances(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	runCreateAndUpdateInstanceSample(t, createInstance, updateInstance)
	runCreateAndUpdateInstanceSample(t, createInstanceWithoutDefaultBackupSchedule, updateInstanceDefaultBackupScheduleType)
	runCreateInstanceSample(t, createInstanceWithProcessingUnits)
	runCreateInstanceSample(t, createInstanceWithAutoscalingConfig)
}

func runCreateAndUpdateInstanceSample(t *testing.T, createFunc, updateFunc instanceSampleFunc) {
	projectID := getSampleProjectId(t)
	instanceID := fmt.Sprintf("go-sample-test-%s", uuid.New().String()[:8])
	out := runInstanceSample(t, createFunc, projectID, instanceID, "failed to create an instance")
	assertContains(t, out, fmt.Sprintf("Created instance [%s]", instanceID))
	out = runInstanceSample(t, updateFunc, projectID, instanceID, "failed to update an instance")
	assertContains(t, out, fmt.Sprintf("Updated instance [%s]", instanceID))
	if err := cleanupInstance(projectID, instanceID); err != nil {
		t.Logf("cleanupInstance error: %s", err)
	}
}

func runCreateInstanceSample(t *testing.T, f instanceSampleFunc) {
	projectID := getSampleProjectId(t)
	instanceID := fmt.Sprintf("go-sample-test-%s", uuid.New().String()[:8])
	out := runInstanceSample(t, f, projectID, instanceID, "failed to create an instance")
	if err := cleanupInstance(projectID, instanceID); err != nil {
		t.Logf("cleanupInstance error: %s", err)
	}
	assertContains(t, out, fmt.Sprintf("Created instance [%s]", instanceID))
}

func TestSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	var out string
	mustRunSample(t, createDatabase, dbName, "failed to create a database")
	runSample(t, createClients, dbName, "failed to create clients")
	runSample(t, write, dbName, "failed to insert data")
	out = runSample(t, batchWrite, dbName, "failed to write data using BatchWrite")
	assertNotContains(t, out, "could not be applied with error")
	runSampleWithContext(ctx, t, addNewColumn, dbName, "failed to add new column")
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

	runSample(t, addAndDropDatabaseRole, dbName, "failed to add database role")
	out = runSample(t, func(w io.Writer, dbName string) error { return readDataWithDatabaseRole(w, dbName, "parent") }, dbName, "failed to read data with database role")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, listDatabaseRoles, dbName, "failed to list database roles")
	assertContains(t, out, "parent")
	assertContains(t, out, "public")
	assertContains(t, out, "spanner_info_reader")
	assertContains(t, out, "spanner_sys_reader")

	out = runSample(t, read, dbName, "failed to read data")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, query, dbName, "failed to query data")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, queryRequestPriority, dbName, "failed to query data with RequestPriority")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, queryWithTag, dbName, "failed to query data with request tag set")
	assertContains(t, out, "1 1 Total Junk")

	runSampleWithContext(ctx, t, addIndex, dbName, "failed to add index")
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

	runSampleWithContext(ctx, t, addStoringIndex, dbName, "failed to add storing index")

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

	out = runSample(t, readRequestPriority, dbName, "failed to read with RequestPriority")
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	assertContains(t, out, "Green")

	out = runSample(t, readBatchData, dbName, "failed to read batch data")
	assertContains(t, out, "1 Marc Richards")

	out = runSample(t, readBatchDataRequestPriority, dbName, "failed to read batch data with RequestPriority")
	assertContains(t, out, "1 Marc Richards")

	runSampleWithContext(ctx, t, addCommitTimestamp, dbName, "failed to add commit timestamp")
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

	runSampleWithContext(ctx, t, createTableDocumentsWithTimestamp, dbName, "failed to create documents table with timestamp")
	runSample(t, writeToDocumentsTable, dbName, "failed to write to documents table")
	runSample(t, updateDocumentsTable, dbName, "failed to update documents table")

	out = runSample(t, queryDocumentsTable, dbName, "failed to query documents table")
	assertContains(t, out, "Hello World 1 Updated")

	runSampleWithContext(ctx, t, createTableDocumentsWithHistoryTable, dbName, "failed to create documents table with history table")
	runSample(t, writeWithHistory, dbName, "failed to write with history")
	runSample(t, updateWithHistory, dbName, "failed to update with history")

	out = runSample(t, queryWithHistory, dbName, "failed to query with history")
	assertContains(t, out, "1 1 Hello World 1 Updated")

	out = runSample(t, insertUsingDML, dbName, "failed to insert using DML")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, insertUsingDMLReturning, dbName, "failed to insert using DML with returning clause")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, insertUsingDMLRequestPriority, dbName, "failed to insert using DML with RequestPriority")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, setCustomTimeoutAndRetry, dbName, "failed to insert using DML with custom timeout and retry")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, setStatementTimeout, dbName, "failed to execute statement with a timeout")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, transactionTimeout, dbName, "failed to run transaction with a timeout")
	assertContains(t, out, "Transaction with timeout was executed successfully")

	out = runSample(t, updateUsingDML, dbName, "failed to update using DML")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, updateUsingDMLReturning, dbName, "failed to update using DML with returning clause")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, deleteUsingDML, dbName, "failed to delete using DML")
	assertContains(t, out, "record(s) deleted")

	out = runSample(t, deleteUsingDMLReturning, dbName, "failed to delete using DML with returning clause")
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
	assertContains(t, out, "4 mutations in transaction")

	out = runSample(t, maxCommitDelay, dbName, "failed to set max commit delay")
	assertContains(t, out, "4 mutations in transaction")

	out = runSample(t, queryWithParameter, dbName, "failed to query with parameter")
	assertContains(t, out, "12 Melissa Garcia")

	out = runSample(t, updateUsingPartitionedDML, dbName, "failed to update using partitioned DML")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, updateUsingPartitionedDMLRequestPriority, dbName, "failed to update using partitioned DML with RequestPriority")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, deleteUsingPartitionedDML, dbName, "failed to delete using partitioned DML")
	assertContains(t, out, "record(s) deleted")

	out = runSample(t, updateUsingBatchDML, dbName, "failed to update using batch DML")
	assertContains(t, out, "Executed 2 SQL statements using Batch DML.")

	out = runSample(t, updateUsingBatchDMLRequestPriority, dbName, "failed to update using batch DML with RequestPriority")
	assertContains(t, out, "Executed 2 SQL statements using Batch DML.")

	out = runSampleWithContext(ctx, t, createTableWithDatatypes, dbName, "failed to create table with data types")
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

	out = runSample(t, readWriteTransactionWithTag, dbName, "failed to perform read-write transaction with tag")
	assertContains(t, out, "Venue capacities updated.")
	assertContains(t, out, "New venue inserted.")
	out = runSample(t, queryWithInt, dbName, "failed to query with int")
	assertContains(t, out, "19 Venue 19 6300")
	assertNotContains(t, out, "42 Venue 42 3000")

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

	out = runSample(t, queryWithGFELatency, dbName, "failed to query with GFE latency")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, queryWithGRPCMetric, dbName, "failed to query with gRPC metric")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, queryWithQueryStats, dbName, "failed to query with query stats")
	assertContains(t, out, "1 1 Total Junk")

	runSample(t, dropColumn, dbName, "failed to drop column")
	runSampleWithContext(ctx, t, addNumericColumn, dbName, "failed to add numeric column")
	runSample(t, updateDataWithNumericColumn, dbName, "failed to update data with numeric")
	out = runSample(t, queryWithNumericParameter, dbName, "failed to query with numeric parameter")
	assertContains(t, out, "4 ")
	assertContains(t, out, "35000")

	out = runSample(t, addJsonColumn, dbName, "failed to add json column")
	assertContains(t, out, "Added VenueDetails column\n")
	out = runSample(t, updateDataWithJsonColumn, dbName, "failed to update data with json")
	assertContains(t, out, "Updated data to VenueDetails column\n")
	out = runSample(t, queryWithJsonParameter, dbName, "failed to query with json parameter")
	assertContains(t, out, "The venue details for venue id 19")

	out = runSample(t, createSequence, dbName, "failed to create table with bit reverse sequence enabled")
	assertContains(t, out, "Created Seq sequence and Customers table, where the key column CustomerId uses the sequence as a default value\n")
	assertContains(t, out, "Inserted customer record with CustomerId")
	assertContains(t, out, "Number of customer records inserted is: 3")
	out = runSample(t, alterSequence, dbName, "failed to alter table with bit reverse sequence enabled")
	assertContains(t, out, "Altered Seq sequence to skip an inclusive range between 1000 and 5000000\n")
	assertContains(t, out, "Inserted customer record with CustomerId")
	assertContains(t, out, "Number of customer records inserted is: 3")
	out = runSample(t, dropSequence, dbName, "failed to drop bit reverse sequence column")
	assertContains(t, out, "Altered Customers table to drop DEFAULT from CustomerId column and dropped the Seq sequence\n")

	out = runSample(t, directedReadOptions, dbName, "failed to read using directed read options")
	assertContains(t, out, "1 1 Total Junk")
}

func TestBackupSample(t *testing.T) {
	t.Skip("https://github.com/GoogleCloudPlatform/golang-samples/issues/2333")
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip("GOLANG_SAMPLES_E2E_TEST not set")
	}
	_ = testutil.SystemTest(t)
	t.Parallel()

	id := randomID()
	instName, dbName, cleanup := initTest(t, id)
	defer cleanup()
	restoreDBName, backupID, cancelledBackupID := initBackupTest(t, id, instName)

	var out string
	// Set up the database for testing backup operations.
	mustRunSample(t, createDatabase, dbName, "failed to create a database")
	runSample(t, write, dbName, "failed to insert data")

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	// Start testing backup operations.
	versionTime := getVersionTime(t, dbName)
	out = runCreateBackupSample(ctx, t, createBackup, dbName, backupID, versionTime, "failed to create a backup")
	assertContains(t, out, fmt.Sprintf("backups/%s", backupID))

	out = runBackupSample(ctx, t, cancelBackup, dbName, cancelledBackupID, "failed to cancel a backup")
	assertContains(t, out, "Backup cancelled.")

	out = runBackupSample(ctx, t, listBackups, dbName, backupID, "failed to list backups")
	assertContains(t, out, fmt.Sprintf("/backups/%s", backupID))
	assertContains(t, out, "Backups listed.")

	out = runBackupSampleWithoutContext(t, listBackupOperations, dbName, backupID, "failed to list backup operations")
	assertContains(t, out, fmt.Sprintf("on database %s", dbName))
	assertContains(t, out, fmt.Sprintf("copied from %s", backupID))

	out = runBackupSampleWithoutContext(t, updateBackup, dbName, backupID, "failed to update a backup")
	assertContains(t, out, fmt.Sprintf("Updated backup %s", backupID))

	out = runBackupSampleWithRetry(ctx, t, restoreBackup, restoreDBName, backupID, "failed to restore a backup", 10)
	assertContains(t, out, fmt.Sprintf("Source database %s restored from backup", dbName))

	// This sample should run after a restore operation.
	out = runSampleWithContext(ctx, t, listDatabaseOperations, restoreDBName, "failed to list database operations")
	assertContains(t, out, fmt.Sprintf("Database %s restored from backup", restoreDBName))

	out = runBackupSample(ctx, t, deleteBackup, dbName, backupID, "failed to delete a backup")
	assertContains(t, out, fmt.Sprintf("Deleted backup %s", backupID))
}

func TestBackupScheduleSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	// Set up the database for testing backup schedule operations.
	mustRunSample(t, createDatabase, dbName, "failed to create a database")

	var out string
	out = runBackupScheduleSample(t, createFullBackupSchedule, dbName, "full-backup-schedule", "failed to create full backup schedule")
	assertContains(t, out, "Created full backup schedule")
	assertContains(t, out, fmt.Sprintf("%s/backupSchedules/%s", dbName, "full-backup-schedule"))

	out = runBackupScheduleSample(t, createIncrementalBackupSchedule, dbName, "incremental-backup-schedule", "failed to create incremental backup schedule")
	assertContains(t, out, "Created incremental backup schedule")
	assertContains(t, out, fmt.Sprintf("%s/backupSchedules/%s", dbName, "incremental-backup-schedule"))

	out = runSample(t, listBackupSchedules, dbName, "failed to list backup schedule")
	assertContains(t, out, fmt.Sprintf("%s/backupSchedules/%s", dbName, "full-backup-schedule"))
	assertContains(t, out, fmt.Sprintf("%s/backupSchedules/%s", dbName, "incremental-backup-schedule"))

	out = runBackupScheduleSample(t, getBackupSchedule, dbName, "full-backup-schedule", "failed to get backup schedule")
	assertContains(t, out, fmt.Sprintf("%s/backupSchedules/%s", dbName, "full-backup-schedule"))

	out = runBackupScheduleSample(t, updateBackupSchedule, dbName, "full-backup-schedule", "failed to update backup schedule")
	assertContains(t, out, "Updated backup schedule")
	assertContains(t, out, fmt.Sprintf("%s/backupSchedules/%s", dbName, "full-backup-schedule"))

	out = runBackupScheduleSample(t, deleteBackupSchedule, dbName, "full-backup-schedule", "failed to delete backup schedule")
	assertContains(t, out, "Deleted backup schedule")
}

func TestInstancePartitionSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	id := randomID()
	projectID, instanceID, instancePartitionID, cleanup := initInstancePartitionTest(t, id)
	defer cleanup()

	var out string
	out = runInstancePartitionSample(t, createInstancePartition, projectID, instanceID, instancePartitionID, "failed to create an instance partition")
	assertContains(t, out, fmt.Sprintf("Created instance partition [%s]", instancePartitionID))
}

func TestCreateDatabaseWithRetentionPeriodSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	wantRetentionPeriod := "7d"
	out := runSampleWithContext(ctx, t, createDatabaseWithRetentionPeriod, dbName, "failed to create a database with a retention period")
	assertContains(t, out, fmt.Sprintf("Created database [%s] with version retention period %q", dbName, wantRetentionPeriod))
}

func TestCustomerManagedEncryptionKeys(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip("GOLANG_SAMPLES_E2E_TEST not set")
	}
	tc := testutil.SystemTest(t)
	t.Parallel()
	startTime := time.Now()
	instName, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	var b bytes.Buffer

	locationId := "us-west1"
	keyRingId := "spanner-test-keyring"
	keyId := "spanner-test-key"

	// Create an encryption key if it does not already exist.
	if err := maybeCreateKey(tc.ProjectID, locationId, keyRingId, keyId); err != nil {
		t.Errorf("failed to create encryption key: %v", err)
	}
	kmsKeyName := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		tc.ProjectID,
		locationId,
		keyRingId,
		keyId,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	// Create an encrypted database. The database is automatically deleted by the cleanup function.
	if err := createDatabaseWithCustomerManagedEncryptionKey(ctx, &b, dbName, kmsKeyName); err != nil {
		t.Errorf("failed to create database with customer managed encryption key: %v", err)
	}
	out := b.String()
	assertContains(t, out, fmt.Sprintf("Created database [%s] using encryption key %q", dbName, kmsKeyName))
	t.Logf("create database operation took: %v\n", time.Since(startTime))

	// Try to create a backup of the encrypted database and delete it after the test.
	backupId := fmt.Sprintf("enc-backup-%s", randomID())
	b.Reset()
	if err := createBackupWithCustomerManagedEncryptionKey(ctx, &b, dbName, backupId, kmsKeyName); err != nil {
		t.Errorf("failed to create backup with customer managed encryption key: %v", err)
	}
	out = b.String()
	assertContains(t, out, fmt.Sprintf("backups/%s", backupId))
	assertContains(t, out, fmt.Sprintf("using encryption key %s", kmsKeyName))
	t.Logf("create backup operation took: %v\n", time.Since(startTime))

	// Try to restore the encrypted database and delete the restored database after the test.
	restoredName := fmt.Sprintf("%s/databases/rest-enc-%s", instName, randomID())
	restoreFunc := func(ctx context.Context, w io.Writer, dbName, backupID string) error {
		return restoreBackupWithCustomerManagedEncryptionKey(ctx, w, dbName, backupId, kmsKeyName)
	}
	out = runBackupSampleWithRetry(ctx, t, restoreFunc, restoredName, backupId, "failed to restore database with customer managed encryption key", 10)
	assertContains(t, out, fmt.Sprintf("Database %s restored", dbName))
	assertContains(t, out, fmt.Sprintf("using encryption key %s", kmsKeyName))
	t.Logf("restore backup operation took: %v\n", time.Since(startTime))
}

func TestCustomerManagedMultiRegionEncryptionKeys(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip("GOLANG_SAMPLES_E2E_TEST not set")
	}
	tc := testutil.SystemTest(t)
	t.Parallel()
	startTime := time.Now()
	instName, dbName, cleanup := initTestWithConfig(t, randomID(), "nam3")
	defer cleanup()

	var b bytes.Buffer
	var kmsKeyNames []string
	keyRingId := "spanner-test-keyring"
	keyId := "spanner-test-key"
	for _, locationId := range []string{"us-central1", "us-east1", "us-east4"} {
		// Create an encryption key if it does not already exist.
		if err := maybeCreateKey(tc.ProjectID, locationId, keyRingId, keyId); err != nil {
			t.Errorf("failed to create encryption key: %v", err)
		}
		kmsKeyNames = append(kmsKeyNames, fmt.Sprintf(
			"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
			tc.ProjectID,
			locationId,
			keyRingId,
			keyId,
		))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	// Create an encrypted database. The database is automatically deleted by the cleanup function.
	if err := createDatabaseWithCustomerManagedMultiRegionEncryptionKey(ctx, &b, dbName, kmsKeyNames); err != nil {
		t.Errorf("failed to create database with customer managed multi-region encryption keys: %v", err)
	}
	out := b.String()
	assertContains(t, out, fmt.Sprintf("Created database [%s] using multi-region encryption keys %q", dbName, kmsKeyNames))
	t.Logf("create database operation took: %v\n", time.Since(startTime))

	// Try to create a backup of the encrypted database and delete it after the test.
	backupId := fmt.Sprintf("enc-backup-%s", randomID())
	b.Reset()
	if err := createBackupWithCustomerManagedMultiRegionEncryptionKey(ctx, &b, dbName, backupId, kmsKeyNames); err != nil {
		t.Errorf("failed to create backup with customer managed multi-region encryption keys: %v", err)
	}
	out = b.String()
	assertContains(t, out, fmt.Sprintf("backups/%s", backupId))
	assertContains(t, out, fmt.Sprintf("using multi-region encryption keys %q", kmsKeyNames))
	t.Logf("create backup operation took: %v\n", time.Since(startTime))

	// Try to create copy of a backup of the encrypted database and delete it after the test.
	copyBackupId := fmt.Sprintf("copy-enc-backup-%s", randomID())
	b.Reset()
	if err := copyBackupWithMultiRegionEncryptionKey(&b, instName, copyBackupId, fmt.Sprintf("%s/backups/%s", instName, backupId), kmsKeyNames); err != nil {
		t.Errorf("failed to copy backup with customer managed multi-region encryption keys: %v", err)
	}
	out = b.String()
	assertContains(t, out, fmt.Sprintf("backups/%s", backupId))
	assertContains(t, out, "multi-region encryption keys\n")
	t.Logf("copy backup operation took: %v\n", time.Since(startTime))

	// Try to restore the encrypted database and delete the restored database after the test.
	restoredName := fmt.Sprintf("%s/databases/rest-enc-%s", instName, randomID())
	restoreFunc := func(ctx context.Context, w io.Writer, dbName, backupID string) error {
		return restoreBackupWithCustomerManagedMultiRegionEncryptionKey(ctx, w, dbName, backupId, kmsKeyNames)
	}
	out = runBackupSampleWithRetry(ctx, t, restoreFunc, restoredName, backupId, "failed to restore database with customer managed multi-region encryption keys", 10)
	assertContains(t, out, fmt.Sprintf("Database %s restored", dbName))
	assertContains(t, out, fmt.Sprintf("using multi-region encryption key %s", kmsKeyNames))
	t.Logf("restore backup operation took: %v\n", time.Since(startTime))
}

func TestCreateDatabaseWithDefaultLeaderSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	instName, dbName, cleanup := initTestWithConfig(t, randomID(), "nam3")
	defer cleanup()

	projectID := getSampleProjectId(t)
	var b bytes.Buffer

	// Try to get Instance Configs
	config := fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "nam3")
	if err := getInstanceConfig(&b, config); err != nil {
		t.Errorf("failed to create get instance configs: %v", err)
	}
	out := b.String()
	assertContains(t, out, "Available leader options for instance config")

	// Try to list Instance Configs
	b.Reset()
	if err := listInstanceConfigs(&b, "projects/"+projectID); err != nil {
		t.Errorf("failed to list instance configs: %v", err)
	}
	out = b.String()
	assertContains(t, out, "Available leader options for instance config")

	// Try to get list of Databases
	b.Reset()
	if err := listDatabases(&b, instName); err != nil {
		t.Errorf("failed to get list of Databases: %v", err)
	}
	out = b.String()
	assertContains(t, out, "Databases for instance")

	// Try to create Database with Default Leader
	b.Reset()
	defaultLeader := "us-east1"
	if err := createDatabaseWithDefaultLeader(&b, dbName, defaultLeader); err != nil {
		t.Errorf("failed to create database with default leader: %v", err)
	}
	out = b.String()
	assertContains(t, out, fmt.Sprintf("Created database [%s] with default leader%q\n", dbName, defaultLeader))

	// Try to update Database with Default Leader
	b.Reset()
	defaultLeader = "us-east4"
	if err := updateDatabaseWithDefaultLeader(&b, dbName, defaultLeader); err != nil {
		t.Errorf("failed to update database with default leader: %v", err)
	}
	out = b.String()
	assertContains(t, out, "Updated the default leader\n")

	// Try to get Database DDL
	b.Reset()
	if err := getDatabaseDdl(&b, dbName); err != nil {
		t.Errorf("failed to get Database DDL: %v", err)
	}
	out = b.String()
	assertContains(t, out, "Database DDL is as follows")

	// Try to Query Information Schema Database Options
	b.Reset()
	if err := queryInformationSchemaDatabaseOptions(&b, dbName); err != nil {
		t.Errorf("failed to query information schema database options: %v", err)
	}
	out = b.String()
	assertContains(t, out, "The result of the query to get")
}

func TestUpdateDatabaseSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	var out string
	mustRunSample(t, createDatabase, dbName, "failed to create a database")
	out = runSampleWithContext(ctx, t, updateDatabase, dbName, "failed to update database")
	assertContains(t, out, fmt.Sprintf("Updated database [%s]\n", dbName))

	databaseAdmin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatalf("cannot create databaseAdmin client: %v", err)
	}

	database, err := databaseAdmin.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbName})
	if err != nil {
		log.Fatalf("error during GetDatabase for db %v: %v", dbName, err)
	}
	// Verify if the EnableDropProtection field got enabled
	if !database.GetEnableDropProtection() {
		t.Errorf("got output %t; want it to contain %t", database.GetEnableDropProtection(), true)
	}

	// Disable Drop db protection for the cleanup to delete the databases
	testutil.Retry(t, 20, time.Minute, func(r *testutil.R) {
		opUpdate, err := databaseAdmin.UpdateDatabase(ctx, &adminpb.UpdateDatabaseRequest{
			Database: &adminpb.Database{
				Name:                 dbName,
				EnableDropProtection: false,
			},
			UpdateMask: &field_mask.FieldMask{
				Paths: []string{"enable_drop_protection"},
			},
		})
		if err != nil {
			// Retry if the database could not be updated due to some transient error
			r.Errorf("UpdateDatabase operation to DB %v failed: %v", dbName, err)
			return
		}
		if _, err := opUpdate.Wait(ctx); err != nil {
			t.Fatalf("UpdateDatabase operation to DB %v failed: %v", dbName, err)
		}
	})
}

func TestCustomInstanceConfigSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	projectID := getSampleProjectId(t)
	defer cleanupInstanceConfigs(projectID)

	var b bytes.Buffer
	userConfigID := fmt.Sprintf("custom-golang-samples-config-%v", randomID())
	if err := createInstanceConfig(&b, projectID, userConfigID, "nam11"); err != nil {
		t.Fatalf("failed to create instance configuration: %v", err)
	}
	out := b.String()
	assertContains(t, out, "Created instance configuration")

	b.Reset()
	if err := updateInstanceConfig(&b, projectID, userConfigID); err != nil {
		t.Errorf("failed to update instance configuration: %v", err)
	}
	out = b.String()
	assertContains(t, out, "Updated instance configuration")

	b.Reset()
	if err := listInstanceConfigOperations(&b, projectID); err != nil {
		t.Errorf("failed to list instance configuration operations: %v", err)
	}
	out = b.String()
	assertContains(t, out, "List instance config operations")

	b.Reset()
	if err := deleteInstanceConfig(&b, projectID, userConfigID); err != nil {
		t.Errorf("failed to delete instance configuration: %v", err)
	}
	out = b.String()
	assertContains(t, out, "Deleted instance configuration")
}

func TestForeignKeyDeleteCascadeSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	mustRunSample(t, createDatabase, dbName, "failed to create a database")

	var out string

	out = runSample(t, createTableWithForeignKeyDeleteCascade, dbName, "failed to create table with foreign key delete constraint")
	assertContains(t, out, "Created Customers and ShoppingCarts table with FKShoppingCartsCustomerId foreign key constraint")
	out = runSample(t, alterTableWithForeignKeyDeleteCascade, dbName, "failed to alter table with foreign key delete constraint")
	assertContains(t, out, "Altered ShoppingCarts table with FKShoppingCartsCustomerName foreign key constraint")
	out = runSample(t, dropForeignKeyDeleteCascade, dbName, "failed to drop foreign key delete constraint")
	assertContains(t, out, "Altered ShoppingCarts table to drop FKShoppingCartsCustomerName foreign key constraint")
}

func TestPgSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	out := runSampleWithContext(ctx, t, pgCreateDatabase, dbName, "failed to create a Spanner PG database")
	assertContains(t, out, fmt.Sprintf("Created Spanner PostgreSQL database [%s]", dbName))

	out = runSampleWithContext(ctx, t, pgAddNewColumn, dbName, "failed to add new column in Spanner PG database")
	assertContains(t, out, "Added MarketingBudget column")

	runSample(t, write, dbName, "failed to insert data in Spanner PG database")
	runSample(t, update, dbName, "failed to update data in Spanner PG database")

	runSampleWithContext(ctx, t, pgAddStoringIndex, dbName, "failed to add storing index in Spanner PG database")
	out = runSample(t, readStoringIndex, dbName, "failed to read storing index in Spanner PG database")
	assertContains(t, out, "500000")

	out = runSample(t, pgWriteWithTransactionUsingDML, dbName, "failed to write with transaction using DML in Spanner PG database")
	assertContains(t, out, "Moved 200000 from Album2's MarketingBudget to Album1")
	out = runSample(t, pgQueryNewColumn, dbName, "failed to query new column in Spanner PG database")
	assertContains(t, out, "1 1 300000")
	assertContains(t, out, "2 2 300000")

	client, err := spanner.NewClient(context.Background(), dbName)
	if err != nil {
		t.Fatalf("failed to create Spanner client: %v", err)
	}
	defer client.Close()
	_, err = client.Apply(context.Background(), []*spanner.Mutation{
		spanner.InsertMap("Venues", map[string]interface{}{
			"VenueId": 4,
			"Name":    "Venue 4",
		}),
		spanner.InsertMap("Venues", map[string]interface{}{
			"VenueId": 19,
			"Name":    "Venue 19",
		}),
		spanner.InsertMap("Venues", map[string]interface{}{
			"VenueId": 42,
			"Name":    "Venue 42",
		}),
	})
	if err != nil {
		t.Fatalf("failed to insert test records: %v", err)
	}
	out = runSample(t, addJsonBColumn, dbName, "failed to add jsonB column")
	assertContains(t, out, "Added VenueDetails column\n")
	out = runSample(t, updateDataWithJsonBColumn, dbName, "failed to update data with jsonB")
	assertContains(t, out, "Updated data to VenueDetails column\n")
	out = runSample(t, queryWithJsonBParameter, dbName, "failed to query with jsonB parameter")
	assertContains(t, out, "The venue details for venue id 19")

	out = runSample(t, pgCreateSequence, dbName, "failed to create table with bit reverse sequence enabled in Spanner PG database")
	assertContains(t, out, "Created Seq sequence and Customers table, where its key column CustomerId uses the sequence as a default value\n")
	assertContains(t, out, "Inserted customer record with CustomerId")
	assertContains(t, out, "Number of customer records inserted is: 3")
	out = runSample(t, pgAlterSequence, dbName, "failed to alter table with bit reverse sequence enabled in Spanner PG database")
	assertContains(t, out, "Altered Seq sequence to skip an inclusive range between 1000 and 5000000\n")
	assertContains(t, out, "Inserted customer record with CustomerId")
	assertContains(t, out, "Number of customer records inserted is: 3")
	out = runSample(t, dropSequence, dbName, "failed to drop bit reverse sequence column in Spanner PG database")
	assertContains(t, out, "Altered Customers table to drop DEFAULT from CustomerId column and dropped the Seq sequence\n")
}

func TestPgQueryParameter(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(
		dbName,
		`CREATE TABLE Singers (
		   SingerId  bigint NOT NULL PRIMARY KEY,
		   FirstName varchar(1024),
		   LastName  varchar(1024)
		 )`)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	client, err := spanner.NewClient(context.Background(), dbName)
	if err != nil {
		t.Fatalf("failed to create Spanner client: %v", err)
	}
	defer client.Close()
	_, err = client.Apply(context.Background(), []*spanner.Mutation{
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			"SingerId":  1,
			"FirstName": "Bruce",
			"LastName":  "Allison",
		}),
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			"SingerId":  2,
			"FirstName": "Alice",
			"LastName":  "Bruxelles",
		}),
		spanner.InsertOrUpdateMap("Singers", map[string]interface{}{
			"SingerId":  12,
			"FirstName": "Melissa",
			"LastName":  "Garcia",
		}),
	})
	if err != nil {
		t.Fatalf("failed to insert test records: %v", err)
	}

	out := runSample(t, pgQueryParameter, dbName, "failed to execute PG query with parameter")
	assertContains(t, out, "12 Melissa Garcia")
	assertNotContains(t, out, "2 Alice Bruxelles")
}

func TestPgDmlSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(
		dbName,
		`CREATE TABLE Singers (
		   SingerId  bigint NOT NULL PRIMARY KEY,
		   FirstName varchar(1024),
		   LastName  varchar(1024),
		   FullName  varchar(2048)
		     GENERATED ALWAYS AS (FirstName || ' ' || LastName) STORED
		 )`,
		`CREATE TABLE Albums (
			SingerId         bigint NOT NULL,
			AlbumId          bigint NOT NULL,
			AlbumTitle       varchar(1024),
			MarketingBudget  bigint,
			PRIMARY KEY (SingerId, AlbumId)
		) INTERLEAVE IN PARENT Singers ON DELETE CASCADE`)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgWriteUsingDML, dbName, "failed to execute PG DML")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, pgDmlWithParameters, dbName, "failed to execute PG DML with parameter")
	assertContains(t, out, "Inserted 2 singers")

	out = runSample(t, pgUpdateUsingDMLReturning, dbName, "failed to execute PG DML update with returning clause")
	assertContains(t, out, "record(s) updated")

	out = runSample(t, pgInsertUsingDMLReturning, dbName, "failed to execute PG DML insert with returning clause")
	assertContains(t, out, "record(s) inserted")

	out = runSample(t, pgDeleteUsingDMLReturning, dbName, "failed to execute PG DML delete with returning clause")
	assertContains(t, out, "record(s) deleted")
}

func TestPgNumericDataType(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgNumericDataType, dbName, "failed to execute PG Numeric sample")
	assertContains(t, out, "Inserted 1 venue(s)")
	assertContains(t, out, "Revenues of Venue 1: 3150.25")
	assertContains(t, out, "Revenues of Venue 2: <null>")
	assertContains(t, out, "Revenues of Venue 3: NaN")
	assertContains(t, out, "Inserted 2 Venues using mutations")
}

func TestPgFunctions(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgFunctions, dbName, "failed to execute PG functions sample")
	assertContains(t, out, "1284352323 seconds after epoch is 2010-09-13 04:32:03 +0000 UTC")
}

func TestPgInformationSchema(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgInformationSchema, dbName, "failed to execute PG INFORMATION_SCHEMA sample")
	assertContains(t, out, "Table: public.venues (User defined type: null)")
}

func TestPgCastDataType(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgCastDataType, dbName, "failed to execute PG cast data type sample")
	assertContains(t, out, "String: 1")
	assertContains(t, out, "Int: 2")
	assertContains(t, out, "Decimal: 3")
	assertContains(t, out, "Bytes: 4")
	assertContains(t, out, "Bool: true")
	assertContains(t, out, "Timestamp: 2021-11-03 09:35:01 +0000 UTC")
}

func TestPgInterleavedTable(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgInterleavedTable, dbName, "failed to execute PG interleaved table sample")
	assertContains(t, out, "Created interleaved table hierarchy using PostgreSQL dialect")
}

func TestPgBatchDml(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(
		dbName,
		`CREATE TABLE Singers (
		   SingerId  bigint NOT NULL PRIMARY KEY,
		   FirstName varchar(1024),
		   LastName  varchar(1024)
		 )`)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgBatchDml, dbName, "failed to execute PG Batch DML sample")
	assertContains(t, out, "Inserted [1 1] singers")
}

func TestPgPartitionedDml(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(
		dbName,
		`CREATE TABLE users (
			user_id   bigint NOT NULL PRIMARY KEY,
			user_name varchar(1024),
			active    boolean
		)`)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	client, err := spanner.NewClient(context.Background(), dbName)
	if err != nil {
		t.Fatalf("failed to create Spanner client: %v", err)
	}
	defer client.Close()
	_, err = client.Apply(context.Background(), []*spanner.Mutation{
		spanner.InsertOrUpdateMap("users", map[string]interface{}{
			"user_id":   1,
			"user_name": "User 1",
			"active":    false,
		}),
		spanner.InsertOrUpdateMap("users", map[string]interface{}{
			"user_id":   2,
			"user_name": "User 2",
			"active":    false,
		}),
		spanner.InsertOrUpdateMap("users", map[string]interface{}{
			"user_id":   3,
			"user_name": "User 3",
			"active":    true,
		}),
	})
	if err != nil {
		t.Fatalf("failed to insert test records: %v", err)
	}

	out := runSample(t, pgPartitionedDml, dbName, "failed to execute PG Partitioned DML sample")
	assertContains(t, out, "Deleted at least 2 inactive users")
}

func TestPgCaseSensitivity(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgCaseSensitivity, dbName, "failed to execute PG case sensitivity sample")
	assertContains(t, out, "SingerId: 1, FirstName: Bruce, LastName: Allison")
	assertContains(t, out, "SingerId: 1, FullName: Bruce Allison")
}

func TestPgOrderNulls(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()
	dbCleanup, err := createTestPgDatabase(dbName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	out := runSample(t, pgOrderNulls, dbName, "failed to execute PG order nulls sample")
	assertContains(t, out, "Singers ORDER BY Name\n\tAlice\n\tBruce\n\t<null>")
	assertContains(t, out, "Singers ORDER BY Name DESC\n\t<null>\n\tBruce\n\tAlice")
	assertContains(t, out, "Singers ORDER BY Name NULLS FIRST\n\t<null>\n\tAlice\n\tBruce")
	assertContains(t, out, "Singers ORDER BY Name DESC NULLS LAST\n\tBruce\n\tAlice\n\t<null>")
}

func TestProtoColumnSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	statements := []string{
		`CREATE TABLE Singers (
				SingerId   INT64 NOT NULL,
				FirstName  STRING(1024),
				LastName   STRING(1024),
			) PRIMARY KEY (SingerId)`,
		`CREATE TABLE Albums (
				SingerId     INT64 NOT NULL,
				AlbumId      INT64 NOT NULL,
				AlbumTitle   STRING(MAX)
			) PRIMARY KEY (SingerId, AlbumId),
			INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
	}
	dbCleanup, err := createTestDatabase(dbName, statements...)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer dbCleanup()

	var out string
	runSample(t, write, dbName, "failed to insert data")
	runSampleWithContext(ctx, t, addProtoColumn, dbName, "failed to update db for proto columns")
	out = runSample(t, updateDataWithProtoColumn, dbName, "failed to update data with proto columns")
	assertContains(t, out, "Data updated\n")
	out = runSample(t, updateDataWithProtoColumnWithDml, dbName, "failed to update data with proto columns using dml")
	assertContains(t, out, "record(s) updated.")
	out = runSample(t, queryWithProtoParameter, dbName, "failed to query with proto parameter")
	assertContains(t, out, "2 singer_id:2")
}

func TestGraphSample(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	_, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	out := runSampleWithContext(
		ctx, t, createDatabaseWithPropertyGraph, dbName,
		"failed to create a Spanner database with a property graph")
	assertContains(t, out, fmt.Sprintf("Created database [%s]", dbName))

	out = runSample(t, insertGraphData, dbName, "")

	out = runSample(t, insertGraphDataWithDml, dbName, "")
	assertContains(t, out, "2 Account record(s) inserted")
	assertContains(t, out, "2 AccountTransferAccount record(s) inserted")

	out = runSample(t, queryGraphData, dbName, "")
	assertContains(t, out, "Dana Alex 500.000000 2020-10-04T16:55:05Z")
	assertContains(t, out, "Lee Dana 300.000000 2020-09-25T02:36:14Z")
	assertContains(t, out, "Alex Lee 300.000000 2020-08-29T15:28:58Z")
	assertContains(t, out, "Alex Lee 100.000000 2020-10-04T16:55:05Z")
	assertContains(t, out, "Dana Lee 200.000000 2020-10-17T03:59:40Z")

	out = runSample(t, queryGraphDataWithParameter, dbName, "")
	assertContains(t, out, "Dana Alex 500.000000 2020-10-04T16:55:05Z")

	out = runSample(t, updateGraphDataWithDml, dbName, "")
	assertContains(t, out, "1 Account record(s) updated.")
	assertContains(t, out, "1 AccountTransferAccount record(s) updated.")

	out = runSample(t, updateGraphDataWithGraphQueryInDml, dbName, "")
	assertContains(t, out, "2 Account record(s) updated.")

	out = runSample(t, deleteGraphDataWithDml, dbName, "")
	assertContains(t, out, "1 AccountTransferAccount record(s) deleted.")
	assertContains(t, out, "1 Account record(s) deleted.")

	out = runSample(t, deleteGraphData, dbName, "")
}

func maybeCreateKey(projectId, locationId, keyRingId, keyId string) error {
	client, err := kms.NewKeyManagementClient(context.Background())
	if err != nil {
		return err
	}

	// Try to create a key ring
	createKeyRingRequest := kmspb.CreateKeyRingRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s", projectId, locationId),
		KeyRingId: keyRingId,
		KeyRing:   &kmspb.KeyRing{},
	}
	_, err = client.CreateKeyRing(context.Background(), &createKeyRingRequest)
	if err != nil {
		if status, ok := status.FromError(err); !ok || status.Code() != codes.AlreadyExists {
			return err
		}
	}

	// Try to create a key
	createKeyRequest := kmspb.CreateCryptoKeyRequest{
		Parent:      fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", projectId, locationId, keyRingId),
		CryptoKeyId: keyId,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
		},
	}
	_, err = client.CreateCryptoKey(context.Background(), &createKeyRequest)
	if err != nil {
		if status, ok := status.FromError(err); !ok || status.Code() != codes.AlreadyExists {
			return err
		}
	}

	return nil
}

func runSample(t *testing.T, f sampleFunc, dbName, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, dbName); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runSampleWithContext(ctx context.Context, t *testing.T, f sampleFuncWithContext, dbName, errMsg string) string {
	var b bytes.Buffer
	if err := f(ctx, &b, dbName); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runCreateBackupSample(ctx context.Context, t *testing.T, f createBackupSampleFunc, dbName string, backupID string, versionTime time.Time, errMsg string) string {
	var b bytes.Buffer
	if err := f(ctx, &b, dbName, backupID, versionTime); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runBackupSample(ctx context.Context, t *testing.T, f backupSampleFunc, dbName, backupID, errMsg string) string {
	var b bytes.Buffer
	if err := f(ctx, &b, dbName, backupID); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runBackupSampleWithoutContext(t *testing.T, f backupSampleFuncWithoutContext, dbName, backupID, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, dbName, backupID); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runBackupSampleWithRetry(ctx context.Context, t *testing.T, f backupSampleFunc, dbName, backupID, errMsg string, maxAttempts int) string {
	var b bytes.Buffer
	testutil.Retry(t, maxAttempts, time.Minute, func(r *testutil.R) {
		b.Reset()
		if err := f(ctx, &b, dbName, backupID); err != nil {
			if strings.Contains(err.Error(), "Please retry the operation once the pending restores complete") {
				r.Errorf("%s: %v", errMsg, err)
			} else {
				t.Fatalf("%s: %v", errMsg, err)
			}
		}
	})
	return b.String()
}

func runBackupScheduleSample(t *testing.T, f backupScheduleSampleFunc, dbName string, scheduleId string, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, dbName, scheduleId); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runInstancePartitionSample(t *testing.T, f instancePartitionSampleFunc, projectID, instanceID, instancePartitionID, errMsg string) string {
	var b bytes.Buffer
	if err := f(&b, projectID, instanceID, instancePartitionID); err != nil {
		t.Errorf("%s: %v", errMsg, err)
	}
	return b.String()
}

func runInstanceSample(t *testing.T, f instanceSampleFunc, projectID, instanceID, errMsg string) string {
	var b bytes.Buffer

	testutil.Retry(t, 20, time.Minute, func(r *testutil.R) {
		b.Reset()
		if err := f(&b, projectID, instanceID); err != nil {
			// Retry if the instance could not be created because there have
			// been too many create requests in the past minute.
			if spanner.ErrCode(err) == codes.ResourceExhausted && strings.Contains(err.Error(), "Quota exceeded for quota metric 'Instance create requests'") {
				r.Errorf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
				return
			} else {
				t.Fatalf("%s: %v", errMsg, err)
			}
		}
	})
	return b.String()
}

func mustRunSample(t *testing.T, f sampleFuncWithContext, dbName, errMsg string) string {
	var b bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := f(ctx, &b, dbName); err != nil {
		t.Fatalf("%s: %v", errMsg, err)
	}
	return b.String()
}

func createTestInstance(t *testing.T, projectID string, instanceConfigName string) (instanceID string, cleanup func()) {
	ctx := context.Background()
	instanceID = fmt.Sprintf("go-sample-%s", uuid.New().String()[:16])
	instanceName := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		t.Fatalf("failed to create InstanceAdminClient: %v", err)
	}
	databaseAdmin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		t.Fatalf("failed to create DatabaseAdminClient: %v", err)
	}

	// Cleanup old test instances that might not have been deleted.
	iter := instanceAdmin.ListInstances(ctx, &instancepb.ListInstancesRequest{
		Parent: fmt.Sprintf("projects/%v", projectID),
		Filter: "labels.cloud_spanner_samples_test:true",
	})
	for {
		instance, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("failed to list existing instances: %v", err)
		}
		if createTimeString, ok := instance.Labels["create_time"]; ok {
			seconds, err := strconv.ParseInt(createTimeString, 10, 64)
			if err != nil {
				t.Logf("could not parse create time %v: %v", createTimeString, err)
				continue
			}
			createTime := time.Unix(seconds, 0)
			diff := time.Now().Sub(createTime)
			if diff > time.Hour*24 {
				t.Logf("deleting stale test instance %v", instance.Name)
				deleteInstanceAndBackups(t, instance.Name, instanceAdmin, databaseAdmin)
			}
		}
	}

	instanceConfigName = fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, instanceConfigName)

	testutil.Retry(t, 20, time.Minute, func(r *testutil.R) {
		op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
			Parent:     fmt.Sprintf("projects/%s", projectID),
			InstanceId: instanceID,
			Instance: &instancepb.Instance{
				Config:      instanceConfigName,
				DisplayName: instanceID,
				NodeCount:   1,
				Labels: map[string]string{
					"cloud_spanner_samples_test": "true",
					"create_time":                fmt.Sprintf("%v", time.Now().Unix()),
				},
				Edition: instancepb.Instance_ENTERPRISE_PLUS,
			},
		})
		if err != nil {
			// Retry if the instance could not be created because there have
			// been too many create requests in the past minute.
			if spanner.ErrCode(err) == codes.ResourceExhausted && strings.Contains(err.Error(), "Quota exceeded for quota metric 'Instance create requests'") {
				r.Errorf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
				return
			} else {
				t.Fatalf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
			}
		}
		_, err = op.Wait(ctx)
		if err != nil {
			t.Fatalf("waiting for instance creation to finish failed: %v", err)
		}
	})

	return instanceID, func() {
		deleteInstanceAndBackups(t, instanceName, instanceAdmin, databaseAdmin)
		instanceAdmin.Close()
		databaseAdmin.Close()
	}
}

func createTestPgDatabase(db string, extraStatements ...string) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	m := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if m == nil || len(m) != 3 {
		return func() {}, fmt.Errorf("invalid database id %s", db)
	}

	client, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return func() {}, err
	}
	defer client.Close()

	opCreate, err := client.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          m[1],
		DatabaseDialect: adminpb.DatabaseDialect_POSTGRESQL,
		CreateStatement: `CREATE DATABASE "` + m[2] + `"`,
	})
	if err != nil {
		return func() {}, err
	}
	if _, err := opCreate.Wait(ctx); err != nil {
		return func() {}, err
	}
	dropDb := func() {
		client, err := database.NewDatabaseAdminClient(ctx)
		if err != nil {
			return
		}
		defer client.Close()
		client.DropDatabase(context.Background(), &adminpb.DropDatabaseRequest{
			Database: db,
		})
	}
	if len(extraStatements) > 0 {
		opUpdate, err := client.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
			Database:   db,
			Statements: extraStatements,
		})
		if err != nil {
			return dropDb, err
		}
		if err := opUpdate.Wait(ctx); err != nil {
			return dropDb, err
		}
	}
	return dropDb, nil
}

func createTestDatabase(db string, extraStatements ...string) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	m := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if m == nil || len(m) != 3 {
		return func() {}, fmt.Errorf("invalid database id %s", db)
	}

	client, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return func() {}, err
	}
	defer client.Close()

	opCreate, err := client.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          m[1],
		CreateStatement: "CREATE DATABASE `" + m[2] + "`",
		ExtraStatements: extraStatements,
	})
	if err != nil {
		return func() {}, err
	}
	if _, err := opCreate.Wait(ctx); err != nil {
		return func() {}, err
	}
	dropDb := func() {
		client, err := database.NewDatabaseAdminClient(ctx)
		if err != nil {
			return
		}
		defer client.Close()
		client.DropDatabase(context.Background(), &adminpb.DropDatabaseRequest{
			Database: db,
		})
	}
	return dropDb, nil
}

func deleteInstanceAndBackups(
	t *testing.T,
	instanceName string,
	instanceAdmin *instance.InstanceAdminClient,
	databaseAdmin *database.DatabaseAdminClient) {
	ctx := context.Background()
	// Delete all backups before deleting the instance.
	iter := databaseAdmin.ListBackups(ctx, &adminpb.ListBackupsRequest{
		Parent: instanceName,
	})
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Failed to list backups for instance %s: %v", instanceName, err)
		}
		databaseAdmin.DeleteBackup(ctx, &adminpb.DeleteBackupRequest{Name: resp.Name})
	}
	instanceAdmin.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: instanceName})
}

func getSampleProjectId(t *testing.T) string {
	// These tests get the project id from the environment variable
	// GOLANG_SAMPLES_SPANNER that is also used by other integration tests for
	// Spanner samples. The tests in this file create a separate instance for
	// each test, so only the project id is used, and the rest of the instance
	// name is ignored.
	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("Skipping spanner integration test. Set GOLANG_SAMPLES_SPANNER.")
	}
	if !strings.HasPrefix(instance, "projects/") {
		t.Fatal("Spanner instance ref must be in the form of 'projects/PROJECT_ID/instances/INSTANCE_ID'")
	}
	projectId, _, err := parseInstanceName(instance)
	if err != nil {
		t.Fatalf("Could not parse project id from instance name %q: %v", instance, err)
	}
	return projectId
}

// getSamplesInstanceConfig specifies the instance config used to create an instance for testing.
// It can be changed by setting the environment variable
// GOLANG_SAMPLES_SPANNER_INSTANCE_CONFIG.
func getSamplesInstanceConfig() string {
	return os.Getenv("GOLANG_SAMPLES_SPANNER_INSTANCE_CONFIG")
}

func assertContains(t *testing.T, out string, sub string) {
	t.Helper()
	if !strings.Contains(out, sub) {
		t.Errorf("got output %q; want it to contain %q", out, sub)
	}
}

func assertNotContains(t *testing.T, out string, sub string) {
	t.Helper()
	if strings.Contains(out, sub) {
		t.Errorf("got output %q; want it to not contain %q", out, sub)
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
	return cleanupInstanceWithName(fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID))
}

func cleanupInstanceWithName(instanceName string) error {
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot create instance databaseAdmin client: %w", err)
	}
	defer instanceAdmin.Close()

	if err := instanceAdmin.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: instanceName}); err != nil {
		return fmt.Errorf("failed to delete instance %s (error %w), might need a manual removal",
			instanceName, err)
	}
	return nil
}

func cleanupInstanceConfigs(projectID string) error {
	// Delete all custom instance configurations.
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot create instance admin client: %w", err)
	}
	defer instanceAdmin.Close()
	configIter := instanceAdmin.ListInstanceConfigs(ctx, &instancepb.ListInstanceConfigsRequest{
		Parent: "projects/" + projectID,
	})
	for {
		resp, err := configIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if strings.Contains(resp.Name, "custom-golang-samples") {
			instanceAdmin.DeleteInstanceConfig(ctx, &instancepb.DeleteInstanceConfigRequest{Name: resp.Name})
		}
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
