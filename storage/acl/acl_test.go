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
	"io"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
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
		bucket                = testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, "samples-acl-bucket-1")
		object                = "foo.txt"
		allAuthenticatedUsers = storage.AllAuthenticatedUsers
	)

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
	if err := addBucketDefaultOwner(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("addBucketDefaultOwner: %v", err)
	}
	if err := printBucketACLForUser(io.Discard, bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("printBucketACLForUser: %v", err)
	}
	if err := removeBucketDefaultOwner(bucket, allAuthenticatedUsers); err != nil {
		t.Errorf("removeBucketDefaultOwner: %v", err)
	}
	if err := addFileOwner(bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("addFileOwner: %v", err)
	}
	if err := printFileACL(io.Discard, bucket, object); err != nil {
		t.Errorf("printFileACL: %v", err)
	}
	if err := printFileACLForUser(io.Discard, bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("printFileACLForUser: %v", err)
	}
	if err := removeFileOwner(bucket, object, allAuthenticatedUsers); err != nil {
		t.Errorf("removeFileOwner: %v", err)
	}
}
