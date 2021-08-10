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

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sampleFunc func(w io.Writer, dbName string) error
type sampleFuncWithContext func(ctx context.Context, w io.Writer, dbName string) error
type instanceSampleFunc func(w io.Writer, projectID, instanceID string) error
type backupSampleFunc func(ctx context.Context, w io.Writer, dbName, backupID string) error
type createBackupSampleFunc func(ctx context.Context, w io.Writer, dbName, backupID string, versionTime time.Time) error

var (
	validInstancePattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)$")
)

func initTest(t *testing.T, id string) (instName, dbName string, cleanup func()) {
	projectID := getSampleProjectId(t)
	instName, cleanup = createTestInstance(t, projectID, "regional-us-central1")
	dbID := validLength(fmt.Sprintf("smpl-%s", id), t)
	dbName = fmt.Sprintf("%s/databases/%s", instName, dbID)

	return
}

func initTestWithConfig(t *testing.T, id string, instanceConfigName string) (instName, dbName string, cleanup func()) {
	projectID := getSampleProjectId(t)
	instName, cleanup = createTestInstance(t, projectID, instanceConfigName)
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

func TestCreateInstances(t *testing.T) {
	_ = testutil.SystemTest(t)
	t.Parallel()

	runCreateInstanceSample(t, createInstance)
	runCreateInstanceSample(t, createInstanceWithProcessingUnits)
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

	out = runSample(t, read, dbName, "failed to read data")
	assertContains(t, out, "1 1 Total Junk")
	out = runSample(t, query, dbName, "failed to query data")
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

	out = runSample(t, readBatchData, dbName, "failed to read batch data")
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
}

func TestBackupSample(t *testing.T) {
	t.Skip("https://github.com/GoogleCloudPlatform/golang-samples/issues/2143")
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

	out = runSampleWithContext(ctx, t, listBackupOperations, dbName, "failed to list backup operations")
	assertContains(t, out, fmt.Sprintf("on database %s", dbName))

	out = runBackupSample(ctx, t, updateBackup, dbName, backupID, "failed to update a backup")
	assertContains(t, out, fmt.Sprintf("Updated backup %s", backupID))

	out = runBackupSampleWithRetry(ctx, t, restoreBackup, restoreDBName, backupID, "failed to restore a backup", 10)
	assertContains(t, out, fmt.Sprintf("Source database %s restored from backup", dbName))

	// This sample should run after a restore operation.
	out = runSampleWithContext(ctx, t, listDatabaseOperations, restoreDBName, "failed to list database operations")
	assertContains(t, out, fmt.Sprintf("Database %s restored from backup", restoreDBName))

	out = runBackupSample(ctx, t, deleteBackup, dbName, backupID, "failed to delete a backup")
	assertContains(t, out, fmt.Sprintf("Deleted backup %s", backupID))
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

	instName, dbName, cleanup := initTest(t, randomID())
	defer cleanup()

	var b bytes.Buffer

	locationId := "us-central1"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// Create an encrypted database. The database is automatically deleted by the cleanup function.
	if err := createDatabaseWithCustomerManagedEncryptionKey(ctx, &b, dbName, kmsKeyName); err != nil {
		t.Errorf("failed to create database with customer managed encryption key: %v", err)
	}
	out := b.String()
	assertContains(t, out, fmt.Sprintf("Created database [%s] using encryption key %q", dbName, kmsKeyName))

	// Try to create a backup of the encrypted database and delete it after the test.
	backupId := fmt.Sprintf("enc-backup-%s", randomID())
	b.Reset()
	if err := createBackupWithCustomerManagedEncryptionKey(ctx, &b, dbName, backupId, kmsKeyName); err != nil {
		t.Errorf("failed to create backup with customer managed encryption key: %v", err)
	}
	out = b.String()
	assertContains(t, out, fmt.Sprintf("backups/%s", backupId))
	assertContains(t, out, fmt.Sprintf("using encryption key %s", kmsKeyName))

	// Try to restore the encrypted database and delete the restored database after the test.
	restoredName := fmt.Sprintf("%s/databases/rest-enc-%s", instName, randomID())
	restoreFunc := func(ctx context.Context, w io.Writer, dbName, backupID string) error {
		return restoreBackupWithCustomerManagedEncryptionKey(ctx, w, dbName, backupId, kmsKeyName)
	}
	out = runBackupSampleWithRetry(ctx, t, restoreFunc, restoredName, backupId, "failed to restore database with customer managed encryption key", 10)
	assertContains(t, out, fmt.Sprintf("Database %s restored", dbName))
	assertContains(t, out, fmt.Sprintf("using encryption key %s", kmsKeyName))
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

func createTestInstance(t *testing.T, projectID string, instanceConfigName string) (instanceName string, cleanup func()) {
	ctx := context.Background()
	instanceID := fmt.Sprintf("go-sample-%s", uuid.New().String()[:16])
	instanceName = fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
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

	return instanceName, func() {
		deleteInstanceAndBackups(t, instanceName, instanceAdmin, databaseAdmin)
		instanceAdmin.Close()
		databaseAdmin.Close()
	}
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
	return cleanupInstanceWithName(fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID))
}

func cleanupInstanceWithName(instanceName string) error {
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot create instance databaseAdmin client: %v", err)
	}
	defer instanceAdmin.Close()

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
