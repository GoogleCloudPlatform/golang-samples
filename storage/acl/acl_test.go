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
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(m *testing.M) {
	// These functions are noisy.
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

// TestAcls runs all of the package tests.
func TestAcls(t *testing.T) {
	tc := testutil.SystemTest(t)

	var (
		bucket  = tc.ProjectID + "-samples-object-bucket-1"
		object1 = "foo.txt"
	)

	var buf bytes.Buffer

	if err := addBucketACL(bucket); err != nil {
		t.Errorf("cannot add bucket acl: %v", err)
	}
	if err := addDefaultBucketACL(bucket); err != nil {
		t.Errorf("cannot add bucket default acl: %v", err)
	}
	if err := bucketListACL(&buf, bucket); err != nil {
		t.Errorf("cannot get bucket acl: %v", err)
	}
	if err := bucketListACLFiltered(&buf, bucket, storage.AllAuthenticatedUsers); err != nil {
		t.Errorf("cannot filter bucket acl: %v", err)
	}
	if err := deleteDefaultBucketACL(bucket); err != nil {
		t.Errorf("cannot delete bucket default acl: %v", err)
	}
	if err := deleteBucketACL(bucket); err != nil {
		t.Errorf("cannot delete bucket acl: %v", err)
	}
	if err := addObjectACL(bucket, object1); err != nil {
		t.Errorf("cannot add object acl: %v", err)
	}
	if err := objectListACL(bucket, object1); err != nil {
		t.Errorf("cannot get object acl: %v", err)
	}
	if err := objectListACLFiltered(bucket, object1, storage.AllAuthenticatedUsers); err != nil {
		t.Errorf("cannot filter object acl: %v", err)
	}
	if err := deleteObjectACL(bucket, object1); err != nil {
		t.Errorf("cannot delete object acl: %v", err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		// Cleanup, this part won't be executed if Fatal happens.
		// TODO(jbd): Implement garbage cleaning.
		if err := client.Bucket(bucket).Delete(ctx); err != nil {
			r.Errorf("cleanup of bucket failed: %v", err)
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
				t.Fatalf("Bucket.Objects(%q): %v", bucket, err)
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
			t.Fatalf("Bucket.Delete(%q): %v", bucket, err)
		}
	}
	if err := b.Create(ctx, projectID, nil); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucket, err)
	}
}
