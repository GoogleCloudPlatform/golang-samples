// Copyright 2025 Google LLC
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

package topics

// [START pubsub_old_version_publisher_flow_control]
import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

func publishWithFlowControlSettings(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	t.PublishSettings.FlowControlSettings = pubsub.FlowControlSettings{
		MaxOutstandingMessages: 100,                     // default 1000
		MaxOutstandingBytes:    10 * 1024 * 1024,        // default 0 (unlimited)
		LimitExceededBehavior:  pubsub.FlowControlBlock, // default Ignore, other options: Block and SignalError
	}

	var wg sync.WaitGroup
	var totalErrors uint64

	numMsgs := 1000
	// Rapidly publishing 1000 messages in a loop may be constrained by flow control.
	for i := 0; i < numMsgs; i++ {
		wg.Add(1)
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("message #" + strconv.Itoa(i)),
		})
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			_, err := res.Get(ctx)
			if err != nil {
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, numMsgs)
	}
	fmt.Fprintf(w, "finished publishing %d messages", numMsgs)
	return nil
}

// [END pubsub_old_version_publisher_flow_control]
