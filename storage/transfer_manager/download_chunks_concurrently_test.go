// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDownloadChunksConcurrently(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-tm-bucket-test"
	blobName := "tm-blob-test"
	fileName := "tm-file-test"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatalf("Bucket(%q).Create: %v", bucketName, err)
	}
	obj := bucket.Object(blobName)
	w := obj.NewWriter(ctx)
	if _, err := fmt.Fprint(w, "hello world"); err != nil {
		t.Fatalf("Writer.Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}

	var buf bytes.Buffer
	if err := downloadChunksConcurrently(&buf, bucketName, blobName, fileName); err != nil {
		t.Errorf("downloadChunksConcurrently: %v", err)
	}
	defer os.Remove(fileName)

	if got, want := buf.String(), fmt.Sprintf("Downloaded %v to %v", blobName, fileName); !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}

	// Clean up.
	if err := testutil.DeleteBucketIfExists(ctx, client, bucketName); err != nil {
		t.Fatalf("testutil.DeleteBucketIfExists: %v", err)
	}
}
