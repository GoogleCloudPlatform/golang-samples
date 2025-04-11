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

package subscriptions

// [START pubsub_create_subscription_with_exactly_once_delivery]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func createSubscriptionWithExactlyOnceDelivery(w io.Writer, projectID, subID string, topic *pubsub.Topic) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:                     topic,
		EnableExactlyOnceDelivery: true,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Created a subscription with exactly once delivery enabled: %v\n", sub)
	return nil
}

// [END pubsub_create_subscription_with_exactly_once_delivery]
