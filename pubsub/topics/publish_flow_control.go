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

package topics

// [START pubsub_publish_flow_control]
import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"

	"cloud.google.com/go/pubsub"
)

func publish_with_flow_control_settings(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	t.PublishSettings.FlowControlSettings = pubsub.FlowControlSettings{
		MaxOutstandingMessages: 100,                           // default 1000
		MaxOutstandingBytes:    10 * 1024 * 1024,              // default 0 (unlimited)
		LimitExceededBehavior:  pubsub.FlowControlSignalError, // default FlowControlBlock
	}

	var wg sync.WaitGroup
	numMsgs := 1000
	// Rapidly publishing 1000 messages in a loop may be constrained by flow control.
	for i := 0; i < numMsgs; i++ {
		wg.Add(1)
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("message #" + strconv.Itoa(i)),
		})
		go func() {
			if _, err := result.Get(ctx); err != nil {
				fmt.Fprintf(w, "result.Get(): %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Fprintf(w, "Published %d messages with flow control enabled", numMsgs)
	return nil
}

// [END pubsub_publish_flow_control]
