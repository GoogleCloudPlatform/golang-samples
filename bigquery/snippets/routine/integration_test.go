// Copyright 2021 Google LLC
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

// Package routines demonstrates interactions with BigQuery Routine resources,
// which includes stored procedures and user defined functions.
package routine

import (
	"context"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRoutines(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	testDatasetID, err := bqtestutil.UniqueBQName("snippet_routine_tests")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}
	if err := client.Dataset(testDatasetID).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID).DeleteWithContents(ctx)

	testRoutineID, err := bqtestutil.UniqueBQName("testroutine")
	if err != nil {
		t.Fatalf("couldn't generate unique routine id: %v", err)
	}
	if err := createRoutineDDL(tc.ProjectID, testDatasetID, testRoutineID); err != nil {
		t.Fatalf("createRoutineDDL(%q %q): %v", testDatasetID, testRoutineID, err)
	}

	testRoutineID, err = bqtestutil.UniqueBQName("testroutine")
	if err != nil {
		t.Fatalf("couldn't generate unique routine id: %v", err)
	}
	if err := createRoutine(tc.ProjectID, testDatasetID, testRoutineID); err != nil {
		t.Fatalf("createRoutine(%q %q): %v", testDatasetID, testRoutineID, err)
	}
	if err := getRoutine(ioutil.Discard, tc.ProjectID, testDatasetID, testRoutineID); err != nil {
		t.Fatalf("getRoutine(%q %q): %v", testDatasetID, testRoutineID, err)
	}
	if err := listRoutines(ioutil.Discard, tc.ProjectID, testDatasetID); err != nil {
		t.Fatalf("listRoutines(%q): %v", testDatasetID, err)
	}
	if err := updateRoutine(tc.ProjectID, testDatasetID, testRoutineID); err != nil {
		t.Fatalf("updateRoutine(%q %q): %v", testDatasetID, testRoutineID, err)
	}
	if err := deleteRoutine(tc.ProjectID, testDatasetID, testRoutineID); err != nil {
		t.Fatalf("deleteRoutine(%q %q): %v", testDatasetID, testRoutineID, err)
	}
}
