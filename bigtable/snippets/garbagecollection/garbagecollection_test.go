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

package garbagecollection

import (
	"bytes"
	"strings"
	"testing"

	"cloud.google.com/go/bigtable"
	"context"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"log"
	"os"
)

func TestSetup(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE")
	adminClient, err := bigtable.NewAdminClient(ctx, tc.ProjectID, instance)

	tableName := "gc-table"
	if err = adminClient.DeleteTable(ctx, tableName); err != nil {
		log.Printf("Could not delete table %s: %v", tableName, err)
	}

	if err := adminClient.CreateTable(ctx, tableName); err != nil {
		log.Fatalf("Could not create table %s: %v", tableName, err)
	}
}

func TestMaxAge(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := createFamilyGCMaxAge(buf, tc.ProjectID, "test-inst", "gc-table")
	if err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got := buf.String()
	if want := "created column family cf1 with policy: age() > 5d\n"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaxVersions(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := createFamilyGCMaxVersions(buf, tc.ProjectID, "test-inst", "gc-table")
	if err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got := buf.String()
	if want := "created column family cf2 with policy: versions() > 2\n"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestUnion(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := createFamilyGCUnion(buf, tc.ProjectID, "test-inst", "gc-table")
	if err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got := buf.String()
	if want := "created column family cf3 with policy: (versions() > 2 || age() > 5d)\n"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIntersect(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := createFamilyGCIntersect(buf, tc.ProjectID, "test-inst", "gc-table")
	if err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got := buf.String()
	if want := "created column family cf4 with policy: (versions() > 2 && age() > 5d)\n"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
func TestNested(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := createFamilyGCNested(buf, tc.ProjectID, "test-inst", "gc-table")
	if err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got := buf.String()
	if want := "created column family cf5 with policy: (versions() > 10 || (versions() > 2 && age() > 5d))\n"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTeardown(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE")
	adminClient, err := bigtable.NewAdminClient(ctx, tc.ProjectID, instance)

	tableName := "gc-table"
	if err = adminClient.DeleteTable(ctx, tableName); err != nil {
		log.Printf("Could not delete table %s: %v", tableName, err)
	}
}
