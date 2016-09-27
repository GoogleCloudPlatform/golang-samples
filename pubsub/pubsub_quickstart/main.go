// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START pubsub_quickstart]
// Sample pubsub_quickstart creates a Google Cloud Pub/Sub topic.
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"

	// Imports the Google Cloud Pub/Sub client package
	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	// Your Google Cloud Platform project ID
	projectID := "YOUR_PROJECT_ID"

	// Creates a client
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// The name for the new topic
	topicName := "my-new-topic"

	// Creates the new topic
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	fmt.Printf("Topic %v created.", topic)
}

// [END pubsub_quickstart]
