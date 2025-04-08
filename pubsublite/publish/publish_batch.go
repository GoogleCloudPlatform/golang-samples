// Copyright 2021 Google LLC
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

package publish

// [START pubsublite_publish_batch]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite/pscompat"
)

func publishWithBatchSettings(w io.Writer, projectID, zone, topicID string, messageCount int) error {
	// projectID := "my-project-id"
	// zone := "us-central1-a"
	// topicID := "my-topic"
	// messageCount := 10
	ctx := context.Background()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, zone, topicID)

	// Batch settings control how the publisher batches messages. These settings
	// apply per partition.
	// See https://pkg.go.dev/cloud.google.com/go/pubsublite/pscompat#pkg-variables
	// for DefaultPublishSettings.
	settings := pscompat.PublishSettings{
		ByteThreshold:  5 * 1024, // 5 KiB
		CountThreshold: 1000,     // 1,000 messages
		DelayThreshold: 100 * time.Millisecond,
	}

	// Create the publisher client.
	publisher, err := pscompat.NewPublisherClientWithSettings(ctx, topicPath, settings)
	if err != nil {
		return fmt.Errorf("pscompat.NewPublisherClientWithSettings error: %w", err)
	}

	// Ensure the publisher will be shut down.
	defer publisher.Stop()

	// Publish requests are sent to the server based on request size, message
	// count and time since last publish, whichever condition is met first.
	var results []*pubsub.PublishResult
	for i := 0; i < messageCount; i++ {
		r := publisher.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("message-%d", i)),
		})
		results = append(results, r)
	}

	// Print publish results.
	var publishedCount int
	for _, r := range results {
		// Get blocks until the result is ready.
		id, err := r.Get(ctx)
		if err != nil {
			// NOTE: A failed PublishResult indicates that the publisher client
			// encountered a fatal error and has permanently terminated. After the
			// fatal error has been resolved, a new publisher client instance must be
			// created to republish failed messages.
			fmt.Fprintf(w, "Publish error: %v\n", err)
			continue
		}
		fmt.Fprintf(w, "Published: %v\n", id)
		publishedCount++
	}

	fmt.Fprintf(w, "Published %d messages with batch settings\n", publishedCount)
	return publisher.Error()
}

// [END pubsublite_publish_batch]
