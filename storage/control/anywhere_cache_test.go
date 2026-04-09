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
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAnywhereCache(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	zone := os.Getenv("GOOGLE_SAMPLES_ZONE")
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

	anywhereCacheID := zone
	anywhereCachePath := fmt.Sprintf("buckets/%s/anywhereCaches/%s", bucketName, anywhereCacheID)

	// Create anywhere cache.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := createAnywhereCache(buf, bucketName, zone); err != nil {
			r.Errorf("createAnywhereCache: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("createAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to create anywhere cache; can't continue")
	}

	// Get anywhere cache.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
			r.Errorf("getAnywhereCache: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("getAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to get anywhere cache; can't continue")
	}

	// List anywhere caches.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := listAnywhereCaches(buf, bucketName); err != nil {
			r.Errorf("listAnywhereCaches: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("listAnywhereCaches: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to list anywhere caches; can't continue")
	}

	// Update anywhere cache.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := updateAnywhereCache(buf, bucketName, anywhereCacheID, "admit-on-second-miss"); err != nil {
			r.Errorf("updateAnywhereCache: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("updateAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to update anywhere cache; can't continue")
	}

	// Pause anywhere cache.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := pauseAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
			r.Errorf("pauseAnywhereCache: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("pauseAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to pause anywhere cache; can't continue")
	}

	// Resume anywhere cache.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := resumeAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
			r.Errorf("resumeAnywhereCache: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("resumeAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to resume anywhere cache; can't continue")
	}

	// Disable anywhere cache.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := disableAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
			r.Errorf("disableAnywhereCache: %v", err)
		}
		if got, want := buf.String(), anywhereCachePath; !strings.Contains(got, want) {
			r.Errorf("disableAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to disable anywhere cache; can't continue")
	}
}
