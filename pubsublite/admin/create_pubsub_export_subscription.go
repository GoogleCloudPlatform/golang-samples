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

package admin

// [START pubsublite_create_pubsub_export_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

func createPubsubExportSubscription(w io.Writer, projectID, region, location, topicID, subID, pubsubTopicID string) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// NOTE: location can be either a region ("us-central1") or a zone ("us-central1-a")
	// For a list of valid locations, see https://cloud.google.com/pubsub/lite/docs/locations.
	// location := "us-central1"
	// NOTE: topic and subscription must be in the same region/zone (e.g. "us-central1-a")
	// topicID := "my-topic"
	// subID := "my-subscription"
	// pubsubTopicID := "destination-topic-id"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %w", err)
	}
	defer client.Close()

	// Initialize the subscription to the oldest retained messages for each
	// partition.
	targetLocation := pubsublite.AtTargetLocation(pubsublite.Beginning)

	sub, err := client.CreateSubscription(ctx, pubsublite.SubscriptionConfig{
		Name:                fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projectID, location, subID),
		Topic:               fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, location, topicID),
		DeliveryRequirement: pubsublite.DeliverImmediately, // Can also be DeliverAfterStored.
		// Configures an export subscription that writes messages to a Pub/Sub topic.
		ExportConfig: &pubsublite.ExportConfig{
			DesiredState: pubsublite.ExportActive, // Can also be ExportPaused.
			Destination: &pubsublite.PubSubDestinationConfig{
				Topic: fmt.Sprintf("projects/%s/topics/%s", projectID, pubsubTopicID),
			},
		},
	}, targetLocation)
	if err != nil {
		return fmt.Errorf("client.CreateSubscription got err: %w", err)
	}
	fmt.Fprintf(w, "Created export subscription: %s\n", sub.Name)
	return nil
}

// [END pubsublite_create_pubsub_export_subscription]
