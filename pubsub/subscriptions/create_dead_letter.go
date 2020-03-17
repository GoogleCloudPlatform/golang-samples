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

// [START pubsub_create_subscription_dead_letter]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

func createSubWithDeadLetter(w io.Writer, projectID, subID string, topicID string, deadLetterTopicName string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topicID := "my-topic"
	// deadLetterTopicName := "projects/my-project/topics/my-dead-letter-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicID)

	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     deadLetterTopicName,
			MaxDeliveryAttempts: 10,
		},
	})
	if err != nil {
		return fmt.Errorf("CreateSubscription: %v", err)
	}
	fmt.Fprintf(w, "Created subscription (%s) with dead letter topic (%s)\n", sub.String(), deadLetterTopicName)
	return nil
}

// [END pubsub_create_subscription_dead_letter]
