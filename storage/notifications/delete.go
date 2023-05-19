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

package notifications

// [START storage_delete_bucket_notification]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// deleteBucketNotification deletes a notification configuration for a bucket.
func deleteBucketNotification(w io.Writer, bucketName, notificationID string) error {
	// bucketName := "bucket-name"
	// notificationID := "notification-id"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	if err := bucket.DeleteNotification(ctx, notificationID); err != nil {
		return fmt.Errorf("Bucket.DeleteNotification: %w", err)
	}
	fmt.Fprintf(w, "Successfully deleted notification with ID %s for bucket %s.\n", notificationID, bucketName)
	return nil
}

// [END storage_delete_bucket_notification]
