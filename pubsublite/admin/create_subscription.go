// Copyright 2021 Google LLC
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

// [START pubsublite_create_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

func createSubscription(w io.Writer, projectID, region, zone, topicID, subID string) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// zone := "us-central1-a"
	// NOTE: topic and subscription must be in the same zone (i.e. "us-central1-a")
	// topicID := "my-topic"
	// subID := "my-subscription"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %v", err)
	}
	defer client.Close()

	sub, err := client.CreateSubscription(ctx, pubsublite.SubscriptionConfig{
		Name:                fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projectID, zone, subID),
		Topic:               fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, zone, topicID),
		DeliveryRequirement: pubsublite.DeliverImmediately, // can also be DeliverAfterStored
	})
	if err != nil {
		return fmt.Errorf("client.CreateSubscription got err: %v", err)
	}
	fmt.Fprintf(w, "Created subscription: %s\n", sub.Name)
	return nil
}

// [END pubsublite_create_subscription]
