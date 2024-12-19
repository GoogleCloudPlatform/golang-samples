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

// [START pubsub_old_version_dead_letter_remove]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

// removeDeadLetterTopic removes the dead letter policy from a subscription.
func removeDeadLetterTopic(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	subConfig, err := client.Subscription(subID).Update(ctx, pubsub.SubscriptionConfigToUpdate{
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{},
	})
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	fmt.Fprintf(w, "Updated subscription config: %+v\n", subConfig)
	return nil
}

// [END pubsub_old_version_dead_letter_remove]
