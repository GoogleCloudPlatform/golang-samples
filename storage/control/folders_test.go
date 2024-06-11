// Copyright 2024 Google LLC
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

package control

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	testPrefix      = "storage-control-test"
	bucketExpiryAge = time.Hour * 24
)

var (
	client *storage.Client
)

func TestMain(m *testing.M) {
	// Create shared storage client
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer c.Close()
	client = c

	// Run tests
	exit := m.Run()

	// Delete old buckets whose name begins with our test prefix
	tc, _ := testutil.ContextMain(m)

	if err := testutil.DeleteExpiredBuckets(c, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
		// Don't fail the test if cleanup fails
		log.Printf("Post-test cleanup failed: %v", err)
	}
	os.Exit(exit)
}

func TestFolders(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	// Create HNS bucket.
	bucketName := testutil.UniqueBucketName(testPrefix)
	b := client.Bucket(bucketName)
	attrs := &storage.BucketAttrs{
		HierarchicalNamespace: &storage.HierarchicalNamespace{
			Enabled: true,
		},
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	}
	if err := b.Create(ctx, tc.ProjectID, attrs); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	t.Cleanup(func() {
		if err := testutil.DeleteBucketIfExists(ctx, client, bucketName); err != nil {
			log.Printf("Bucket.Delete(%q): %v", bucketName, err)
		}
	})

	folderName := "foo"
	folderPath := fmt.Sprintf("projects/_/buckets/%v/folders/%v", bucketName, folderName)
	newFolderName := "bar"
	newFolderPath := fmt.Sprintf("projects/_/buckets/%v/folders/%v", bucketName, newFolderName)

	// Create folder.
	buf := &bytes.Buffer{}
	if err := createFolder(buf, bucketName, folderName); err != nil {
		t.Fatalf("createFolder: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Errorf("createFolder: got %q, want to contain %q", got, want)
	}

	// Get folder.
	buf = &bytes.Buffer{}
	if err := getFolder(buf, bucketName, folderName); err != nil {
		t.Fatalf("getFolder: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Errorf("getFolder: got %q, want to contain %q", got, want)
	}

	// List folders.
	buf = &bytes.Buffer{}
	if err := listFolders(buf, bucketName); err != nil {
		t.Fatalf("listFolders: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Errorf("listFolders: got %q, want to contain %q", got, want)
	}

	// Rename folder.
	buf = &bytes.Buffer{}
	if err := renameFolder(buf, bucketName, folderName, newFolderName); err != nil {
		t.Fatalf("renameFolder: %v", err)
	}
	if got, want := buf.String(), newFolderPath; !strings.Contains(got, want) {
		t.Errorf("listFolders: got %q, want to contain %q", got, want)
	}

	// Delete folder.
	buf = &bytes.Buffer{}
	if err := deleteFolder(buf, bucketName, newFolderName); err != nil {
		t.Fatalf("deleteFolder: %v", err)
	}
	if got, want := buf.String(), newFolderPath; !strings.Contains(got, want) {
		t.Errorf("deleteFolder: got %q, want to contain %q", got, want)
	}
}
