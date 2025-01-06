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

// [START pubsub_old_version_dead_letter_create_subscription]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

// createSubWithDeadLetter creates a subscription with a dead letter policy.
func createSubWithDeadLetter(w io.Writer, projectID, subID string, topicID string, fullyQualifiedDeadLetterTopic string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topicID := "my-topic"
	// fullyQualifiedDeadLetterTopic := "projects/my-project/topics/my-dead-letter-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	topic := client.Topic(topicID)

	subConfig := pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     fullyQualifiedDeadLetterTopic,
			MaxDeliveryAttempts: 10,
		},
	}

	sub, err := client.CreateSubscription(ctx, subID, subConfig)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Created subscription (%s) with dead letter topic (%s)\n", sub.String(), fullyQualifiedDeadLetterTopic)
	fmt.Fprintln(w, "To process dead letter messages, remember to add a subscription to your dead letter topic.")
	return nil
}

// [END pubsub_old_version_dead_letter_create_subscription]
