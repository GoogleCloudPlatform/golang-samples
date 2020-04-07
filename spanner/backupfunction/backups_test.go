// Copyright 2019 Google LLC
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

// [START spanner_functions_backup_util]

// Package backupfunction is a Cloud function that can be periodically triggered to create a backup
// for the specified database.
package backupfunction

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/google/uuid"
	// apioption "google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

var (
	// testProjectID specifies the project used for testing. It can be changed
	// by setting environment variable GCLOUD_TESTS_GOLANG_PROJECT_ID.
	testProjectID    = projID()
	testInstanceName = os.Getenv("GCLOUD_TESTS_GOLANG_INSTANCE_NAME")
	testEndpoint     = os.Getenv("GCLOUD_TESTS_GOLANG_ENDPOINT")

	databaseAdmin *database.DatabaseAdminClient
	instanceAdmin *instance.InstanceAdminClient
)

const (
	envProjID     = "GCLOUD_TESTS_GOLANG_PROJECT_ID"
	envPrivateKey = "GCLOUD_TESTS_GOLANG_KEY"
)

// ProjID returns the project ID to use in integration tests, or the empty
// string if none is configured.
func projID() string {
	return os.Getenv(envProjID)
}

func initIntegrationTests(t *testing.T) (cleanup func()) {
	ctx := context.Background()
	flag.Parse() // Needed for testing.Short().

	if testing.Short() {
		t.Log("Integration tests skipped in -short mode.")
		return func() {}
	}

	if testProjectID == "" {
		t.Log("Integration tests skipped: GCLOUD_TESTS_GOLANG_PROJECT_ID is missing")
		return func() {}
	}

	var err error

	// Check if a specific endpoint is set for the integration test
	// var opts apioption.ClientOption
	// if testEndpoint != "" {
	// 	t.Logf("Running integration test with endpoint %s", testEndpoint)
	// 	opts = apioption.WithEndpoint(testEndpoint)
	// }

	// Create InstanceAdmin and DatabaseAdmin clients.
	instanceAdmin, err = instance.NewInstanceAdminClient(ctx) //, opts)
	if err != nil {
		t.Fatalf("cannot create instance databaseAdmin client: %v", err)
	}
	databaseAdmin, err = database.NewDatabaseAdminClient(ctx) //, opts)
	if err != nil {
		t.Fatalf("cannot create databaseAdmin client: %v", err)
	}

	// If a specific instance was selected for testing, use that.  Otherwise create a new instance for testing and
	// tear it down after the test.
	createInstanceForTest := testInstanceName == ""
	if createInstanceForTest {
		testInstanceName = fmt.Sprintf("go-test-%s", uuid.New())
		// limit testInstanceName to length  of 40
		testInstanceName = testInstanceName[0:39]
		// Get the list of supported instance configs for the project that is used
		// for the integration tests. The supported instance configs can differ per
		// project. The integration tests will use the first instance config that
		// is returned by Cloud Spanner. This will normally be the regional config
		// that is physically the closest to where the request is coming from.
		configIterator := instanceAdmin.ListInstanceConfigs(ctx, &instancepb.ListInstanceConfigsRequest{
			Parent: fmt.Sprintf("projects/%s", testProjectID),
		})
		config, err := configIterator.Next()
		if err != nil {
			t.Fatalf("Cannot get any instance configurations.\nPlease make sure the Cloud Spanner API is enabled for the test project.\nGet error: %v", err)
		}

		// Create a test instance to use for this test run.
		op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
			Parent:     fmt.Sprintf("projects/%s", testProjectID),
			InstanceId: testInstanceName,
			Instance: &instancepb.Instance{
				Config:      config.Name,
				DisplayName: testInstanceName,
				NodeCount:   1,
			},
		})
		if err != nil {
			t.Fatalf("could not create instance with id %s: %v", fmt.Sprintf("projects/%s/instances/%s", testProjectID, testInstanceName), err)
		}
		// Wait for the instance creation to finish.
		i, err := op.Wait(ctx)
		if err != nil {
			t.Fatalf("waiting for instance creation to finish failed: %v", err)
		}
		if i.State != instancepb.Instance_READY {
			t.Logf("instance state is not READY, it might be that the test instance will cause problems during tests. Got state %v\n", i.State)
		}
	}

	return func() {
		if createInstanceForTest {
			if err := instanceAdmin.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: testInstanceName}); err != nil {
				t.Logf("failed to drop instance %s (error %v), might need a manual removal",
					testInstanceName, err)
			}
			// Delete other test instances that may be lingering around.
			cleanupInstances(t, testInstanceName)
		}

		databaseAdmin.Close()
		instanceAdmin.Close()
	}
}

// Prepare initializes Cloud Spanner testing DB and clients.
func prepareIntegrationTest(ctx context.Context, t *testing.T) (string, func()) {
	if databaseAdmin == nil {
		t.Skip("Integration tests skipped")
	}
	// Construct a unique test DB name.
	dbName := fmt.Sprintf("test%s", uuid.New())
	// limit dbName to length  of 30
	dbName = dbName[0:10]
	dbPath := fmt.Sprintf("projects/%v/instances/%v/databases/%v", testProjectID, testInstanceName, dbName)
	/// Create database and tables.
	op, err := databaseAdmin.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%v/instances/%v", testProjectID, testInstanceName),
		CreateStatement: "CREATE DATABASE " + dbName,
		ExtraStatements: []string{
			`CREATE TABLE Singers (
				SingerId	INT64 NOT NULL,
				FirstName	STRING(1024),
				LastName	STRING(1024),
				SingerInfo	BYTES(MAX)
			) PRIMARY KEY (SingerId)`,
			`CREATE INDEX SingerByName ON Singers(FirstName, LastName)`,
			`CREATE TABLE Accounts (
				AccountId	INT64 NOT NULL,
				Nickname	STRING(100),
				Balance		INT64 NOT NULL,
			) PRIMARY KEY (AccountId)`,
			`CREATE INDEX AccountByNickname ON Accounts(Nickname) STORING (Balance)`,
			`CREATE TABLE Types (
				RowID		INT64 NOT NULL,
				String		STRING(MAX),
				StringArray	ARRAY<STRING(MAX)>,
				Bytes		BYTES(MAX),
				BytesArray	ARRAY<BYTES(MAX)>,
				Int64a		INT64,
				Int64Array	ARRAY<INT64>,
				Bool		BOOL,
				BoolArray	ARRAY<BOOL>,
				Float64		FLOAT64,
				Float64Array	ARRAY<FLOAT64>,
				Date		DATE,
				DateArray	ARRAY<DATE>,
				Timestamp	TIMESTAMP,
				TimestampArray	ARRAY<TIMESTAMP>,
			) PRIMARY KEY (RowID)`,
		},
	})
	if err != nil {
		t.Fatalf("cannot create testing DB %v: %v", dbPath, err)
	}
	if _, err := op.Wait(ctx); err != nil {
		t.Fatalf("cannot create testing DB %v: %v", dbPath, err)
	}

	return dbPath, func() {
		err := databaseAdmin.DropDatabase(ctx, &adminpb.DropDatabaseRequest{
			Database: dbPath,
		})
		if err != nil {
			t.Fatalf("cannot drop testing DB %v: %v", dbPath, err)
		}
	}
}

func cleanupInstances(t *testing.T, instanceName string) {
	ctx := context.Background()
	if testInstanceName == "" {
		t.Logf("Deleting instance %s", instanceName)

		if err := instanceAdmin.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: instanceName}); err != nil {
			t.Logf("failed to delete instance %s (error %v), might need a manual removal",
				instanceName, err)
		}
	}
}

func TestIntegrationCreateBackup(t *testing.T) {
	ctx := context.Background()
	instanceCleanup := initIntegrationTests(t)
	defer instanceCleanup()
	if databaseAdmin == nil {
		t.Skip("Integration tests skipped")
	}
	testDatabaseName, cleanup := prepareIntegrationTest(ctx, t)
	defer cleanup()

	backupPrefix := "go-test-"
	expires := time.Duration(time.Hour * 7)
	op, err := CreateBackup(ctx, os.Stdout, databaseAdmin, testDatabaseName, expires, backupPrefix)
	if err != nil {
		t.Fatal(err)
	}
	backup, completionErr := op.Wait(ctx)
	if completionErr != nil {
		t.Logf("Error completing backup: %v", completionErr)
	}

	defer func() {
		err := databaseAdmin.DeleteBackup(ctx, &adminpb.DeleteBackupRequest{Name: backup.Name})
		if err != nil {
			t.Fatal(err)
		}
	}()
}
