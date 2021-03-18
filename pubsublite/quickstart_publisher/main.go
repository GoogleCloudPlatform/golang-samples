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

// [START pubsublite_quickstart_publisher]

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite/pscompat"
)

func main() {
	// NOTE: Set these flags for an existing Pub/Sub Lite topic when running this
	// sample.
	projectID := flag.String("project_id", "", "Cloud Project ID")
	zone := flag.String("zone", "", "Cloud Zone where the topic resides, e.g. us-central1-a")
	topicID := flag.String("topic_id", "", "Existing Pub/Sub Lite topic")
	messageCount := flag.Int("message_count", 100, "The number of messages to send")
	flag.Parse()

	ctx := context.Background()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", *projectID, *zone, *topicID)

	// Create the publisher client.
	publisher, err := pscompat.NewPublisherClient(ctx, topicPath)
	if err != nil {
		log.Fatalf("pscompat.NewPublisherClient error: %v", err)
	}

	// Ensure the publisher will be shut down.
	defer publisher.Stop()

	// Publish messages. Messages are automatically batched.
	var results []*pubsub.PublishResult
	for i := 0; i < *messageCount; i++ {
		r := publisher.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("message-%d", i)),
		})
		results = append(results, r)
	}

	// Print publish results.
	for _, r := range results {
		// Get blocks until the result is ready.
		id, err := r.Get(ctx)
		if err != nil {
			// NOTE: The publisher will terminate upon first error. Create a new
			// publisher to republish failed messages.
			log.Fatalf("Publish error: %v", err)
		}

		// Metadata decoded from the id contains the partition and offset.
		metadata, err := pscompat.ParseMessageMetadata(id)
		if err != nil {
			log.Fatalf("Failed to parse %q: %v", id, err)
		}
		fmt.Printf("Published: partition=%d, offset=%d\n", metadata.Partition, metadata.Offset)
	}

	fmt.Printf("Published %d messages\n", *messageCount)
}

// [END pubsublite_quickstart_publisher]
