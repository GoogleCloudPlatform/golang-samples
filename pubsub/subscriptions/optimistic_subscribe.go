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
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// optimisticSubscribe shows the recommended pattern for checking if a subscription exists
// prior to starting the subscriber application.
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

	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				s, err := client.CreateSubscription(ctx2, subID, pubsub.SubscriptionConfig{
					Topic: client.Topic(topicID),
				})
				fmt.Fprintf(w, "Created subscription: %q\n", subID)
				if err != nil {
					s.Receive(ctx2, func(ctx context.Context, msg *pubsub.Message) {
						fmt.Fprintf(w, "Got : %q\n", string(msg.Data))
						msg.Ack()
					})
				}
			}
		}
	}

	return nil
}

// [END pubsub_optimistic_subscribe]
