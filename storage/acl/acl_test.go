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
	"ioutil"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestACL runs all of the package tests.
func TestACL(t *testing.T) {
	tc := testutil.SystemTest(t)
	var (
		bucket                = tc.ProjectID + "-samples-object-bucket-1"
		object                = "foo.txt"
		allAuthenticatedUsers = storage.AllAuthenticatedUsers
		roleReader            = storage.RoleReader
	)

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
