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
	t.Skip("Skipping due to project permissions changes, see: b/445769988")
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

	// Create folder. Retry because there is no automatic retry in the client
	// for this op.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := createFolder(buf, bucketName, folderName); err != nil {
			r.Errorf("createFolder: %v", err)
		}
		if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
			r.Errorf("createFolder: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to create folder; can't continue")
	}

	// Get folder. Retry because there is no automatic retry in the client
	// for this op.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getFolder(buf, bucketName, folderName); err != nil {
			r.Errorf("getFolder: %v", err)
		}
		if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
			r.Errorf("getFolder: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to get folder; can't continue")
	}

	// List folders.
	buf := &bytes.Buffer{}
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

func TestManagedFolders(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.UniqueBucketName(testPrefix + "mf")
	b := client.Bucket(bucketName)
	attrs := &storage.BucketAttrs{
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

	folderName := "managed-foo"
	folderPath := fmt.Sprintf("projects/_/buckets/%v/managedFolders/%v/", bucketName, folderName)
	buf := &bytes.Buffer{}

	// Create Managed folder. Retry because there is no automatic retry in the client
	// for this op.
	buf.Reset()
	if err := createManagedFolder(buf, bucketName, folderName); err != nil {
		t.Fatalf("createManagedFolder: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Fatalf("createManagedFolder: got %q, want to contain %q", got, want)
	}

	// Get managed folder. Retry because there is no automatic retry in the client
	// for this op.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getManagedFolder(buf, bucketName, folderName); err != nil {
			r.Errorf("getManagedFolder: %v", err)
		}
		if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
			r.Errorf("getManagedFolder: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to get managed folder; can't continue")
	}

	// List managed folders.
	buf.Reset()
	if err := listManagedFolders(buf, bucketName); err != nil {
		t.Fatalf("listManagedFolders: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Errorf("listManagedFolders: got %q, want to contain %q", got, want)
	}

	// Delete managed folder.
	buf.Reset()
	if err := deleteManagedFolder(buf, bucketName, folderName); err != nil {
		t.Fatalf("deleteManagedFolder: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Errorf("deleteManagedFolder: got %q, want to contain %q", got, want)
	}
}

func TestAnywhereCaches(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	zone := os.Getenv("GOLANG_SAMPLES_ZONE")
	if zone == "" {
		zone = "us-central1-a"
	}

	bucketName := testutil.UniqueBucketName(testPrefix + "ac")
	b := client.Bucket(bucketName)
	attrs := &storage.BucketAttrs{
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

	admissionPolicy := "admit-on-first-miss"
	updatedAdmissionPolicy := "admit-on-second-miss"
	ttl := "3600s"

	// Create Anywhere Cache.
	buf := &bytes.Buffer{}
	if err := createAnywhereCache(buf, bucketName, zone, admissionPolicy, ttl); err != nil {
		t.Fatalf("createAnywhereCache: %v", err)
	}
	// The resource name should contain the bucket name.
	if got, want := buf.String(), bucketName; !strings.Contains(got, want) {
		t.Errorf("createAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Extract anywhereCacheID from the output.
	// Output format: "Anywhere Cache created: projects/_/buckets/bucket-name/anywhereCaches/anywhere-cache-id\n"
	output := buf.String()
	parts := strings.Split(strings.TrimSpace(output), "/")
	anywhereCacheID := parts[len(parts)-1]

	anywhereCachePath := fmt.Sprintf("buckets/%v/anywhereCaches/%v", bucketName, anywhereCacheID)

	// Get Anywhere Cache.
	buf.Reset()
	if err := getAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Fatalf("getAnywhereCache: %v", err)
	}
	if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
		t.Errorf("getAnywhereCache: got %q, want to contain %q", got, want)
	}

	// List Anywhere Caches.
	buf.Reset()
	if err := listAnywhereCaches(buf, bucketName); err != nil {
		t.Fatalf("listAnywhereCaches: %v", err)
	}
	if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
		t.Errorf("listAnywhereCaches: got %q, want to contain %q", got, want)
	}

	// Update Anywhere Cache.
	buf.Reset()
	if err := updateAnywhereCache(buf, bucketName, anywhereCacheID, updatedAdmissionPolicy); err != nil {
		t.Fatalf("updateAnywhereCache: %v", err)
	}
	if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
		t.Errorf("updateAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Pause Anywhere Cache.
	buf.Reset()
	if err := pauseAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Fatalf("pauseAnywhereCache: %v", err)
	}
	if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
		t.Errorf("pauseAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Resume Anywhere Cache.
	buf.Reset()
	if err := resumeAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Fatalf("resumeAnywhereCache: %v", err)
	}
	if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
		t.Errorf("resumeAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Disable Anywhere Cache.
	buf.Reset()
	if err := disableAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Fatalf("disableAnywhereCache: %v", err)
	}
	if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
		t.Errorf("disableAnywhereCache: got %q, want to contain %q", got, want)
	}
}
