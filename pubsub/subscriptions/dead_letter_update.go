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

// [START pubsub_dead_letter_update_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateDeadLetter updates an existing subscription with a dead letter policy.
func updateDeadLetter(w io.Writer, projectID, subscription, deadLetterTopic string) error {
	// projectID := "my-project-id"
	// subID := "projects/my-project-id/subscriptions/my-sub"
	// deadLetterTopic := "projects/my-project-id/topics/my-dead-letter-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub, err := client.SubscriptionAdminClient.UpdateSubscription(ctx, &pubsubpb.UpdateSubscriptionRequest{
		Subscription: &pubsubpb.Subscription{
			Name: subscription,
			DeadLetterPolicy: &pubsubpb.DeadLetterPolicy{
				MaxDeliveryAttempts: 20,
				DeadLetterTopic:     deadLetterTopic,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"dead_letter_policy"},
		},
	})
	if err != nil {
		return fmt.Errorf("UpdateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Updated subscription: %+v\n", sub)
	return nil
}

// [END pubsub_dead_letter_update_subscription]
