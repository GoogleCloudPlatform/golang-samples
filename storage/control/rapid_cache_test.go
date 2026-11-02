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
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/googleapi"
)

const (
	defaultRapidCacheZone = "us-central1-a"

	// Poll for up to 2 minutes for the cache to reach the "disabled" state.
	disablePollAttempts = 12
	disablePollInterval = 10 * time.Second

	// A bucket cannot be deleted until its disabled Rapid Cache has been
	// fully removed, which can take more than 30 minutes. Cleanup is
	// best-effort: retry briefly, then leave the bucket for the stale bucket
	// sweep in TestMain.
	bucketDeleteTimeout = 2 * time.Minute
	bucketDeleteBackoff = 15 * time.Second
)

func TestRapidCache(t *testing.T) {
	// Rapid Cache requires an allowlisted project, so these tests are skipped
	// unless explicitly enabled.
	if os.Getenv("STORAGE_RUN_RAPID_CACHE_TESTS") == "" {
		t.Skip("STORAGE_RUN_RAPID_CACHE_TESTS not set")
	}
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	zone := os.Getenv("GOLANG_SAMPLES_RAPID_CACHE_ZONE")
	if zone == "" {
		zone = defaultRapidCacheZone
	}
	// The bucket must be in the same region as the cache zone,
	// e.g. "us-central1-a" -> "us-central1".
	idx := strings.LastIndex(zone, "-")
	if idx <= 0 {
		t.Fatalf("invalid zone %q", zone)
	}
	region := zone[:idx]

	// Rapid Cache Ultra requires a regional bucket with hierarchical namespace
	// enabled.
	bucketName := testutil.UniqueBucketName(testPrefix + "-rc")
	attrs := &storage.BucketAttrs{
		Location: region,
		HierarchicalNamespace: &storage.HierarchicalNamespace{
			Enabled: true,
		},
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	}
	if err := client.Bucket(bucketName).Create(ctx, tc.ProjectID, attrs); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	t.Cleanup(func() { deleteRapidCacheBucket(ctx, bucketName) })

	// Create Rapid Cache.
	buf := &bytes.Buffer{}
	if err := createRapidCache(buf, bucketName, zone); err != nil {
		t.Fatalf("createRapidCache: %v", err)
	}
	const prefix = "Created Rapid Cache: "
	out := buf.String()
	i := strings.Index(out, prefix)
	if i < 0 {
		t.Fatalf("createRapidCache: got %q, want to contain %q", out, prefix)
	}
	// The cache ID is assigned by the server; parse it from the cache name.
	fields := strings.Fields(out[i+len(prefix):])
	if len(fields) == 0 {
		t.Fatalf("createRapidCache: no cache name in %q", out)
	}
	cacheName := fields[0]
	cacheID := path.Base(cacheName)
	if cacheID == "" || cacheID == "." || cacheID == "/" {
		t.Fatalf("could not parse cache ID from %q", cacheName)
	}

	// Disable the Rapid Cache even if an intermediate step fails, so the
	// bucket can eventually be deleted.
	disabled := false
	t.Cleanup(func() {
		if disabled {
			return
		}
		if err := disableRapidCache(&bytes.Buffer{}, bucketName, cacheID); err != nil {
			log.Printf("disableRapidCache(%q): %v", cacheName, err)
		}
	})

	// Get Rapid Cache.
	buf.Reset()
	if err := getRapidCache(buf, bucketName, cacheID); err != nil {
		t.Fatalf("getRapidCache: %v", err)
	}
	got := buf.String()
	for _, want := range []string{cacheName, "Zone: " + zone} {
		if !strings.Contains(got, want) {
			t.Errorf("getRapidCache: got %q, want to contain %q", got, want)
		}
	}

	// List Rapid Caches.
	buf.Reset()
	if err := listRapidCaches(buf, bucketName); err != nil {
		t.Fatalf("listRapidCaches: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, cacheName) {
		t.Errorf("listRapidCaches: got %q, want to contain %q", got, cacheName)
	}

	// Update Rapid Cache.
	buf.Reset()
	if err := updateRapidCache(buf, bucketName, cacheID, 48*time.Hour); err != nil {
		t.Errorf("updateRapidCache: %v", err)
	} else {
		got = buf.String()
		for _, want := range []string{"Updated Rapid Cache: " + cacheName, "TTL: 172800 seconds"} {
			if !strings.Contains(got, want) {
				t.Errorf("updateRapidCache: got %q, want to contain %q", got, want)
			}
		}
	}

	// Disable Rapid Cache.
	buf.Reset()
	if err := disableRapidCache(buf, bucketName, cacheID); err != nil {
		t.Fatalf("disableRapidCache: %v", err)
	}
	disabled = true
	if got, want := buf.String(), "Disabling Rapid Cache: "+cacheName; !strings.Contains(got, want) {
		t.Errorf("disableRapidCache: got %q, want to contain %q", got, want)
	}

	// The cache moves to the "disabled" state shortly after the request.
	for attempt := 0; attempt < disablePollAttempts; attempt++ {
		buf.Reset()
		if err := getRapidCache(buf, bucketName, cacheID); err != nil {
			t.Fatalf("getRapidCache: %v", err)
		}
		if strings.Contains(buf.String(), "State: disabled") {
			return
		}
		time.Sleep(disablePollInterval)
	}
	t.Errorf("getRapidCache: got %q, want to contain %q", buf.String(), "State: disabled")
}

// deleteRapidCacheBucket makes a best-effort attempt to delete a bucket that
// had a Rapid Cache. Deletion fails with 409 Conflict until the disabled cache
// has been fully removed, so after a short retry the bucket is left behind for
// the stale bucket sweep in TestMain. The test bucket holds no objects, so the
// bucket is deleted directly to inspect the returned error code.
func deleteRapidCacheBucket(ctx context.Context, bucketName string) {
	deadline := time.Now().Add(bucketDeleteTimeout)
	for {
		err := client.Bucket(bucketName).Delete(ctx)
		var gErr *googleapi.Error
		isGErr := errors.As(err, &gErr)
		if err == nil || errors.Is(err, storage.ErrBucketNotExist) || (isGErr && gErr.Code == http.StatusNotFound) {
			return
		}
		if !isGErr || gErr.Code != http.StatusConflict || time.Now().After(deadline) {
			log.Printf("Could not delete bucket %q; leaving it for later cleanup: %v", bucketName, err)
			return
		}
		time.Sleep(bucketDeleteBackoff)
	}
}
