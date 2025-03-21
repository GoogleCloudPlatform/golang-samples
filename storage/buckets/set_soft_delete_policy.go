// Copyright 2025 Google LLC
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

// [START storage_set_soft_delete_policy]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// Sets a bucket's soft delete policy with a 10-day retention period. Returns an error if the operation fails.
func setSoftDeletePolicy(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)

	// Set soft delete policy with 10-day retention period.
	_, err = bucket.Update(ctx, storage.BucketAttrsToUpdate{
		SoftDeletePolicy: &storage.SoftDeletePolicy{
			RetentionDuration: 10 * 24 * time.Hour, // 10 days in hours
		},
	})
	if err != nil {
		return fmt.Errorf("bucket.Update: %w", err)
	}

	fmt.Fprintf(w, "Soft delete policy for %s was set to a 10-day retention period\n", bucketName)
	return nil
}

// [END storage_set_soft_delete_policy]
