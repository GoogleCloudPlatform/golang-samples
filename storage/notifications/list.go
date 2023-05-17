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

// [START storage_list_bucket_notifications]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// listBucketNotifications lists notification configurations for a bucket.
func listBucketNotifications(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"

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

	for nID, n := range notifications {
		fmt.Fprintf(w, "Notification topic %s with ID %s\n", n.TopicID, nID)
	}

	return nil
}

// [END storage_list_bucket_notifications]
