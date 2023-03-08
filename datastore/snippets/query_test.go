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
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var projectID string

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("test project not set up properly")
		return
	}

	projectID = tc.ProjectID

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// Run the sample test
	m.Run()

	// Do teardown tasks
	err = client.Delete(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
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
