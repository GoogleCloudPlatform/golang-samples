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

package buckets

// [START storage_create_bucket_class_location]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// createBucketClassLocation creates a new bucket in the project with Storage class and
// location.
func createBucketClassLocation(w io.Writer, projectID, bucketName string) error {
	// projectID := "my-project-id"
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	storageClassAndLocation := &storage.BucketAttrs{
		StorageClass: "COLDLINE",
		Location:     "asia",
	}
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, storageClassAndLocation); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Created bucket %v in %v with storage class %v\n", bucketName, storageClassAndLocation.Location, storageClassAndLocation.StorageClass)
	return nil
}

// [END storage_create_bucket_class_location]
