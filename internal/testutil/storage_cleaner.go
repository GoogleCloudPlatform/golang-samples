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
	"testing"
	"time"

	"cloud.google.com/go/storage"
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
	deleteBucketIfExists(ctx, t, client, bucket)

	b := client.Bucket(bucket)

	// Now create the bucket.
	// Retry because the bucket can take time to fully delete.
	Retry(t, 10, 10*time.Second, func(r *R) {
		if err := b.Create(ctx, projectID, nil); err != nil {
			r.Errorf("Bucket.Create(%q): %v", bucket, err)
		}
	})
	return nil
}

func deleteBucketIfExists(ctx context.Context, t *testing.T, client *storage.Client, bucket string) {
	t.Helper()

	b := client.Bucket(bucket)

	// Check if the bucket does not exist, return nil.
	if _, err := b.Attrs(ctx); err != nil {
		return
	}

	// Delete all of the elements in the already existent bucket.
	it := b.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Errorf("Bucket.Objects(%q): %v", bucket, err)
		}
		if attrs.EventBasedHold || attrs.TemporaryHold {
			if _, err := b.Object(attrs.Name).Update(ctx, storage.ObjectAttrsToUpdate{
				TemporaryHold:  false,
				EventBasedHold: false,
			}); err != nil {
				t.Errorf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
			}
		}
		if err := b.Object(attrs.Name).Delete(ctx); err != nil {
			t.Errorf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
		}
	}

	// Then delete the bucket itself.
	if err := b.Delete(ctx); err != nil {
		t.Errorf("Bucket.Delete(%q): %v", bucket, err)
	}
}
