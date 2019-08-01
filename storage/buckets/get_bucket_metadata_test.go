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

package inspect

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetBucketMetadata(t *testing.T) {
	testutil.SystemTest(t)
	setup(t)
	bucketName = tc.ProjectID + "-storage-buckets-tests"
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Clean up bucket before running tests
	deleteBucket(storageClient, bucketName)
	if err := create(storageClient, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucketName, err)
	}

	// Run test
	bucketMetadataBuf := new(bytes.Buffer)
	if err := getBucketMetadata(bucketMetadataBuf, storageClient, bucketName); err != nil {
		t.Errorf("getBucketMetadata: %#v", err)
	}

	got := bucketMetadataBuf.String()
	if want := "BucketName:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}