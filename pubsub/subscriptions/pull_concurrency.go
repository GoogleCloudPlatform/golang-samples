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

// [START pubsub_subscriber_concurrency_control]
import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
)

func pullMsgsConcurrenyControl(w io.Writer, projectID, subName string, numGoroutines int) error {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	// numGoroutines := 4
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subName)
	// Must set ReceiveSettings.Synchronous to false to enable concurrency settings.
	// Otherwise, NumGoroutines will be set to 1.
	sub.ReceiveSettings.Synchronous = false
	// NumGoroutines is the number of goroutines sub.Receive will spawn to pull messages concurrently.
	sub.ReceiveSettings.NumGoroutines = numGoroutines

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var numMsgs uint64
	// Create a channel to send messages to as they come in.
	cm := make(chan *pubsub.Message)

	// Handle individual messages in a goroutine.
	go func() {
		for {
			select {
			case msg := <-cm:
				_ = msg // TODO: handle message
				atomic.AddUint64(&numMsgs, 1)
				msg.Ack()
			case <-ctx.Done():
				fmt.Fprintf(w, "Received %d messages\n", numMsgs)
				return
			}
		}
	}()

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil {
		return fmt.Errorf("Error in Receive: %v", err)
	}

	return nil
}

// [END pubsub_subscriber_concurrency_control]
