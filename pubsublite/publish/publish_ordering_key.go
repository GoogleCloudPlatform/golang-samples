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

// [START pubsublite_publish_ordering_key]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite/pscompat"
)

func publishWithOrderingKey(w io.Writer, projectID, zone, topicID string, messageCount int) error {
	// projectID := "my-project-id"
	// zone := "us-central1-a"
	// topicID := "my-topic"
	// messageCount := 10
	ctx := context.Background()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, zone, topicID)

	// Create the publisher client.
	publisher, err := pscompat.NewPublisherClient(ctx, topicPath)
	if err != nil {
		return fmt.Errorf("pscompat.NewPublisherClient error: %w", err)
	}

	// Ensure the publisher will be shut down.
	defer publisher.Stop()

	// Messages of the same ordering key will always get published to the same
	// partition. When OrderingKey is unset, messages can get published to
	// different partitions if more than one partition exists for the topic.
	var results []*pubsub.PublishResult
	for i := 0; i < messageCount; i++ {
		r := publisher.Publish(ctx, &pubsub.Message{
			OrderingKey: "test_ordering_key",
			Data:        []byte(fmt.Sprintf("message-%d", i)),
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

		// Metadata decoded from the id contains the partition and offset.
		metadata, err := pscompat.ParseMessageMetadata(id)
		if err != nil {
			return fmt.Errorf("failed to parse message metadata %q: %w", id, err)
		}
		fmt.Fprintf(w, "Published: partition=%d, offset=%d\n", metadata.Partition, metadata.Offset)
		publishedCount++
	}

	fmt.Fprintf(w, "Published %d messages with ordering key\n", publishedCount)
	return publisher.Error()
}

// [END pubsublite_publish_ordering_key]
