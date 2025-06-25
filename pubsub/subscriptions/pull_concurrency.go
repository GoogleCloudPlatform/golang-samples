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

	"cloud.google.com/go/pubsub/v2"
)

func pullMsgsConcurrencyControl(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	// client.Subscriber can be passed a subscription ID (e.g. "my-sub") or
	// a fully qualified name (e.g. "projects/my-project/subscriptions/my-sub").
	// If a subscription ID is provided, the project ID from the client is used.
	sub := client.Subscriber(subID)
	// NumGoroutines determines the number of streams sub.Receive will spawn to pull
	// messages. It is recommended to set this to 1, unless your throughput
	// is greater than 10 MB/s, as even having 1 stream can still result in
	// messages being handled asynchronously.
	sub.ReceiveSettings.NumGoroutines = 1
	// MaxOutstandingMessages limits the number of concurrent handlers of messages.
	// In this case, up to 8 unacked messages can be handled concurrently.
	sub.ReceiveSettings.MaxOutstandingMessages = 8

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received int32

	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive returned error: %w", err)
	}
	fmt.Fprintf(w, "Received %d messages\n", received)

	return nil
}

// [END pubsub_subscriber_concurrency_control]
