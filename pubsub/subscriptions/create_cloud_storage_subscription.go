// Copyright 2023 Google LLC
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

// [START pubsub_create_cloud_storage_subscription]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// createCloudStorageSubscription creates a Pub/Sub subscription that exports messages to Cloud Storage.
func createCloudStorageSubscription(w io.Writer, projectID, topic, subscription, bucket string) error {
	// projectID := "my-project-id"
	// topic := "projects/my-project-id/topics/my-topic"
	// subscription := "projects/my-project/subscriptions/my-sub"
	// bucket := "my-bucket" // bucket must not have the gs:// prefix
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subscription,
		Topic: topic,
		CloudStorageConfig: &pubsubpb.CloudStorageConfig{
			Bucket:         bucket,
			FilenamePrefix: "log_events_",
			FilenameSuffix: ".avro",
			OutputFormat: &pubsubpb.CloudStorageConfig_AvroConfig_{
				AvroConfig: &pubsubpb.CloudStorageConfig_AvroConfig{
					WriteMetadata: true,
				},
			},
			MaxDuration: durationpb.New(1 * time.Minute),
			MaxBytes:    1e8,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create cloud storage sub: %w", err)
	}
	fmt.Fprintf(w, "Created Cloud Storage subscription: %v\n", sub)

	return nil
}

// [END pubsub_create_cloud_storage_subscription]
