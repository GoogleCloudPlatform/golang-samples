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
	var publishedCount int
	for _, r := range results {
		// Get blocks until the result is ready.
		id, err := r.Get(ctx)
		if err != nil {
			// NOTE: A failed PublishResult indicates that the publisher client
			// encountered a fatal error and has permanently terminated. After the
			// fatal error has been resolved, a new publisher client instance must be
			// created to republish failed messages.
			fmt.Printf("Publish error: %v\n", err)
			continue
		}

		// Metadata decoded from the id contains the partition and offset.
		metadata, err := pscompat.ParseMessageMetadata(id)
		if err != nil {
			log.Fatalf("Failed to parse message metadata %q: %v", id, err)
		}

		fmt.Printf("Published: partition=%d, offset=%d\n", metadata.Partition, metadata.Offset)
		publishedCount++
	}

	fmt.Printf("Published %d messages\n", publishedCount)

	if err := publisher.Error(); err != nil {
		fmt.Printf("Publisher client terminated due to error: %v\n", publisher.Error())
	}
}

// [END pubsublite_quickstart_publisher]
