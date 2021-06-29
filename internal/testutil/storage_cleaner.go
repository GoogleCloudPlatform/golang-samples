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

package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

// CleanBucket creates a new bucket. If the bucket already exists, it will be
// deleted and recreated.
func CleanBucket(ctx context.Context, t *testing.T, projectID, bucket string) error {
	t.Helper()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}

	// Delete the bucket if it exists.
	if err := deleteBucketIfExists(ctx, client, bucket); err != nil {
		return fmt.Errorf("error deleting bucket: %v", err)
	}

	b := client.Bucket(bucket)

	// Now create the bucket.
	// Retry because the bucket can take time to fully delete.
	Retry(t, 10, 30*time.Second, func(r *R) {
		if err := b.Create(ctx, projectID, nil); err != nil {
			if err, ok := err.(*googleapi.Error); ok {
				// Just in case...
				if err.Code == 409 {
					deleteBucketIfExists(ctx, client, bucket) // Ignore error.
				}
			}
			r.Errorf("Bucket.Create(%q): %v", bucket, err)
		}
	})

	// Wait until the bucket exists.
	Retry(t, 10, 30*time.Second, func(r *R) {
		if _, err := b.Attrs(ctx); err != nil {
			// Bucket does not exist.
			r.Errorf("Bucket was not created")
			return
		}
	})

	return nil
}

func deleteBucketIfExists(ctx context.Context, client *storage.Client, bucket string) error {
	b := client.Bucket(bucket)

	// Check if the bucket does not exist, return nil.
	if _, err := b.Attrs(ctx); err != nil {
		return nil
	}

	// Delete all of the elements in the already existent bucket, including noncurrent objects.
	it := b.Objects(ctx, &storage.Query{
		// Versions true to output all generations of objects.
		Versions: true,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket.Objects(%q): %v", bucket, err)
		}
		if attrs.EventBasedHold || attrs.TemporaryHold {
			if _, err := b.Object(attrs.Name).Update(ctx, storage.ObjectAttrsToUpdate{
				TemporaryHold:  false,
				EventBasedHold: false,
			}); err != nil {
				return fmt.Errorf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
			}
		}
		obj := b.Object(attrs.Name).Generation(attrs.Generation)
		if err := obj.Delete(ctx); err != nil {
			return fmt.Errorf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
		}

	}

	// Then delete the bucket itself.
	if err := b.Delete(ctx); err != nil {
		return fmt.Errorf("Bucket.Delete(%q): %v", bucket, err)
	}

	retries := 10
	delay := 10 * time.Second

	for i := 0; i < retries; i++ {
		if _, err := b.Attrs(ctx); err != nil {
			// Deletion successful.
			return nil
		}
		// Deletion not complete.
		time.Sleep(delay)
	}

	return fmt.Errorf("failed to delete bucket %q", bucket)
}
