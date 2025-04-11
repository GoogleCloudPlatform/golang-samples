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

// [START storage_print_pubsub_bucket_notification]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// printPubsubBucketNotification gets a notification configuration for a bucket.
func printPubsubBucketNotification(w io.Writer, bucketName, notificationID string) error {
	// bucketName := "bucket-name"
	// notificationID := "notification-id"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	notifications, err := client.Bucket(bucketName).Notifications(ctx)
	if err != nil {
		return fmt.Errorf("Bucket.Notifications: %w", err)
	}

	n := notifications[notificationID]
	fmt.Fprintf(w, "Notification: %+v", n)

	return nil
}

// [END storage_print_pubsub_bucket_notification]
