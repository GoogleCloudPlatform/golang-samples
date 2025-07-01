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

// [START pubsub_dead_letter_delivery_attempt]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub/v2"
)

func pullMsgsDeadLetterDeliveryAttempt(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// client.Subscriber can be passed a subscription ID (e.g. "my-sub") or
	// a fully qualified name (e.g. "projects/my-project/subscriptions/my-sub").
	// If a subscription ID is provided, the project ID from the client is used.
	sub := client.Subscriber(subID)
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		// When dead lettering is enabled, the delivery attempt field is a pointer to the
		// the number of times the service has attempted to delivery a message.
		// Otherwise, the field is nil.
		if msg.DeliveryAttempt != nil {
			fmt.Fprintf(w, "message: %s, delivery attempts: %d", msg.Data, *msg.DeliveryAttempt)
		}
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("got error in Receive: %w", err)
	}
	return nil
}

// [END pubsub_dead_letter_delivery_attempt]
