// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSample(t *testing.T) {
	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("Skipping spanner integration test. Set GOLANG_SAMPLES_SPANNER.")
	}
	if !strings.HasPrefix(instance, "projects/") {
		t.Fatal("Spanner instance ref must be in the form of 'projects/PROJECT_ID/instances/INSTANCE_ID'")
	}
	dbName := fmt.Sprintf("%s/databases/test-%d", instance, time.Now().Unix())

	ctx := context.Background()
	adminClient, dataClient := createClients(ctx, dbName)
	defer adminClient.Close()
	defer dataClient.Close()

	assertContains := func(out string, sub string) {
		if !strings.Contains(out, sub) {
			t.Errorf("got output %q; want it to contain %s", out, sub)
		}
	}
	runCommand := func(t *testing.T, cmd string, dbName string) string {
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName); err != nil {
			t.Errorf("run(%q, %q): %v", cmd, dbName, err)
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
	runCommand(t, "createdatabase", dbName)
	runCommand(t, "write", dbName)

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
}
