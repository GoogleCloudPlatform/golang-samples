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

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type runCommandFunc func(t *testing.T, cmd, dbName string) string
type runBackupCommandFunc func(t *testing.T, cmd, dbName, backupID string) string

func initTest(t *testing.T, projectID string) (dbName string, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client, runCommand runCommandFunc, mustRunCommand runCommandFunc, cleanup func()) {
	instance := getInstance(t)
	databaseID := validLength(fmt.Sprintf("test-%s", projectID), t)
	dbName = fmt.Sprintf("%s/databases/%s", instance, databaseID)

	ctx := context.Background()
	adminClient, dataClient = createClients(ctx, dbName)

	// Check for database existance prior to test start and delete, as resources may not have
	// been cleaned up from previous invocations.
	if db, err := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbName}); err == nil {
		t.Logf("database %s exists in state %s. delete result: %v", db.GetName(), db.GetState().String(),
			adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName}))
	}

	runCommand = func(t *testing.T, cmd, dbName string) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName, ""); err != nil {
			t.Errorf("run(%q, %q): %v", cmd, dbName, err)
		}
		return b.String()
	}
	mustRunCommand = func(t *testing.T, cmd, dbName string) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName, ""); err != nil {
			t.Fatalf("run(%q, %q): %v", cmd, dbName, err)
		}
		return b.String()
	}
	cleanup = func() {
		dataClient.Close()
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

func initBackupTest(t *testing.T, projectID, dbName string, adminClient *database.DatabaseAdminClient, dataClient *spanner.Client) (restoreDBName, backupID, cancelledBackupID string, runBackupCommand runBackupCommandFunc, cleanup func()) {
	instance := getInstance(t)
	restoreDatabaseID := validLength(fmt.Sprintf("restore-%s", projectID), t)
	restoreDBName = fmt.Sprintf("%s/databases/%s", instance, restoreDatabaseID)
	backupID = validLength(fmt.Sprintf("backup-%s", projectID), t)
	cancelledBackupID = validLength(fmt.Sprintf("cancel-%s", projectID), t)

	ctx := context.Background()
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

	runBackupCommand = func(t *testing.T, cmd, dbName, backupID string) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName, backupID); err != nil {
			t.Errorf("run(%q, %q): %v", cmd, dbName, err)
		}
		return b.String()
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

func TestSample(t *testing.T) {
	tc := testutil.SystemTest(t)

	dbName, _, _, runCommand, mustRunCommand, cleanup := initTest(t, tc.ProjectID)
	defer cleanup()

	// We execute all the commands of the tutorial code. These commands have to be run in a specific
	// order since in many cases earlier commands setup the database for the subsequent commands.
	mustRunCommand(t, "createdatabase", dbName)
	runCommand(t, "write", dbName)
	runCommand(t, "addnewcolumn", dbName)

	runCommand(t, "delete", dbName)
	runCommand(t, "write", dbName)
	runCommand(t, "update", dbName)
	out := runCommand(t, "dmlwritetxn", dbName)
	assertContains(t, out, "Moved 200000 from Album2's MarketingBudget to Album1")
	out = runCommand(t, "querynewcolumn", dbName)
	assertContains(t, out, "1 1 300000")
	assertContains(t, out, "2 2 300000")

	runCommand(t, "delete", dbName)
	runCommand(t, "write", dbName)
	runCommand(t, "update", dbName)
	out = runCommand(t, "writetransaction", dbName)
	assertContains(t, out, "Moved 200000 from Album2's MarketingBudget to Album1")
	out = runCommand(t, "querynewcolumn", dbName)
	assertContains(t, out, "1 1 300000")
	assertContains(t, out, "2 2 300000")

	runCommand(t, "delete", dbName)
	runCommand(t, "write", dbName)
	writeTime := time.Now()

	assertContains(t, runCommand(t, "read", dbName), "1 1 Total Junk")

	assertContains(t, runCommand(t, "query", dbName), "1 1 Total Junk")

	runCommand(t, "addindex", dbName)
	out = runCommand(t, "queryindex", dbName)
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	if strings.Contains(out, "Green") {
		t.Errorf("got output %q; should not contain Green", out)
	}

	out = runCommand(t, "readindex", dbName)
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	assertContains(t, out, "Green")

	runCommand(t, "delete", dbName)
	runCommand(t, "write", dbName)
	runCommand(t, "update", dbName)
	runCommand(t, "addstoringindex", dbName)
	assertContains(t, runCommand(t, "readstoringindex", dbName), "500000")
	out = runCommand(t, "readonlytransaction", dbName)
	if strings.Count(out, "Total Junk") != 2 {
		t.Errorf("got output %q; wanted it to contain 2 occurrences of Total Junk", out)
	}

	// Wait at least 15 seconds since the write.
	time.Sleep(time.Until(writeTime.Add(16 * time.Second)))
	out = runCommand(t, "readstaledata", dbName)
	assertContains(t, out, "Go, Go, Go")
	assertContains(t, out, "Forever Hold Your Peace")
	assertContains(t, out, "Green")

	assertContains(t, runCommand(t, "readbatchdata", dbName), "1 Marc Richards")

	runCommand(t, "addcommittimestamp", dbName)
	runCommand(t, "updatewithtimestamp", dbName)
	out = runCommand(t, "querywithtimestamp", dbName)
	assertContains(t, out, "1000000")

	runCommand(t, "writestructdata", dbName)
	assertContains(t, runCommand(t, "querywithstruct", dbName), "6")
	out = runCommand(t, "querywitharrayofstruct", dbName)
	assertContains(t, out, "6")
	assertContains(t, out, "7")
	assertContains(t, out, "8")
	assertContains(t, runCommand(t, "querywithstructfield", dbName), "6")
	out = runCommand(t, "querywithnestedstructfield", dbName)
	assertContains(t, out, "6 Imagination")
	assertContains(t, out, "9 Imagination")

	runCommand(t, "createtabledocswithtimestamp", dbName)
	runCommand(t, "writetodocstable", dbName)
	runCommand(t, "updatedocstable", dbName)

	assertContains(t, runCommand(t, "querydocstable", dbName), "Hello World 1 Updated")

	runCommand(t, "createtabledocswithhistorytable", dbName)
	runCommand(t, "writewithhistory", dbName)
	runCommand(t, "updatewithhistory", dbName)

	out = runCommand(t, "querywithhistory", dbName)
	assertContains(t, out, "1 1 Hello World 1 Updated")

	out = runCommand(t, "dmlinsert", dbName)
	assertContains(t, out, "record(s) inserted")

	out = runCommand(t, "dmlupdate", dbName)
	assertContains(t, out, "record(s) updated")

	out = runCommand(t, "dmldelete", dbName)
	assertContains(t, out, "record(s) deleted")

	out = runCommand(t, "dmlwithtimestamp", dbName)
	assertContains(t, out, "record(s) updated")

	out = runCommand(t, "dmlwriteread", dbName)
	assertContains(t, out, "Found record name with ")

	out = runCommand(t, "dmlupdatestruct", dbName)
	assertContains(t, out, "record(s) inserted")

	out = runCommand(t, "dmlwrite", dbName)
	assertContains(t, out, "record(s) inserted")

	out = runCommand(t, "querywithparameter", dbName)
	assertContains(t, out, "12 Melissa Garcia")

	out = runCommand(t, "dmlupdatepart", dbName)
	assertContains(t, out, "record(s) updated")

	out = runCommand(t, "dmldeletepart", dbName)
	assertContains(t, out, "record(s) deleted")

	out = runCommand(t, "dmlbatchupdate", dbName)
	assertContains(t, out, "Executed 2 SQL statements using Batch DML.")

	out = runCommand(t, "createtablewithdatatypes", dbName)
	assertContains(t, out, "Created Venues table")

	out = runCommand(t, "writedatatypesdata", dbName)
	out = runCommand(t, "querywitharray", dbName)
	assertContains(t, out, "19 Venue 19 2020-11-01")
	assertContains(t, out, "42 Venue 42 2020-10-01")

	out = runCommand(t, "querywithbool", dbName)
	assertContains(t, out, "19 Venue 19 true")

	out = runCommand(t, "querywithbytes", dbName)
	assertContains(t, out, "4 Venue 4")

	out = runCommand(t, "querywithdate", dbName)
	assertContains(t, out, "4 Venue 4 2018-09-02")
	assertContains(t, out, "42 Venue 42 2018-10-01")

	out = runCommand(t, "querywithfloat", dbName)
	assertContains(t, out, "4 Venue 4 0.8")
	assertContains(t, out, "19 Venue 19 0.9")

	out = runCommand(t, "querywithint", dbName)
	assertContains(t, out, "19 Venue 19 6300")
	assertContains(t, out, "42 Venue 42 3000")

	out = runCommand(t, "querywithstring", dbName)
	assertContains(t, out, "42 Venue 42")

	// Wait 5 seconds to avoid a time drift issue for the next query:
	// https://github.com/GoogleCloudPlatform/golang-samples/issues/1146.
	time.Sleep(time.Second * 5)
	out = runCommand(t, "querywithtimestampparameter", dbName)
	assertContains(t, out, "4 Venue 4")
	assertContains(t, out, "19 Venue 19")
	assertContains(t, out, "42 Venue 42")
	out = runCommand(t, "querywithqueryoptions", dbName)
	assertContains(t, out, "4 Venue 4")
	assertContains(t, out, "19 Venue 19")
	assertContains(t, out, "42 Venue 42")
	out = runCommand(t, "createclientwithqueryoptions", dbName)
	assertContains(t, out, "4 Venue 4")
	assertContains(t, out, "19 Venue 19")
	assertContains(t, out, "42 Venue 42")
}

func TestBackupSample(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	dbName, adminClient, dataClient, runCommand, mustRunCommand, cleanup := initTest(t, tc.ProjectID)
	defer cleanup()
	restoreDBName, backupID, cancelledBackupID, runBackupCommand, cleanupBackup := initBackupTest(t, tc.ProjectID, dbName, adminClient, dataClient)

	var out string
	// Set up the database for testing backup operations.
	mustRunCommand(t, "createdatabase", dbName)
	runCommand(t, "write", dbName)
	runCommand(t, "writestructdata", dbName)
	runCommand(t, "createtabledocswithtimestamp", dbName)
	runCommand(t, "writetodocstable", dbName)
	runCommand(t, "createtabledocswithhistorytable", dbName)
	runCommand(t, "writewithhistory", dbName)
	runCommand(t, "createtablewithdatatypes", dbName)
	runCommand(t, "writedatatypesdata", dbName)

	// Start testing backup operations.
	out = runBackupCommand(t, "createbackup", dbName, backupID)
	assertContains(t, out, fmt.Sprintf("backups/%s", backupID))

	out = runBackupCommand(t, "cancelbackup", dbName, cancelledBackupID)
	assertContains(t, out, "Backup cancelled.")

	out = runBackupCommand(t, "listbackups", dbName, backupID)
	assertContains(t, out, fmt.Sprintf("/backups/%s", backupID))
	assertContains(t, out, "Backups listed.")

	out = runCommand(t, "listbackupoperations", dbName)
	assertContains(t, out, fmt.Sprintf("on database %s", dbName))

	out = runBackupCommand(t, "updatebackup", dbName, backupID)
	assertContains(t, out, fmt.Sprintf("Updated backup %s", backupID))

	out = runBackupCommand(t, "restorebackup", restoreDBName, backupID)
	assertContains(t, out, fmt.Sprintf("Source database %s restored from backup", dbName))

	// This command should run after a restore operation.
	out = runCommand(t, "listdatabaseoperations", restoreDBName)
	assertContains(t, out, fmt.Sprintf("Database %s restored from backup", restoreDBName))

	// Delete the restore DB.
	cleanupBackup()

	out = runBackupCommand(t, "deletebackup", dbName, backupID)
	assertContains(t, out, fmt.Sprintf("Deleted backup %s", backupID))
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
