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
	"google.golang.org/api/iterator"
)

// TestBucket creates a new bucket with the given prefix and registers a cleanup
// function to delete the bucket and any objects it contains when the test finishes.
// TestBucket returns the bucket name. It fails the test if bucket creation fails.
func TestBucket(ctx context.Context, t *testing.T, projectID, prefix string) string {
	t.Helper()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })
	return CreateTestBucket(ctx, t, client, projectID, prefix)
}

// CreateTestBucket creates a new bucket with the given prefix and registers a
// cleanup function to delete the bucket and any objects it contains.
// It is equivalent to TestBucket but allows Storage Client re-use.
func CreateTestBucket(ctx context.Context, t *testing.T, client *storage.Client, projectID, prefix string) string {
	t.Helper()
	bucketName := UniqueBucketName(prefix)

	b := client.Bucket(bucketName)
	if err := b.Create(ctx, projectID, nil); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}

	t.Cleanup(func() {
		if err := DeleteBucketIfExists(ctx, client, bucketName); err != nil {
			log.Printf("Bucket.Delete(%q): %v", bucketName, err)
		}
	})
	return bucketName
}

// DeleteBucketIfExists deletes a bucket and all its objects.
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
		obj := b.Object(attrs.Name)
		// Objects with a hold must have the hold released
		if attrs.EventBasedHold || attrs.TemporaryHold {
			if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{
				TemporaryHold:  false,
				EventBasedHold: false,
			}); err != nil {
				return fmt.Errorf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
			}
		}
		// Objects with a retention policy must must have the policy removed.
		if attrs.Retention != nil {
			_, err = obj.OverrideUnlockedRetention(true).Update(ctx, storage.ObjectAttrsToUpdate{
				Retention: &storage.ObjectRetention{},
			})
			if err != nil {
				return fmt.Errorf("failed to remove retention from object(%q): %v", attrs.Name, err)
			}
		}

		if err := obj.Generation(attrs.Generation).Delete(ctx); err != nil {
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

// UniqueBucketName returns a unique name with the test prefix.
func UniqueBucketName(prefix string) string {
	return strings.Join([]string{prefix, uuid.New().String()}, "-")
}

// DeleteExpiredBuckets deletes old testing buckets that weren't cleaned previously.
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
