// Copyright 2023 Google LLC
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

// Package datastore_snippets contains snippet code for the Cloud Datastore API.
// The code is not runnable.

package datastore_snippets

import (
	"bytes"
	"context"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/datastore"
	admin "cloud.google.com/go/datastore/admin/apiv1"
	"cloud.google.com/go/datastore/admin/apiv1/adminpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var projectID string

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Print("test project not set up properly")
		return 1
	}

	projectID = tc.ProjectID

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Print(err)
		return 1
	}
	defer client.Close()

	// Do setup tasks
	task := struct {
		Task string
	}{
		Task: "simpleTask",
	}

	key := datastore.IncompleteKey("TaskList", nil)
	key, err = client.Put(ctx, key, &task)
	if err != nil {
		log.Print(err)
		return 1
	}

	// Run the sample test
	code := m.Run()

	// Do teardown tasks
	err = client.Delete(ctx, key)
	if err != nil {
		log.Print(err)
		if code == 0 {
			code = 1
		}
	}
	return code
}

func TestNotEqualQuery(t *testing.T) {
	var buf bytes.Buffer

	err := queryNotEquals(&buf, projectID)
	if err != nil {
		t.Fatal(err)
	}

	result := buf.String()
	if result == "" {
		t.Error("didn't get result")
	}
}

func TestInQuery(t *testing.T) {
	var buf bytes.Buffer

	err := queryIn(&buf, projectID)
	if err != nil {
		t.Fatal(err)
	}

	result := buf.String()
	if result == "" {
		t.Error("didn't get result")
	}
}

func TestNotInQuery(t *testing.T) {
	var buf bytes.Buffer

	err := queryNotIn(&buf, projectID)
	if err != nil {
		t.Fatal(err)
	}

	result := buf.String()
	if result == "" {
		t.Error("didn't get result")
	}
}

func TestMultipleInequalitiesQuery(t *testing.T) {
	keyPrefix := uuid.NewString()
	// Create client
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		client.Close()
	})

	// Create entities
	type Task struct {
		Category    string `datastore:"category"`
		Done        bool   `datastore:"done"`
		Priority    int64  `datastore:"priority"`
		Days        int64  `datastore:"days"`
		Description string `datastore:"description"`
	}
	keys := []*datastore.Key{
		datastore.NameKey("Task", keyPrefix+"-key1", nil),
		datastore.NameKey("Task", keyPrefix+"-key2", nil),
		datastore.NameKey("Task", keyPrefix+"-key3", nil),
	}
	tasks := []Task{
		{
			Category:    "Personal",
			Priority:    4,
			Days:        3,
			Description: "Learn Cloud Datastore",
		},
		{
			Category:    "Personal",
			Priority:    5,
			Days:        5,
			Description: "Integrate Cloud Datastore",
		},
		{
			Category:    "Personal",
			Priority:    5,
			Days:        2,
			Description: "Set Up Cloud Datastore",
		},
	}
	if _, err := client.PutMulti(ctx, keys, tasks); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := client.DeleteMulti(ctx, keys); err != nil {
			t.Error(err)
		}
	})

	// Create required indexes
	adminClient, err := admin.NewDatastoreAdminClient(ctx)
	if err != nil {
		t.Fatalf("admin.NewDatastoreAdminClient: %v", err)
	}
	t.Cleanup(func() {
		adminClient.Close()
	})
	setupIndex(t, ctx, adminClient, projectID)

	// Run query
	var buf bytes.Buffer
	if err = queryMultipleInequality(&buf, projectID); err != nil {
		t.Fatal(err)
	}

	// Compare results
	got := buf.String()
	want := "Key: /Task," + keyPrefix + "-key3. Entity: {5 2 false Personal Set Up Cloud Datastore}\n"
	if !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func setupIndex(t *testing.T, ctx context.Context, adminClient *admin.DatastoreAdminClient, projectID string) {
	t.Helper()
	index := &adminpb.Index{
		Kind:     "Task",
		Ancestor: adminpb.Index_NONE,
		Properties: []*adminpb.Index_IndexedProperty{
			{Name: "days", Direction: adminpb.Index_ASCENDING},
			{Name: "priority", Direction: adminpb.Index_ASCENDING},
		},
	}
	createOp, err := adminClient.CreateIndex(ctx, &adminpb.CreateIndexRequest{
		ProjectId: projectID,
		Index:     index,
	})
	var indexID string
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			it := adminClient.ListIndexes(ctx, &adminpb.ListIndexesRequest{ProjectId: projectID})
			for {
				idx, err := it.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					t.Fatalf("ListIndexes: %v", err)
				}
				if idx.Kind == index.Kind && idx.Ancestor == index.Ancestor && reflect.DeepEqual(idx.Properties, index.Properties) {
					indexID = idx.IndexId
					break
				}
			}
		} else {
			t.Fatalf("CreateIndex: %v", err)
		}
	} else {
		createdIndex, err := createOp.Wait(ctx)
		if err != nil {
			t.Fatalf("CreateIndex Wait: %v", err)
		}
		indexID = createdIndex.IndexId
	}

	if indexID != "" {
		t.Cleanup(func() {
			deleteOp, err := adminClient.DeleteIndex(ctx, &adminpb.DeleteIndexRequest{
				ProjectId: projectID,
				IndexId:   indexID,
			})
			if err != nil {
				t.Errorf("DeleteIndex: %v", err)
			}
			if _, err := deleteOp.Wait(ctx); err != nil {
				t.Errorf("DeleteIndex Wait: %v", err)
			}
		})
	}
}
