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

// [START pubsub_old_version_optimistic_subscribe]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// optimisticSubscribe shows the recommended pattern for optimistically
// assuming a subscription exists prior to receiving messages.
func optimisticSubscribe(w io.Writer, projectID, topicID, subID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Instead of checking if the subscription exists, optimistically try to
	// receive from the subscription.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got from existing subscription: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				// Since the subscription does not exist, create the subscription.
				s, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
					Topic: client.Topic(topicID),
				})
				if err != nil {
					return err
				}
				fmt.Fprintf(w, "Created subscription: %q\n", subID)

				// Pull from the new subscription.
				err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
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

// [END pubsub_old_version_optimistic_subscribe]
