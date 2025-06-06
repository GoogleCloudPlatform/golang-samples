// Copyright 2023 Google LLC
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

// [START pubsub_create_unwrapped_push_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

// createPushNoWrapperSubscription creates a push subscription where messages are delivered in the HTTP body.
func createPushNoWrapperSubscription(w io.Writer, projectID, topic, subscription, endpoint string) error {
	// projectID := "my-project-id"
	// topic := "projects/my-project-id/topics/my-topic"
	// subscription := "projects/my-project/subscriptions/my-sub"
	// endpoint := "https://my-test-project.appspot.com/push"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:               subscription,
		Topic:              topic,
		AckDeadlineSeconds: 10,
		PushConfig: &pubsubpb.PushConfig{
			PushEndpoint: endpoint,
			Wrapper: &pubsubpb.PushConfig_NoWrapper_{
				NoWrapper: &pubsubpb.PushConfig_NoWrapper{
					// Determines if message metadata is added to the HTTP headers of
					// the delivered message.
					WriteMetadata: true,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Created push no wrapper subscription: %v\n", sub)
	return nil
}

// [END pubsub_create_unwrapped_push_subscription]
