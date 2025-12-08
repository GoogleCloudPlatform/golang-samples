// Copyright 2020 Google LLC
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

// [START pubsub_dead_letter_create_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

// createSubWithDeadLetter creates a subscription with a dead letter policy.
func createSubWithDeadLetter(w io.Writer, projectID, topic, subscription, fullyQualifiedDeadLetterTopic string) error {
	// projectID := "my-project-id"
	// topic := "projects/my-project-id/topics/my-topic"
	// subscription := "projects/my-project-id/subscriptions/my-sub"
	// fullyQualifiedDeadLetterTopic := "projects/my-project-id/topics/my-dead-letter-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subscription,
		Topic: topic,
		DeadLetterPolicy: &pubsubpb.DeadLetterPolicy{
			DeadLetterTopic:     fullyQualifiedDeadLetterTopic,
			MaxDeliveryAttempts: 10,
		},
	})
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Created subscription with dead letter topic: (%s)\n", fullyQualifiedDeadLetterTopic)
	fmt.Fprintln(w, "To process dead letter messages, remember to add a subscription to your dead letter topic.")
	return nil
}

// [END pubsub_dead_letter_create_subscription]
