// Copyright 2024 Google LLC
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

// [START storage_create_bucket_with_object_retention]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// createBucketObjectRetention creates a bucket with object retention enabled.
func createBucketObjectRetention(w io.Writer, projectID, bucketName string) error {
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

	bucket := client.Bucket(bucketName).SetObjectRetention(true)

	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}

	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).Attrs: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Created bucket %v with object retention mode: %v\n", bucketName, attrs.ObjectRetentionMode)
	return nil
}

// [END storage_create_bucket_with_object_retention]
