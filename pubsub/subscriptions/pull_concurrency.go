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

	"cloud.google.com/go/pubsub"
)

func pullMsgsConcurrent(w io.Writer, projectID, subName string) error {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	sub := client.Subscription(subName)
	// NumGoroutines is the number of goroutines Receive will spawn to pull messages concurrently.
	sub.ReceiveSettings.NumGoroutines = 4
	// If ReceiveSettings.Synchronous is set to true, NumGoroutines is overriden to 1. To enable
	// concurrency settings, set this to false.
	sub.ReceiveSettings.Synchronous = false
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

// [END pubsub_subscriber_concurrency_control]
