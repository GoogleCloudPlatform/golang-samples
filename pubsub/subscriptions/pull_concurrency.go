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
	"runtime"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
)

func pullMsgsConcurrenyControl(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	// Must set ReceiveSettings.Synchronous to false (or leave as default) to enable
	// concurrency settings. Otherwise, NumGoroutines will be set to 1.
	sub.ReceiveSettings.Synchronous = false
	// NumGoroutines is the number of goroutines sub.Receive will spawn to pull
	// messages concurrently.
	sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var counter int32

	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		// The message handler passed to Receive may be called concurrently
		// so it's okay to process the messages concurrently but make sure
		// to synchronize access to shared memory.
		atomic.AddInt32(&counter, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("pubsub: Receive returned error: %v", err)
	}
	fmt.Fprintf(w, "Received %d messages\n", counter)

	return nil
}

// [END pubsub_subscriber_concurrency_control]
