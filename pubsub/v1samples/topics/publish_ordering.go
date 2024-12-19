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

// [START pubsub_old_version_publish_with_ordering_keys]
import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func publishWithOrderingKey(w io.Writer, projectID, topicID string) {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()

	// Pub/Sub's ordered delivery guarantee only applies when publishes for an ordering key are in the same region.
	// For list of locational endpoints for Pub/Sub, see https://cloud.google.com/pubsub/docs/reference/service_apis_overview#list_of_locational_endpoints
	client, err := pubsub.NewClient(ctx, projectID,
		option.WithEndpoint("us-east1-pubsub.googleapis.com:443"))
	if err != nil {
		fmt.Fprintf(w, "pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)
	t.EnableMessageOrdering = true

	messages := []struct {
		message     string
		orderingKey string
	}{
		{
			message:     "message1",
			orderingKey: "key1",
		},
		{
			message:     "message2",
			orderingKey: "key2",
		},
		{
			message:     "message3",
			orderingKey: "key1",
		},
		{
			message:     "message4",
			orderingKey: "key2",
		},
	}

	for _, m := range messages {
		res := t.Publish(ctx, &pubsub.Message{
			Data:        []byte(m.message),
			OrderingKey: m.orderingKey,
		})

		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			_, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				fmt.Printf("Failed to publish: %s\n", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(res)
	}

	wg.Wait()

	if totalErrors > 0 {
		fmt.Fprintf(w, "%d of 4 messages did not publish successfully", totalErrors)
		return
	}

	fmt.Fprint(w, "Published 4 messages with ordering keys successfully\n")
}

// [END pubsub_old_version_publish_with_ordering_keys]
