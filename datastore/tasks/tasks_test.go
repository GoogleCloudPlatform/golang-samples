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
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var client *datastore.Client

func TestMain(m *testing.M) {
	ctx := context.Background()
	if tc, ok := testutil.ContextMain(m); ok {
		var err error
		client, err = datastore.NewClient(ctx, tc.ProjectID)
		if err != nil {
			log.Fatalf("datastore.NewClient: %v", err)
		}
		defer client.Close()
	}
	os.Exit(m.Run())
}

func makeDesc() string {
	return fmt.Sprintf("t-%d", time.Now().Unix())
}

func TestAddMarkDelete(t *testing.T) {
	tc := testutil.SystemTest(t)

	desc := makeDesc()

	k, err := AddTask(tc.ProjectID, desc)
	if err != nil {
		t.Fatal(err)
	}

	if err := MarkDone(tc.ProjectID, k.ID); err != nil {
		t.Fatal(err)
	}

	if err := DeleteTask(tc.ProjectID, k.ID); err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	t.Skip("Flaky. Eventual consistency. Re-enable once the datastore emulator works with gRPC.")

	tc := testutil.SystemTest(t)

	desc := makeDesc()

	k, err := AddTask(tc.ProjectID, desc)
	if err != nil {
		t.Fatal(err)
	}

	foundTask, err := listAndGetTask(tc.ProjectID, desc)
	if err != nil {
		t.Error(err)
	}
	if got, want := foundTask.id, k.ID; got != want {
		t.Errorf("k.ID: got %d, want %d", got, want)
	}

	if err := MarkDone(tc.ProjectID, foundTask.id); err != nil {
		t.Fatal(err)
	}

	foundTask, err = listAndGetTask(tc.ProjectID, desc)
	if err != nil {
		t.Error(err)
	}
	if !foundTask.Done {
		t.Error("foundTask.Done: got false, want true")
	}

	if err := DeleteTask(tc.ProjectID, foundTask.id); err != nil {
		t.Fatal(err)
	}
}

func TestCreateClientWithDatabase(t *testing.T) {
	test := testutil.SystemTest(t)
	projectID := test.ProjectID
	databaseName := "customdb"

	var buf bytes.Buffer
	err := createClientWithDatabase(&buf, projectID, databaseName)
	if err != nil {
		t.Fatal(err)
	}
}

func listAndGetTask(projectID string, desc string) (*Task, error) {

	tasks, err := ListTasks(projectID)
	if err != nil {
		return nil, err
	}

	var foundTask *Task
	for _, t := range tasks {
		if t.Desc == desc {
			foundTask = t
		}
	}
	if foundTask == nil {
		return nil, fmt.Errorf("Did not find task %s in list.", desc)
	}

	return foundTask, nil
}
