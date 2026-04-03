// Copyright 2025 Google LLC
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

func TestAnywhereCaches(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	zone := os.Getenv("GOOGLE_CLOUD_CPP_TEST_ZONE")
	if zone == "" {
		t.Skip("GOOGLE_CLOUD_CPP_TEST_ZONE not set")
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

	// Expected suffix for resource name comparison
	cacheSuffix := fmt.Sprintf("buckets/%v/anywhereCaches/%v", bucketName, zone)
	buf := &bytes.Buffer{}

	// Create Anywhere Cache.
	// We might need to retry because the bucket might not be immediately available
	// for Storage Control operations.
	if ok := testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := createAnywhereCache(buf, bucketName, zone); err != nil {
			r.Errorf("createAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("createAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to create anywhere cache; can't continue")
	}

	// Get Anywhere Cache.
	if ok := testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := getAnywhereCache(buf, bucketName, zone); err != nil {
			r.Errorf("getAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("getAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Errorf("failed to get anywhere cache")
	}

	// List Anywhere Caches.
	if ok := testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := listAnywhereCaches(buf, bucketName); err != nil {
			r.Errorf("listAnywhereCaches: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("listAnywhereCaches: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Errorf("failed to list anywhere caches")
	}

	// Update Anywhere Cache.
	if ok := testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := updateAnywhereCache(buf, bucketName, zone, "admit-on-second-miss"); err != nil {
			r.Errorf("updateAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("updateAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Errorf("failed to update anywhere cache")
	}

	// Pause Anywhere Cache.
	if ok := testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := pauseAnywhereCache(buf, bucketName, zone); err != nil {
			r.Errorf("pauseAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("pauseAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Errorf("failed to pause anywhere cache")
	}

	// Resume Anywhere Cache.
	if ok := testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := resumeAnywhereCache(buf, bucketName, zone); err != nil {
			r.Errorf("resumeAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("resumeAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Errorf("failed to resume anywhere cache")
	}

	// Disable Anywhere Cache.
	if ok := testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := disableAnywhereCache(buf, bucketName, zone); err != nil {
			r.Errorf("disableAnywhereCache: %v", err)
		}
		if got, want := buf.String(), cacheSuffix; !strings.Contains(got, want) {
			r.Errorf("disableAnywhereCache: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Errorf("failed to disable anywhere cache")
	}
}
