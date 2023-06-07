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

// [START storage_create_bucket_turbo_replication]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// createBucketTurboReplication creates a new dual-region bucket with
// turbo replication enabled.
func createBucketTurboReplication(w io.Writer, projectID, bucketName, location string) error {
	// projectID := "my-project-id"
	// bucketName := "bucket-name"
	// location := "NAM4" // the name of a dual-region location
	// See this documentation for other valid locations:
	// https://cloud.google.com/storage/docs/locations#location-dr

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	storageLocationAndRPO := &storage.BucketAttrs{
		Location: location,
		RPO:      storage.RPOAsyncTurbo,
	}
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, storageLocationAndRPO); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Created bucket %v with turbo replication in %v\n", bucketName, storageLocationAndRPO.Location)
	return nil
}

// [END storage_create_bucket_turbo_replication]
