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
<<<<<<< HEAD
	"crypto/rand"
	"fmt"
	"io"
	"log"
=======
	"fmt"
>>>>>>> df0f473c (Add sample download_chunks_concurrently with test.)
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
<<<<<<< HEAD
	"github.com/google/uuid"
)

const (
	testPrefix     = "storage-objects-test"
	downloadObject = "tm-obj-download"
)

var (
	tmBucketName  string
	storageClient *storage.Client
	downloadData  []byte
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	tc, _ := testutil.ContextMain(m)

	var err error

	// Create fixture client & bucket to use across tests.
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	tmBucketName = fmt.Sprintf("%s-%s", testPrefix, uuid.New().String())
	bucket := storageClient.Bucket(tmBucketName)
	if err := bucket.Create(ctx, tc.ProjectID, nil); err != nil {
		log.Fatalf("Bucket(%q).Create: %v", tmBucketName, err)
	}

	// Create object fixture for download tests.
	w := bucket.Object(downloadObject).NewWriter(ctx)
	downloadData = make([]byte, 2*1024*1024) // 2 MiB
	if _, err := rand.Read(downloadData); err != nil {
		log.Fatalf("rand.Read: %v", err)
	}
	if _, err := io.Copy(w, bytes.NewReader(downloadData)); err != nil {
		log.Fatalf("uploading object: %v", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("closing writer: %v", err)
	}

	// Run tests.
	exitCode := m.Run()

	// Cleanup bucket and objects.
	if err := testutil.DeleteBucketIfExists(ctx, storageClient, tmBucketName); err != nil {
		log.Printf("deleting bucket: %v", err)
	}
	os.Exit(exitCode)
}

func TestDownloadChunksConcurrently(t *testing.T) {
	bucketName := tmBucketName
	blobName := downloadObject

	// Create a temporary file to download to, ensuring we have permissions
	// and the file is cleaned up.
	f, err := os.CreateTemp("", "tm-file-test-")
	if err != nil {
		t.Fatalf("os.CreateTemp: %v", err)
	}
	fileName := f.Name()
	f.Close() // Close the file so the download can write to it.
	defer os.Remove(fileName)
=======
	"google.golang.org/api/iterator"
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
	defer deleteBucket(ctx, t, bucket)

	obj := bucket.Object(blobName)
	w := obj.NewWriter(ctx)
	if _, err := fmt.Fprint(w, "hello world"); err != nil {
		t.Fatalf("Writer.Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}
>>>>>>> df0f473c (Add sample download_chunks_concurrently with test.)

	var buf bytes.Buffer
	if err := downloadChunksConcurrently(&buf, bucketName, blobName, fileName); err != nil {
		t.Errorf("downloadChunksConcurrently: %v", err)
	}
<<<<<<< HEAD
=======
	defer os.Remove(fileName)
>>>>>>> df0f473c (Add sample download_chunks_concurrently with test.)

	if got, want := buf.String(), fmt.Sprintf("Downloaded %v to %v", blobName, fileName); !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
<<<<<<< HEAD

	// Verify that the downloaded data is the same as the uploaded data.
	downloadedBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%q): %v", fileName, err)
	}

	if !bytes.Equal(downloadedBytes, downloadData) {
		t.Errorf("downloaded data does not match uploaded data. got %d bytes, want %d bytes", len(downloadedBytes), len(downloadData))
=======
}

func deleteBucket(ctx context.Context, t *testing.T, bucket *storage.BucketHandle) {
	t.Helper()
	it := bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Logf("Bucket(%v).Objects: %v", bucket, err)
			break
		}
		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			t.Logf("Bucket(%v).Object(%q).Delete: %v", bucket, attrs.Name, err)
		}
	}
	if err := bucket.Delete(ctx); err != nil {
		t.Logf("Bucket(%v).Delete: %v", bucket, err)
>>>>>>> df0f473c (Add sample download_chunks_concurrently with test.)
	}
}
