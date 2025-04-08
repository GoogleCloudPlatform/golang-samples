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

// [START storage_remove_retention_policy]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// removeRetentionPolicy removes bucket retention policy.
func removeRetentionPolicy(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	bucket := client.Bucket(bucketName)
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).Attrs: %w", bucketName, err)
	}
	if attrs.RetentionPolicy.IsLocked {
		return fmt.Errorf("retention policy is locked")
	}

	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Retention period for %v has been removed\n", bucketName)
	return nil
}

// [END storage_remove_retention_policy]
