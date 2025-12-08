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

package subscriptions

// [START pubsub_old_version_create_cloud_storage_subscription]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

// createCloudStorageSubscription creates a Pub/Sub subscription that exports messages to Cloud Storage.
func createCloudStorageSubscription(w io.Writer, projectID, subID string, topic *pubsub.Topic, bucket string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	// note bucket should not have the gs:// prefix
	// bucket := "my-bucket"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic: topic,
		CloudStorageConfig: pubsub.CloudStorageConfig{
			Bucket:         bucket,
			FilenamePrefix: "log_events_",
			FilenameSuffix: ".avro",
			OutputFormat:   &pubsub.CloudStorageOutputFormatAvroConfig{WriteMetadata: true},
			MaxDuration:    1 * time.Minute,
			MaxBytes:       1e8,
		},
	})
	if err != nil {
		return fmt.Errorf("client.CreateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Created Cloud Storage subscription: %v\n", sub)

	return nil
}

// [END pubsub_old_version_create_cloud_storage_subscription]
