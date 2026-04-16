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
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAnywhereCache(t *testing.T) {
	tc := testutil.SystemTest(t)
	zone := os.Getenv("GOLANG_SAMPLES_ZONE")
	if zone == "" {
		zone = "us-central1-a"
	}
	ctx := context.Background()

	bucketName := testutil.UniqueBucketName(testPrefix)
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
		testutil.DeleteBucketIfExists(ctx, client, bucketName)
	})

	var cacheId string

	// Create Anywhere Cache
	t.Run("Create", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := createAnywhereCache(buf, bucketName, zone); err != nil {
			t.Fatalf("createAnywhereCache: %v", err)
		}
		got := buf.String()
		if !strings.Contains(got, "Created anywhere cache:") {
			t.Errorf("got %q, want it to contain 'Created anywhere cache:'", got)
		}
		// Extract cacheId from the name: projects/_/buckets/BUCKET/anywhereCaches/CACHE_ID
		trimmedGot := strings.TrimSpace(got)
		parts := strings.Split(trimmedGot, "/")
		cacheId = parts[len(parts)-1]
	})

	if cacheId == "" {
		t.Fatal("cacheId is empty, cannot continue tests")
	}

	// Get Anywhere Cache
	t.Run("Get", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := getAnywhereCache(buf, bucketName, cacheId); err != nil {
			t.Errorf("getAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheId; !strings.Contains(got, want) {
			t.Errorf("getAnywhereCache: got %q, want to contain %q", got, want)
		}
	})

	// List Anywhere Caches
	t.Run("List", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := listAnywhereCaches(buf, bucketName); err != nil {
			t.Errorf("listAnywhereCaches: %v", err)
		}
		if got, want := buf.String(), cacheId; !strings.Contains(got, want) {
			t.Errorf("listAnywhereCaches: got %q, want to contain %q", got, want)
		}
	})

	// Update Anywhere Cache
	t.Run("Update", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := updateAnywhereCache(buf, bucketName, cacheId, "admit-on-second-miss"); err != nil {
			t.Errorf("updateAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheId; !strings.Contains(got, want) {
			t.Errorf("updateAnywhereCache: got %q, want to contain %q", got, want)
		}
	})

	// Pause Anywhere Cache
	t.Run("Pause", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := pauseAnywhereCache(buf, bucketName, cacheId); err != nil {
			t.Errorf("pauseAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheId; !strings.Contains(got, want) {
			t.Errorf("pauseAnywhereCache: got %q, want to contain %q", got, want)
		}
	})

	// Resume Anywhere Cache
	t.Run("Resume", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := resumeAnywhereCache(buf, bucketName, cacheId); err != nil {
			t.Errorf("resumeAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheId; !strings.Contains(got, want) {
			t.Errorf("resumeAnywhereCache: got %q, want to contain %q", got, want)
		}
	})

	// Disable Anywhere Cache
	t.Run("Disable", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if err := disableAnywhereCache(buf, bucketName, cacheId); err != nil {
			t.Errorf("disableAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheId; !strings.Contains(got, want) {
			t.Errorf("disableAnywhereCache: got %q, want to contain %q", got, want)
		}
	})
}
