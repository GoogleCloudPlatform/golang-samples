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

// [START storage_create_bucket_notifications]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// createBucketNotification creates a notification configuration for a bucket.
func createBucketNotification(w io.Writer, projectID, bucketName, topic string) error {
	// projectID := "my-project-id"
	// bucketName := "bucket-name"
	// topic := "topic-name"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	notification := storage.Notification{
		TopicID:        topic,
		TopicProjectID: projectID,
		PayloadFormat:  storage.JSONPayload,
	}

	createdNotification, err := client.Bucket(bucketName).AddNotification(ctx, &notification)
	if err != nil {
		return fmt.Errorf("Bucket.AddNotification: %w", err)
	}
	fmt.Fprintf(w, "Successfully created notification with ID %s for bucket %s.\n", createdNotification.ID, bucketName)
	return nil
}

// [END storage_create_bucket_notifications]
