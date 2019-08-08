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

// Package topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package topics

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	// List all the topics from the project.
	fmt.Println("Listing all topics from the project:")
	topics, err := list(client)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}
	for _, t := range topics {
		fmt.Println(t)
	}

	const topic = "my-topic"
	// Create a new topic called my-topic.
	if err := create(client, topic); err != nil {
		log.Fatalf("Failed to create a topic: %v", err)
	}

	// Publish a text message on the created topic.
	if err := publish(client, topic, "hello world!"); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	// Publish 10 messages with asynchronous error handling.
	if err := publishThatScales(client, topic, 10); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	// Delete the topic.
	if err := delete(client, topic); err != nil {
		log.Fatalf("Failed to delete the topic: %v", err)
	}
}
