// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package main

import (
	"fmt"
	"log"
	"os"

	// [START imports]
	"golang.org/x/net/context"

	"cloud.google.com/go/pubsub"
	// [END imports]
)

func main() {
	ctx := context.Background()
	// [START auth]
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	// [END auth]

	// List all the topics from the project.
	fmt.Println("Listing all topics from the project:")
	for _, t := range list(client) {
		fmt.Printf("%v\n", t.Name())
	}

	const topic = "example-topic"
	// Create a new topic called example-topic.
	create(client, topic)

	// Publish a text message on the created topic.
	publish(client, topic, "hello world!")

	// Delete the topic.
	delete(client, topic)
}

func create(client *pubsub.Client, topic string) {
	ctx := context.Background()
	// [START create_topic]
	t, err := client.NewTopic(ctx, topic)
	if err != nil {
		log.Fatalf("Could not create a new topic: %v", err)
	}
	fmt.Printf("Topic created: %v\n", t.Name())
	// [END create_topic]
}

func list(client *pubsub.Client) []*pubsub.Topic {
	ctx := context.Background()

	// [START list_topics]
	var topics []*pubsub.Topic

	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == pubsub.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		topics = append(topics, topic)
	}

	return topics
	// [END list_topics]
}

func delete(client *pubsub.Client, topic string) {
	ctx := context.Background()
	// [START delete_topic]
	t := client.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted topic: %v\n", t.Name())
	// [END delete_topic]
}

func publish(client *pubsub.Client, topic, msg string) {
	ctx := context.Background()
	// [START publish]
	t := client.Topic(topic)
	msgIDs, err := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	if err != nil {
		log.Fatalf("Failed to publish the message: %v", err)
	}
	for _, id := range msgIDs {
		fmt.Printf("Published a message; msg ID: %v\n", id)
	}
	// [END publish]
}
