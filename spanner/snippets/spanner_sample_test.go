/*
Copyright 2016 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/admin/database/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"

	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

var (
	// testInstanceID is the instance used for testing.
	testInstanceID = "go-integration-test"

	// adminClient is a client for Cloud Spanner's database admin api.
	adminClient *database.DatabaseAdminClient
	// dataClient is a client for Cloud Spanner's data plane api.
	dataClient *spanner.Client
	// db is the path of the testing database.
	db = "go-samples-test-db"
)

func runCommand(t *testing.T, cmd string, dbName string) string {
	var b bytes.Buffer
	run(context.Background(), adminClient, dataClient, &b, cmd, dbName)
	return b.String()
}

func assertContains(t *testing.T, out string, sub string) {
	if !strings.Contains(out, sub) {
		t.Errorf("got output %q; want it to contain %s", out, sub)
	}
}

func dropDatabase(ctx context.Context, t *testing.T, dbName string) {
	if err := adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{
		Database: dbName,
	}); err != nil {
		log.Printf("dropDatabase failed with %v", err)
	}

}

func TestSample(t *testing.T) {
	c := testutil.SystemTest(t)

	var (
		instance = "projects/" + c.ProjectID + "/instances/" + testInstanceID
		dbName   = instance + "/databases/" + db
		ctx      = context.Background()
	)
	adminClient, dataClient = createClients(ctx, dbName)
	defer adminClient.Close()
	defer dataClient.Close()

	// Try to the drop the database in case previous run failed to clean up
	dropDatabase(ctx, t, dbName)

	// We execute all the commands of the tutorial code. These commands have to be run in a specific
	// order since in many cases earlier commands setup the database for the subsequent commands.
	runCommand(t, "createdatabase", dbName)
	runCommand(t, "write", dbName)

	assertContains(t, runCommand(t, "read", dbName), "1 1 Total Junk")

	assertContains(t, runCommand(t, "query", dbName), "1 1 Total Junk")

	runCommand(t, "addnewcolumn", dbName)
	runCommand(t, "update", dbName)

	runCommand(t, "writetransaction", dbName)
	out := runCommand(t, "querynewcolumn", dbName)
	assertContains(t, out, "1 1 300000")
	assertContains(t, out, "2 2 300000")

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

	runCommand(t, "addstoringindex", dbName)
	assertContains(t, runCommand(t, "readstoringindex", dbName), "300000")
	out = runCommand(t, "readonlytransaction", dbName)
	if strings.Count(out, "Total Junk") != 2 {
		t.Errorf("got output %q; wanted it to contain 2 occurences of Total Junk", out)
	}
	dropDatabase(ctx, t, dbName)
}
