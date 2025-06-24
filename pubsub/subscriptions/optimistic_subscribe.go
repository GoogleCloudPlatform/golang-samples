// Copyright 2024 Google LLC
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

// [START pubsub_optimistic_subscribe]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// optimisticSubscribe shows the recommended pattern for optimistically
// assuming a subscription exists prior to receiving messages.
func optimisticSubscribe(w io.Writer, projectID, topic, subscriptionName string) error {
	// projectID := "my-project-id"
	// topic := "projects/my-project-id/topics/my-topic"
	// subscription := "projects/my-project/subscriptions/my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub := client.Subscriber(subscriptionName)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Instead of checking if the subscription exists, optimistically try to
	// receive from the subscription assuming it exists.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got from existing subscription: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				// If the subscription does not exist, then create the subscription.
				subscription, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
					Name:  subscriptionName,
					Topic: topic,
				})
				if err != nil {
					return err
				}
				fmt.Fprintf(w, "Created subscription: %q\n", subscriptionName)

				sub = client.Subscriber(subscription.GetName())
				err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
					fmt.Fprintf(w, "Got from new subscription: %q\n", string(msg.Data))
					msg.Ack()
				})
				if err != nil && !errors.Is(err, context.Canceled) {
					return err
				}
			}
		}
	}
	return nil
}

// [END pubsub_optimistic_subscribe]
