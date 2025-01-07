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

// [START pubsub_old_version_dead_letter_update_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

// updateDeadLetter updates an existing subscription with a dead letter policy.
func updateDeadLetter(w io.Writer, projectID, subID string, fullyQualifiedDeadLetterTopic string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// fullyQualifiedDeadLetterTopic := "projects/my-project/topics/my-dead-letter-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	updateConfig := pubsub.SubscriptionConfigToUpdate{
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     fullyQualifiedDeadLetterTopic,
			MaxDeliveryAttempts: 20,
		},
	}

	subConfig, err := client.Subscription(subID).Update(ctx, updateConfig)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	fmt.Fprintf(w, "Updated subscription config: %+v\n", subConfig)
	return nil
}

// [END pubsub_old_version_dead_letter_update_subscription]
