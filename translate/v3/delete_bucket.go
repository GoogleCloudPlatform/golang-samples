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

package v3

import (
	"context"
	"testing"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Helper function to delete a bucket with all their containing objects.
func deleteBucket(ctx context.Context, t *testing.T, bucket *storage.BucketHandle) {
	bucketAttrs, err := bucket.Attrs(ctx)
	if err != nil {
		t.Errorf("bucket.Attrs: %v", err)
	}
	bucketName := bucketAttrs.Name
	it := bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Errorf("bucket.Objects: %v", err)
		}
		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			t.Errorf("Bucket(%v).Object(%v).Delete: %v", bucketName, attrs.Name, err)
		}
	}
	if err := bucket.Delete(ctx); err != nil {
		t.Errorf("bucket.Delete: %v", err)
	}
}
