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

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

func createSubscriptionWithExactlyOnceDelivery(w io.Writer, projectID, topic, subscription string) error {
	// projectID := "my-project-id"
	// topic := "projects/my-project-id/topics/my-topic"
	// subscription := "projects/my-project/subscriptions/my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subscription,
		Topic: topic, EnableExactlyOnceDelivery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create exactly once sub: %w", err)
	}
	fmt.Fprintf(w, "Created a subscription with exactly once delivery enabled: %v\n", sub)
	return nil
}

// [END pubsub_create_subscription_with_exactly_once_delivery]
