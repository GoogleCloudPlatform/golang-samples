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
	topics, err := list(client)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}
	for _, t := range topics {
		fmt.Printf("%v\n", t.Name())
	}

	const topic = "example-topic"
	// Create a new topic called example-topic.
	if err := create(client, topic); err != nil {
		log.Fatalf("Failed to create a topic: %v", err)
	}

	// Publish a text message on the created topic.
	if err := publish(client, topic, "hello world!"); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	// Delete the topic.
	if err := delete(client, topic); err != nil {
		log.Fatalf("Failed to delete the topic: %v", err)
	}
}

func create(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START create_topic]
	t, err := client.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	fmt.Printf("Topic created: %v\n", t.Name())
	// [END create_topic]
	return nil
}

func list(client *pubsub.Client) ([]*pubsub.Topic, error) {
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
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
	// [END list_topics]
}

func delete(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START delete_topic]
	t := client.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		return err
	}
	fmt.Printf("Deleted topic: %v\n", t.Name())
	// [END delete_topic]
	return nil
}

func publish(client *pubsub.Client, topic, msg string) error {
	ctx := context.Background()
	// [START publish]
	t := client.Topic(topic)
	msgIDs, err := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	if err != nil {
		return err
	}
	for _, id := range msgIDs {
		fmt.Printf("Published a message; msg ID: %v\n", id)
	}
	// [END publish]
	return nil
}
