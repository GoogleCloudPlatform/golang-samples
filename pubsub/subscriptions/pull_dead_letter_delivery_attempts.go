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

// [START pubsub_subscriber_dead_letter_delivery_attempt]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

func pullMsgsDeadLetterDeliveryAttempt(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	sub := client.Subscription(subID)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// The delivery attempt is an int pointer with the number of times
		// a message has attempted to be delivered if dead lettering is enabled.
		// Otherwise, it is nil.
		fmt.Fprintf(w, "message: %s, delivery attempts: %d", msg.Data, *msg.DeliveryAttempt)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("receive: %v", err)
	}
	return nil
}

// [END pubsub_subscriber_dead_letter_delivery_attempt]
