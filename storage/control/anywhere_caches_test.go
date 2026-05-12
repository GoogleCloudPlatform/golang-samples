// Copyright 2026 Google LLC
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
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAnywhereCaches(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	// Use a local client to ensure isolation as per memory.
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	const testPrefix = "storage-control-ac-test"

	bucketName := testutil.UniqueBucketName(testPrefix)
	b := client.Bucket(bucketName)
	attrs := &storage.BucketAttrs{
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{Enabled: true},
	}
	if err := b.Create(ctx, tc.ProjectID, attrs); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	t.Cleanup(func() {
		if err := testutil.DeleteBucketIfExists(ctx, client, bucketName); err != nil {
			log.Printf("Bucket.Delete(%q): %v", bucketName, err)
		}
	})

	zoneName := "us-central1-a"
	cacheName := fmt.Sprintf("projects/_/buckets/%s/anywhereCaches/%s", bucketName, zoneName)

	// Create
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := createAnywhereCache(buf, bucketName, cacheName, zoneName); err != nil {
			r.Errorf("createAnywhereCache: %v", err)
		}
		// Match partial resource name to handle project numbers/IDs as per memory.
		want := fmt.Sprintf("buckets/%s/anywhereCaches/%s", bucketName, zoneName)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("createAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to create anywhere cache; can't continue")
	}

	// Get
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getAnywhereCache(buf, cacheName); err != nil {
			r.Errorf("getAnywhereCache: %v", err)
		}
		want := fmt.Sprintf("buckets/%s/anywhereCaches/%s", bucketName, zoneName)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("getAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to get anywhere cache; can't continue")
	}

	// List
	buf := &bytes.Buffer{}
	if err := listAnywhereCaches(buf, bucketName); err != nil {
		t.Fatalf("listAnywhereCaches: %v", err)
	}
	want := fmt.Sprintf("buckets/%s/anywhereCaches/%s", bucketName, zoneName)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("listAnywhereCaches: got %q, want to contain %q", got, want)
	}

	// Update
	buf.Reset()
	if err := updateAnywhereCache(buf, cacheName, "admit-on-second-miss"); err != nil {
		t.Fatalf("updateAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("updateAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Pause
	buf.Reset()
	if err := pauseAnywhereCache(buf, cacheName); err != nil {
		t.Fatalf("pauseAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("pauseAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Resume
	buf.Reset()
	if err := resumeAnywhereCache(buf, cacheName); err != nil {
		t.Fatalf("resumeAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("resumeAnywhereCache: got %q, want to contain %q", got, want)
	}

	// Disable
	buf.Reset()
	if err := disableAnywhereCache(buf, cacheName); err != nil {
		t.Fatalf("disableAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("disableAnywhereCache: got %q, want to contain %q", got, want)
	}
}
