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

// [START pubsub_subscriber_sync_pull]
import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/pubsub"
)

func pullMsgsSync(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
	// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
	// the max number of messages the client will hold in memory at a time.
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	// Receive messages for 5 seconds.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			fmt.Fprintf(w, "Got message :%q\n", string(msg.Data))
			msg.Ack()
		}
	}()

	// Receive blocks until the passed in context is done.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil && status.Code(err) != codes.Canceled {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}

// [END pubsub_subscriber_sync_pull]
