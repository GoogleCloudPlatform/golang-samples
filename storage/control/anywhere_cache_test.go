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
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

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
