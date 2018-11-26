// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSample(t *testing.T) {
	tc := testutil.SystemTest(t)

	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("Skipping spanner integration test. Set GOLANG_SAMPLES_SPANNER.")
	}
	if !strings.HasPrefix(instance, "projects/") {
		t.Fatal("Spanner instance ref must be in the form of 'projects/PROJECT_ID/instances/INSTANCE_ID'")
	}
	dbName := fmt.Sprintf("%s/databases/test-%s", instance, tc.ProjectID)

	ctx := context.Background()
	adminClient, dataClient := createClients(ctx, dbName)
	defer adminClient.Close()
	defer dataClient.Close()

	assertContains := func(out string, sub string) {
		if !strings.Contains(out, sub) {
			t.Errorf("got output %q; want it to contain %q", out, sub)
		}
	}
	runCommand := func(t *testing.T, cmd string, dbName string) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName); err != nil {
			t.Errorf("run(%q, %q): %v", cmd, dbName, err)
		}
		return b.String()
	}
	mustRunCommand := func(t *testing.T, cmd string, dbName string) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName); err != nil {
			t.Fatalf("run(%q, %q): %v", cmd, dbName, err)
		}
		return b.String()
	}

	defer func() {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			err := adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName})
			if err != nil {
				r.Errorf("DropDatabase(%q): %v", dbName, err)
			}
		})
	}()

	// We execute all the commands of the tutorial code. These commands have to be run in a specific
	// order since in many cases earlier commands setup the database for the subsequent commands.
	mustRunCommand(t, "createdatabase", dbName)
	runCommand(t, "write", dbName)

	runCommand(t, "delete", dbName)
	runCommand(t, "write", dbName)
	writeTime := time.Now()

	assertContains(runCommand(t, "read", dbName), "1 1 Total Junk")

	assertContains(runCommand(t, "query", dbName), "1 1 Total Junk")

	runCommand(t, "addnewcolumn", dbName)
	runCommand(t, "update", dbName)

	runCommand(t, "writetransaction", dbName)
	out := runCommand(t, "querynewcolumn", dbName)
	assertContains(out, "1 1 300000")
	assertContains(out, "2 2 300000")

	runCommand(t, "addindex", dbName)
	out = runCommand(t, "queryindex", dbName)
	assertContains(out, "Go, Go, Go")
	assertContains(out, "Forever Hold Your Peace")
	if strings.Contains(out, "Green") {
		t.Errorf("got output %q; should not contain Green", out)
	}

	out = runCommand(t, "readindex", dbName)
	assertContains(out, "Go, Go, Go")
	assertContains(out, "Forever Hold Your Peace")
	assertContains(out, "Green")

	runCommand(t, "addstoringindex", dbName)
	assertContains(runCommand(t, "readstoringindex", dbName), "300000")
	out = runCommand(t, "readonlytransaction", dbName)
	if strings.Count(out, "Total Junk") != 2 {
		t.Errorf("got output %q; wanted it to contain 2 occurences of Total Junk", out)
	}

	// Wait at least 15 seconds since the write.
	time.Sleep(time.Until(writeTime.Add(16 * time.Second)))
	out = runCommand(t, "readstaledata", dbName)
	assertContains(out, "Go, Go, Go")
	assertContains(out, "Forever Hold Your Peace")
	assertContains(out, "Green")

	assertContains(runCommand(t, "readbatchdata", dbName), "1 Marc Richards")

	runCommand(t, "addcommittimestamp", dbName)
	runCommand(t, "updatewithtimestamp", dbName)
	out = runCommand(t, "querywithtimestamp", dbName)
	assertContains(out, "1000000")

	runCommand(t, "writestructdata", dbName)
	assertContains(runCommand(t, "querywithstruct", dbName), "6")
	out = runCommand(t, "querywitharrayofstruct", dbName)
	assertContains(out, "6")
	assertContains(out, "7")
	assertContains(out, "8")
	assertContains(runCommand(t, "querywithstructfield", dbName), "6")
	out = runCommand(t, "querywithnestedstructfield", dbName)
	assertContains(out, "6 Imagination")
	assertContains(out, "9 Imagination")

	runCommand(t, "createtabledocswithtimestamp", dbName)
	runCommand(t, "writetodocstable", dbName)
	runCommand(t, "updatedocstable", dbName)

	assertContains(runCommand(t, "querydocstable", dbName), "Hello World 1 Updated")

	runCommand(t, "createtabledocswithhistorytable", dbName)
	runCommand(t, "writewithhistory", dbName)
	runCommand(t, "updatewithhistory", dbName)

	out = runCommand(t, "querywithhistory", dbName)
	assertContains(out, "1 1 Hello World 1 Updated")

	out = runCommand(t, "dmlinsert", dbName)
	assertContains(out, "record(s) inserted")

	out = runCommand(t, "dmlupdate", dbName)
	assertContains(out, "record(s) updated")

	out = runCommand(t, "dmldelete", dbName)
	assertContains(out, "record(s) deleted")

	out = runCommand(t, "dmlwithtimestamp", dbName)
	assertContains(out, "record(s) updated")

	out = runCommand(t, "dmlwriteread", dbName)
	assertContains(out, "Found record name with ")

	out = runCommand(t, "dmlupdatestruct", dbName)
	assertContains(out, "record(s) inserted")

	out = runCommand(t, "dmlwrite", dbName)
	assertContains(out, "record(s) inserted")

	out = runCommand(t, "dmlwritetxn", dbName)
	assertContains(out, "from Album1's MarketingBudget to Album2")

	out = runCommand(t, "dmlupdatepart", dbName)
	assertContains(out, "record(s) updated")

	out = runCommand(t, "dmldeletepart", dbName)
	assertContains(out, "record(s) deleted")
}
