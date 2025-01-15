// Copyright 2022 Google LLC
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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJobWithBucket(t *testing.T) {
	t.Parallel()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	region := "us-central1"
	jobName := fmt.Sprintf("test-job-go-bucket-%v-%v", time.Now().Format("2006-01-02"), r.Int())
	bucketName := fmt.Sprintf("test-bucket-go-batch-%v-%v", time.Now().Format("2006-01-02"), r.Int())
	buf := &bytes.Buffer{}

	if err := createBucket(tc.ProjectID, bucketName); err != nil {
		t.Errorf("Failed to create GCS bucket: createBucket got err: %v", err)
	}

	if err := createScriptJobWithBucket(buf, tc.ProjectID, region, jobName, bucketName); err != nil {
		t.Errorf("createScriptJobWithBucket got err: %v", err)
	}

	succeeded, err := jobSucceeded(tc.ProjectID, region, jobName)
	if err != nil {
		t.Errorf("Could not verify job completion: %v", err)
	}
	if !succeeded {
		t.Errorf("The test job has failed: %v", err)
	}

	// clean up after the test

	// GCS bucket will fail to delete if it is non-empty.
	// There is no 'force=true' argument in the Go GCS client,
	// and the low-level function 'newDeleteRequest' is not public,
	// so we cannot use force deletion here and have to first delete all files,
	// and only then delete the actual bucket.
	for i := 0; i <= 3; i++ {
		deleteFile(bucketName, fmt.Sprintf("output_task_%d.txt", i))
	}
	if err := deleteBucket(bucketName); err != nil {
		t.Errorf("Failed to delete GCS bucket: deleteBucket got err: %v", err)
	}
}

// createBucket creates a new bucket in the project.
func createBucket(projectID, bucketName string) error {
	// based on the storage_create_bucket sample
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}
	return nil
}

func deleteBucket(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := client.Bucket(bucketName)
	if err := bucket.Delete(ctx); err != nil {
		return fmt.Errorf("Bucket(%q).Delete: %w", bucketName, err)
	}
	return nil
}

// deleteFile removes specified object.
func deleteFile(bucket, object string) error {
	// copied from storage_delete_file sample
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}
	return nil
}
