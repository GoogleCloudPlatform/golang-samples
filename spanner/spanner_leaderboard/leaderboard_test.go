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
	"strconv"
	"strings"
	"testing"
	"time"

	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func randomID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%s-%s", strconv.FormatInt(now.Unix(), 10), uuid.New().String()[:8])
}

func TestSample(t *testing.T) {
	testutil.EndToEndTest(t)

	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("Skipping spanner integration test. Set GOLANG_SAMPLES_SPANNER.")
	}
	if !strings.HasPrefix(instance, "projects/") {
		t.Fatal("Spanner instance ref must be in the form of 'projects/PROJECT_ID/instances/INSTANCE_ID'")
	}
	dbName := fmt.Sprintf("%s/databases/lb-%s", instance, randomID())

	ctx := context.Background()
	adminClient, dataClient := createClients(ctx, dbName)
	defer func() {
		dataClient.Close()
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			err := adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName})
			if err != nil {
				r.Errorf("DropDatabase(%q): %v", dbName, err)
			}
		})
		adminClient.Close()
	}()

	// Check for database existence prior to test start and delete, as resources may not have
	// been cleaned up from previous invocations.
	if db, err := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbName}); err == nil {
		t.Logf("database %s exists in state %s. delete result: %v", db.GetName(), db.GetState().String(),
			adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: dbName}))
	}

	assertContains := func(t *testing.T, name, out, sub string) {
		if !strings.Contains(out, sub) {
			t.Errorf("%s failed: got output %q; want it to contain %q", name, out, sub)
		}
	}
	runCommand := func(t *testing.T, cmd string, dbName string, timespan int) string {
		t.Helper()
		var b bytes.Buffer
		// Set timeout to 600 seconds so it should avoid DeadlineExceeded error.
		cctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()
		if err := run(cctx, adminClient, dataClient, &b, cmd, dbName, timespan); err != nil {
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

	// These commands have to be run in a specific order
	// since earlier commands setup the database for the subsequent commands.
	mustRunCommand(t, "createdatabase", dbName, 0)
	assertContains(t, "insertplayers", runCommand(t, "insertplayers", dbName, 0), "Inserted players")
	assertContains(t, "insertscores", runCommand(t, "insertscores", dbName, 0), "Inserted scores")
	assertContains(t, "query", runCommand(t, "query", dbName, 0), "PlayerId: ")
	assertContains(t, "querywithtimespan 18168", runCommand(t, "querywithtimespan", dbName, 18168), "PlayerId: ")
	assertContains(t, "querywithtimespan 18730", runCommand(t, "querywithtimespan", dbName, 18730), "PlayerId: ")
	assertContains(t, "querywithtimespan 186870", runCommand(t, "querywithtimespan", dbName, 186870), "PlayerId: ")
}
