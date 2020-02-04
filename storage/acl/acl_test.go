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
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

// TestACL runs all of the package tests.
func TestACL(t *testing.T) {
	tc := testutil.SystemTest(t)
	var (
		bucket                = tc.ProjectID + "-samples-object-bucket-1"
		object                = "foo.txt"
		buf                   bytes.Buffer
		allAuthenticatedUsers = storage.AllAuthenticatedUsers
		roleReader := storage.RoleReader
	)

	if err := addBucketACL(bucket, allAuthenticatedUsers, roleReader); err != nil {
		t.Errorf("addBucketACL %v", err)
	}
	if err := addDefaultBucketACL(bucket, allAuthenticatedUsers, roleReader); err != nil {
		t.Errorf("addDefaultBucketACL: %v", err)
	}
	if err := bucketListACL(&buf, bucket); err != nil {
		t.Errorf("bucketListACL: %v", err)
	}
	buf.Reset()
	if err := bucketListACLFiltered(&buf, bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("bucketListACLFiltered: %v", err)
	}
	buf.Reset()
	if err := deleteDefaultBucketACL(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("deleteDefaultBucketACL: %v", err)
	}
	if err := deleteBucketACL(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("deleteBucketACL: %v", err)
	}
	if err := addObjectACL(bucket, object, allAuthenticatedUsers, roleReader); err != nil {
		t.Errorf("addObjectACL: %v", err)
	}
	if err := objectListACL(&buf, bucket, object); err != nil {
		t.Errorf("objectListACL: %v", err)
	}
	buf.Reset()
	if err := objectListACLFiltered(&buf, bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("objectListACLFiltered: %v", err)
	}
	if err := deleteObjectACL(bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("deleteObjectACL: %v", err)
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
			r.Errorf("Bucket.Delete: %v", err)
		}
	})
}
