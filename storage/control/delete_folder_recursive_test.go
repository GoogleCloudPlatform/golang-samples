// Copyright 2026 Google LLC
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

package control

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeleteFolderRecursive(t *testing.T) {
	t.Skip("Skipping due to project permissions changes, see: b/445769988")
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	// Create HNS bucket.
	bucketName := testutil.UniqueBucketName(testPrefix)
	b := client.Bucket(bucketName)
	attrs := &storage.BucketAttrs{
		HierarchicalNamespace: &storage.HierarchicalNamespace{
			Enabled: true,
		},
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	}
	if err := b.Create(ctx, tc.ProjectID, attrs); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	t.Cleanup(func() {
		if err := testutil.DeleteBucketIfExists(ctx, client, bucketName); err != nil {
			log.Printf("Bucket.Delete(%q): %v", bucketName, err)
		}
	})

	folderName := "foo-recursive"
	folderPath := fmt.Sprintf("projects/_/buckets/%v/folders/%v", bucketName, folderName)

	// Create folder. Retry because there is no automatic retry in the client
	// for this op.
	if ok := testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := createFolder(buf, bucketName, folderName); err != nil {
			r.Errorf("createFolder: %v", err)
		}
		if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
			r.Errorf("createFolder: got %q, want to contain %q", got, want)
		}
	}); !ok {
		t.Fatalf("failed to create folder; can't continue")
	}

	// Delete folder recursive.
	buf := &bytes.Buffer{}
	if err := deleteFolderRecursive(buf, bucketName, folderName); err != nil {
		t.Fatalf("deleteFolderRecursive: %v", err)
	}
	if got, want := buf.String(), folderPath; !strings.Contains(got, want) {
		t.Errorf("deleteFolderRecursive: got %q, want to contain %q", got, want)
	}
}
