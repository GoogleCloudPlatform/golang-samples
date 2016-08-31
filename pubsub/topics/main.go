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
	proj := os.Getenv("GCLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GCLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	c, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		log.Fatal(err)
	}
	// [END auth]

	// List all the topics from the project.
	list(c)

	const topic = "example-topic"
	// Create a new topic called example-topic.
	create(c, topic)

	// Publish a text message on the created topic.
	publish(c, topic, "hello world!")

	// Delete the topic.
	delete(c, topic)
}

func create(c *pubsub.Client, topic string) {
	ctx := context.Background()
	// [START create_topic]
	t, err := c.NewTopic(ctx, topic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("topic created: %v\n", t.Name())
	// [END create_topic]
}

func list(c *pubsub.Client) {
	// [START list_topics]
	ctx := context.Background()
	it := c.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == pubsub.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(topic.Name())
	}
	// [END list_topics]
}

func delete(c *pubsub.Client, topic string) {
	ctx := context.Background()
	// [START delete_topic]
	t := c.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted topic: %v\n", t.Name())
	// [END delete_topic]
}

func publish(c *pubsub.Client, topic, msg string) {
	ctx := context.Background()
	// [START publish]
	t := c.Topic(topic)
	msgIDs, err := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	if err != nil {
		log.Fatalf("failed to publish the message: %v", err)
	}
	for _, id := range msgIDs {
		fmt.Printf("published a message; msg ID: %v\n", id)
	}
	// [END publish]
}
