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

package buckets

// [START storage_create_bucket_dual_region]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// createBucketDualRegion creates a new dual-region bucket in the project in the
// provided locations.
func createBucketDualRegion(w io.Writer, projectID, bucketName, region1, region2 string) error {
	// projectID := "my-project-id"
	// bucketName := "bucket-name"
	// region1 := "US-EAST1"
	// region2 := "US-WEST1"

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	storageLocation := &storage.BucketAttrs{
		Location: fmt.Sprintf("%s+%s", region1, region2),
	}
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, storageLocation); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %v", bucketName, err)
	}
	fmt.Fprintf(w, "Created bucket %v in %v\n", bucketName, storageLocation.Location)
	return nil
}

// [END storage_create_bucket_dual_region]
