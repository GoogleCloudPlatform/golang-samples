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

// [START pubsub_publish_with_error_handler]
import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub/v2"
)

func publishThatScales(w io.Writer, projectID, topicID string, n int) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	var totalErrors uint64

	// Make sure to reuse this publisher across publishes.
	p := client.Publisher(topicID)

	for i := 0; i < n; i++ {
		result := p.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Fprintf(w, "Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Fprintf(w, "Published message %d; msg ID: %v\n", i, id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}

// [END pubsub_publish_with_error_handler]
