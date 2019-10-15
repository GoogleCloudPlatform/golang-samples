// Copyright 2019 Google LLC
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

// [START pubsub_subscriber_flow_settings]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func pullMsgsSettings(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	sub := client.Subscription(subID)
	sub.ReceiveSettings.Synchronous = true
	// MaxOutstandingMessages is the maximum number of unprocessed messages the
	// client will pull from the server before pausing.
	//
	// This is only guaranteed when ReceiveSettings.Synchronous is set to true.
	// When Synchronous is set to false, the StreamingPull RPC is used which
	// can pull a single large batch of messages at once that is greater than
	// MaxOustandingMessages before pausing. For more info, see
	// https://cloud.google.com/pubsub/docs/pull#streamingpull_dealing_with_large_backlogs_of_small_messages.
	sub.ReceiveSettings.MaxOutstandingMessages = 10
	// MaxOutstandingBytes is the maximum size of unprocessed messages,
	// that the client will pull from the server before pausing. Similar
	// to MaxOutstandingMessages, this may be exceeded with a large batch
	// of messages since we cannot control the size of a batch of messages
	// from the server (even with the synchronous Pull RPC).
	sub.ReceiveSettings.MaxOutstandingBytes = 1e10
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

// [END pubsub_subscriber_flow_settings]
