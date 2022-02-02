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
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

// CreateTestBucket creates a new bucket with the given prefix
func CreateTestBucket(ctx context.Context, t *testing.T, client *storage.Client, projectID, prefix string) (string, error) {
	t.Helper()
	bucketName := UniqueBucketName(prefix)
	return bucketName, cleanBucketWithClient(ctx, t, client, projectID, bucketName)
}

// CleanBucket creates a new bucket. If the bucket already exists, it will be
// deleted and recreated.
func CleanBucket(ctx context.Context, t *testing.T, projectID, bucket string) error {
	t.Helper()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	return cleanBucketWithClient(ctx, t, client, projectID, bucket)
}

// cleanBucketWithClient creates a new bucket. If the bucket already exists, it will be
// deleted and recreated.
// Like CleanBucket but you must provide the storage client.
func cleanBucketWithClient(ctx context.Context, t *testing.T, client *storage.Client, projectID, bucket string) error {
	t.Helper()

	// Delete the bucket if it exists.
	if err := DeleteBucketIfExists(ctx, client, bucket); err != nil {
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
					DeleteBucketIfExists(ctx, client, bucket) // Ignore error.
				}
			}
			r.Errorf("Bucket.Create(%q): %v", bucket, err)
		}
	})

	WaitForBucketToExist(ctx, t, b)

	return nil
}

// DeleteBucketIfExists deletes a bucket and all its objects
func DeleteBucketIfExists(ctx context.Context, client *storage.Client, bucket string) error {
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
		// Objects with a hold must have the hold released
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

	// Waits for a bucket to no longer exist, as it can take time to propagate
	// Errors after 10 successful attempts at retrieving the bucket's attrs
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

// WaitForBucketToExist waits for a bucket to exist, as it can take time to propagate
// Errors after 10 unsuccessful attempts at retrieving the bucket's attrs
func WaitForBucketToExist(ctx context.Context, t *testing.T, b *storage.BucketHandle) {
	t.Helper()
	Retry(t, 10, 30*time.Second, func(r *R) {
		if _, err := b.Attrs(ctx); err != nil {
			// Bucket does not exist
			r.Errorf("Bucket was not created")
			return
		}
	})
}

// UniqueBucketName returns a unique name with the test prefix
// Any bucket created with this prefix may be deleted by DeleteExpiredBuckets
func UniqueBucketName(prefix string) string {
	return strings.Join([]string{prefix, uuid.New().String()}, "-")
}

// DeleteExpiredBuckets deletes old testing buckets that weren't cleaned previously
func DeleteExpiredBuckets(client *storage.Client, projectID, prefix string, expireAge time.Duration) error {
	ctx := context.Background()

	it := client.Buckets(ctx, projectID)
	it.Prefix = prefix
	for {
		bktAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if time.Since(bktAttrs.Created) > expireAge {
			log.Printf("deleting bucket %q, which is more than %s old", bktAttrs.Name, expireAge)
			if err := DeleteBucketIfExists(ctx, client, bktAttrs.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
