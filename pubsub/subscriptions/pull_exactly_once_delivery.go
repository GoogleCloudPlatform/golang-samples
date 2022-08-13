// Copyright 2022 Google LLC
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

// [START pubsub_subscriber_exactly_once]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

// receiveMessagesWithExactlyOnceDeliveryEnabled instantiates a subscriber client.
// This differs from regular subscribing since you must call msg.AckWithResult()
// or msg.NackWithResult() instead of the regular Ack/Nack methods.
// When exactly once delivery is enabled on the subscription, the message is
// guaranteed to not be delivered again if the ack result succeeds.
func receiveMessagesWithExactlyOnceDeliveryEnabled(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	// Set MinExtensionPeriod high to avoid any unintentional
	// acknowledgment expirations (e.g. due to network events).
	// This can lead to high tail latency in case of client crashes.
	sub.ReceiveSettings.MinExtensionPeriod = 600 * time.Second

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		r := msg.AckWithResult()
		// Block until the result is returned and a pubsub.AcknowledgeStatus
		// is returned for the acked message.
		status, err := r.Get(ctx)
		if err != nil {
			fmt.Fprintf(w, "MessageID: %s failed when calling result.Get: %v", msg.ID, err)
		}

		switch status {
		case pubsub.AcknowledgeStatusSuccess:
			fmt.Fprintf(w, "Message successfully acked: %s", msg.ID)
		case pubsub.AcknowledgeStatusInvalidAckID:
			fmt.Fprintf(w, "Message failed to ack with response of Invalid. ID: %s", msg.ID)
		case pubsub.AcknowledgeStatusPermissionDenied:
			fmt.Fprintf(w, "Message failed to ack with response of Permission Denied. ID: %s", msg.ID)
		case pubsub.AcknowledgeStatusFailedPrecondition:
			fmt.Fprintf(w, "Message failed to ack with response of Failed Precondition. ID: %s", msg.ID)
		case pubsub.AcknowledgeStatusOther:
			fmt.Fprintf(w, "Message failed to ack with response of Other. ID: %s", msg.ID)
		default:
		}
	})
	if err != nil {
		return fmt.Errorf("got err from sub.Receive: %v", err)
	}
	return nil
}

// [END pubsub_subscriber_exactly_once]
