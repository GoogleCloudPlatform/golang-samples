// Copyright 2020 Google LLC
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

package acl

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestACL runs all of the package tests.
func TestACL(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	var (
		bucket                = tc.ProjectID + "-samples-acl-bucket-1"
		object                = "foo.txt"
		allAuthenticatedUsers = storage.AllAuthenticatedUsers
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucket)

	b := client.Bucket(bucket)

	// Upload a test object with storage.Writer.
	wc := b.Object(object).NewWriter(ctx)
	if _, err = fmt.Fprint(wc, "Hello\nworld"); err != nil {
		t.Errorf("fmt.Fprint: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Errorf("Writer.Close: %v", err)
	}

	// Run all the tests.
	if err := addBucketOwner(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("addBucketOwner: %v", err)
	}
	if err := addBucketDefaultOwner(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("addBucketDefaultOwner: %v", err)
	}
	if err := printBucketACL(ioutil.Discard, bucket); err != nil {
		t.Errorf("printBucketACL: %v", err)
	}
	if err := printBucketACLForUser(ioutil.Discard, bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("printBucketACLForUser: %v", err)
	}
	if err := removeBucketDefaultOwner(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("removeBucketDefaultOwner: %v", err)
	}
	if err := removeBucketOwner(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("removeBucketOwner: %v", err)
	}
	if err := addFileOwner(bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("addFileOwner: %v", err)
	}
	if err := printFileACL(ioutil.Discard, bucket, object); err != nil {
		t.Errorf("printFileACL: %v", err)
	}
	if err := printFileACLForUser(ioutil.Discard, bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("printFileACLForUser: %v", err)
	}
	if err := removeFileOwner(bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("removeFileOwner: %v", err)
	}

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		// Cleanup, this part won't be executed if Fatal happens.
		// TODO(jbd): Implement garbage cleaning.
		b := client.Bucket(bucket)
		if err := b.Object(object).Delete(ctx); err != nil {
			r.Errorf("Object(%q).Delete: %v", object, err)
		}
		if err := b.Delete(ctx); err != nil {
			r.Errorf("Bucket(%q).Delete: %v", bucket, err)
		}
	})
}

// cleanBucket ensures there's a fresh bucket with a given name, deleting the existing bucket if it already exists.
func cleanBucket(t *testing.T, ctx context.Context, client *storage.Client, projectID, bucket string) {
	b := client.Bucket(bucket)
	_, err := b.Attrs(ctx)
	if err == nil {
		it := b.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("Bucket(%q).Objects: %v", bucket, err)
			}
			if attrs.EventBasedHold || attrs.TemporaryHold {
				if _, err := b.Object(attrs.Name).Update(ctx, storage.ObjectAttrsToUpdate{
					TemporaryHold:  false,
					EventBasedHold: false,
				}); err != nil {
					t.Fatalf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
				}
			}
			if err := b.Object(attrs.Name).Delete(ctx); err != nil {
				t.Fatalf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
			}
		}
		if err := b.Delete(ctx); err != nil {
			t.Fatalf("Bucket(%q).Delete: %v", bucket, err)
		}
	}
	if err := b.Create(ctx, projectID, nil); err != nil && status.Code(err) != codes.AlreadyExists {
		t.Fatalf("Bucket(%q).Create: %v", bucket, err)
	}
}
