// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	gocontext "golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	rawbq "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/iterator"
)

var (
	idCounter uint64
	expiry    = time.Hour
)

const appPrefix = "golang_samples"

// NextDataset returns the next unique dataset ID generated for the given namespace.
func NextDataset(namespace string) string {
	atomic.AddUint64(&idCounter, 1)
	return fmt.Sprintf("%s_%s_dataset%d_%d", appPrefix, namespace, idCounter, time.Now().Unix())
}

// NextTable returns the next unique table ID generated for the given namespace
func NextTable(namespace string) string {
	atomic.AddUint64(&idCounter, 1)
	return fmt.Sprintf("%s_%s_table%d_%d", appPrefix, namespace, idCounter, time.Now().Unix())
}

// NextBucket will return a new unique bucket name for the given namespace.
func NextBucket(namespace string) string {
	atomic.AddUint64(&idCounter, 1)
	return fmt.Sprintf("%s_%s_bucket%d_%d", appPrefix, namespace, idCounter, time.Now().Unix())
}

// NextObject will return a new unique object name for the given namespace.
func NextObject(namespace string) string {
	atomic.AddUint64(&idCounter, 1)
	return fmt.Sprintf("%s_%s_object%d_%d", appPrefix, namespace, idCounter, time.Now().Unix())
}

// CleanupDatasets deletes the given datasets and all expired
// datasets belongs to the tests.
func CleanupDatasets(t *testing.T, name ...string) {
	ctx := gocontext.Background()
	tc := SystemTest(t)

	all := make(map[string]struct{})
	for _, n := range name {
		all[n] = struct{}{}
	}

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Errorf("bigquery.NewClient: %v", err)
		return
	}

	it := client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			t.Errorf("it.Next: %v", err)
			continue // Ignore, will be cleaned up at the next time.
		}
		creation, ok := extractTime(dataset.DatasetID)
		if ok && creation.Before(time.Now().Add(-expiry)) {
			all[dataset.DatasetID] = struct{}{}
		}

	}
	for d := range all {
		deleteDataset(t, d)
	}
}

func deleteDataset(t *testing.T, datasetID string) {
	ctx := gocontext.Background()
	tc := SystemTest(t)

	hc, err := google.DefaultClient(ctx, rawbq.CloudPlatformScope)
	if err != nil {
		t.Errorf("DefaultClient: %v", err)
	}
	s, err := rawbq.New(hc)
	if err != nil {
		t.Errorf("bigquery.New: %v", err)
	}
	call := s.Datasets.Delete(tc.ProjectID, datasetID)
	call.DeleteContents(true)
	call.Context(ctx)
	if err = call.Do(); err != nil {
		t.Errorf("deleteDataset(%q): %v", datasetID, err)
	}
}

// CleanupBuckets deletes all the expired and given buckets.
// If a bucket contains objects, it first deletes all the objects.
func CleanupBuckets(t *testing.T, name ...string) {
	ctx := gocontext.Background()
	tc := SystemTest(t)

	all := make(map[string]struct{})
	for _, n := range name {
		all[n] = struct{}{}
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("Cannot create client: %v", err)
		return // Will be cleanuped when expired
	}

	biter := client.Buckets(ctx, tc.ProjectID)
	for {
		b, err := biter.Next()
		if err != nil {
			break // Ignore errors, deletion will be retried at the next cleanup.
		}
		bucket := b.Name
		t, ok := extractTime(bucket)
		if ok && t.Before(time.Now().Add(-expiry)) {
			all[bucket] = struct{}{}
		}
	}

	for b := range all {
		bucket := client.Bucket(b)
		it := bucket.Objects(ctx, nil)

		// Delete all objects belongs to the bucket.
		// If deletion fails, skip to the next bucket,
		// next cleanup job will delete the garbage bucket.
		for {
			o, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// Ignore, will expire and deleted in the future.
				t.Logf("Cannot iterate objects (bucket=%v): %v", b, err)
				break
			}
			deleteObject(ctx, t, bucket, o.Name)
		}
		deleteBucket(ctx, t, bucket)
	}
}

func deleteObject(ctx gocontext.Context, t *testing.T, bucket *storage.BucketHandle, name string) {
	Retry(t, 10, time.Second, func(r *R) {
		if err := bucket.Object(name).Delete(ctx); err != nil {
			r.Errorf("cannot clean up object (bucket=%v, object=%v): %v", bucket, name, err)
		}
	})
}

func deleteBucket(ctx gocontext.Context, t *testing.T, bucket *storage.BucketHandle) {
	Retry(t, 10, time.Second, func(r *R) {
		if err := bucket.Delete(ctx); err != nil {
			r.Errorf("cannot clean up bucket (bucket=%v): %v", bucket, err)
		}
	})
}

// BucketsMustExist makes sure the given list of buckets exists.
// It creates the given bucket names if there are no matching buckets already.
// If any error occurs during creation or checking the existence, it fails
// the current test case.
func BucketsMustExist(t *testing.T, name ...string) {
	ctx := gocontext.Background()
	tc := SystemTest(t)

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("BucketsMustExist cannot create client: %v", err)
	}

	for _, n := range name {
		b := client.Bucket(n)
		_, err = b.Attrs(ctx)
		if err == storage.ErrBucketNotExist {
			err = b.Create(ctx, tc.ProjectID, nil)
		}
		if err != nil {
			t.Fatalf("bucket ensuring failed: %v", err)
		}
	}
}

func extractTime(s string) (time.Time, bool) {
	if !strings.HasPrefix(s, appPrefix) {
		return time.Time{}, false
	}
	i := strings.LastIndex(s, "_")
	if i < 0 {
		return time.Time{}, false
	}
	nanos, err := strconv.ParseInt(s[i+1:], 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	return time.Unix(0, nanos), true
}
