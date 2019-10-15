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

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func TestSample(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("Skipping spanner integration test. Set GOLANG_SAMPLES_SPANNER.")
	}
	if !strings.HasPrefix(instance, "projects/") {
		t.Fatal("Spanner instance ref must be in the form of 'projects/PROJECT_ID/instances/INSTANCE_ID'")
	}
	dbName := fmt.Sprintf("%s/databases/test--%s", instance, tc.ProjectID)

	ctx := context.Background()
	adminClient, dataClient := createClients(ctx, dbName)
	defer adminClient.Close()
	defer dataClient.Close()

	// Check for database existance prior to test start and delete, as resources may not have
	// been cleaned up from previous invocations.
	if db, err := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbName}); err == nil {
		t.Logf("database %s exists in state %s. delete result: %v", db.GetName(), db.GetState().String(),
			adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName}))
	}

	assertContains := func(t *testing.T, out string, sub string) {
		if !strings.Contains(out, sub) {
			t.Errorf("got output %q; want it to contain %q", out, sub)
		}
	}
	runCommand := func(t *testing.T, cmd string, dbName string, timespan int) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName, timespan); err != nil {
			t.Errorf("run(%q, %q): %v", cmd, dbName, err)
		}
		return b.String()
	}
	mustRunCommand := func(t *testing.T, cmd string, dbName string, timespan int) string {
		t.Helper()
		var b bytes.Buffer
		if err := run(context.Background(), adminClient, dataClient, &b, cmd, dbName, timespan); err != nil {
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

	// These commands have to be run in a specific order
	// since earlier commands setup the database for the subsequent commands.
	mustRunCommand(t, "createdatabase", dbName, 0)
	assertContains(t, runCommand(t, "insertplayers", dbName, 0), "Inserted players")
	assertContains(t, runCommand(t, "insertplayers", dbName, 0), "Inserted players")
	assertContains(t, runCommand(t, "insertscores", dbName, 0), "Inserted scores")
	assertContains(t, runCommand(t, "insertscores", dbName, 0), "Inserted scores")
	assertContains(t, runCommand(t, "query", dbName, 0), "PlayerId: ")
	assertContains(t, runCommand(t, "querywithtimespan", dbName, 168), "PlayerId: ")
	assertContains(t, runCommand(t, "querywithtimespan", dbName, 730), "PlayerId: ")
	assertContains(t, runCommand(t, "querywithtimespan", dbName, 6870), "PlayerId: ")
}
