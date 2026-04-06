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

package firestore

import (
	"bytes"
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"cloud.google.com/go/firestore"
)

func testQueryFilterOrSetup(projectID string) ([]*firestore.DocumentRef, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bw := client.BulkWriter(ctx)
	colName := "users"

	docs := []struct {
		shortName string
		birthYear int
	}{
		{shortName: "aturing", birthYear: 1912},
		{shortName: "alovelace", birthYear: 1815},
		{shortName: "cbabbage", birthYear: 1791},
		{shortName: "ghopper", birthYear: 1906},
	}
	var refs []*firestore.DocumentRef

	for _, d := range docs {
		ref := client.Collection(colName).Doc(d.shortName)
		_, err := bw.Create(ref, map[string]interface{}{"birthYear": d.birthYear})
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	bw.End()
	return refs, nil
}

func testQueryFilterOrCleanup(projectID string, refs []*firestore.DocumentRef) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	// New BulkWriter instance
	bw := client.BulkWriter(ctx)

	for _, d := range refs {
		_, err := bw.Delete(d)
		if err != nil {
			return err
		}
	}
	bw.End()
	return nil
}

var projectIDOnce sync.Once

func getProjectID(t *testing.T) string {
	projectIDOnce.Do(func() {
		if projectID == "" {
			projectID = os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
		}
	})
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}
	return projectID
}

func TestQueryFilterOr(t *testing.T) {
	pid := getProjectID(t)

	refs, err := testQueryFilterOrSetup(pid)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := testQueryFilterOrCleanup(pid, refs)
		if err != nil {
			t.Fatal(err)
		}
	})

	var buf bytes.Buffer
	err = queryFilterOr(&buf, pid)
	if err != nil {
		t.Fatal(err)
	}

	want := "ghopper"
	got := buf.String()
	if !strings.Contains(got, want) {
		t.Errorf("Wanted: %s; got: %s", want, got)
	}
}
