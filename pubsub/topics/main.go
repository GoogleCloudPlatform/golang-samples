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

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
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
		fmt.Println(t)
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
	fmt.Printf("Topic created: %v\n", t)
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
		if err == iterator.Done {
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
	fmt.Printf("Deleted topic: %v\n", t)
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

func getPolicy(c *pubsub.Client, topicName string) *iam.Policy {
	ctx := context.Background()

	// [START pubsub_get_topic_policy]
	policy, err := c.Topic(topicName).IAM().Policy(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, role := range policy.Roles() {
		log.Print(policy.Members(role))
	}
	// [END pubsub_get_topic_policy]
	return policy
}

func addUsers(c *pubsub.Client, topicName string) {
	ctx := context.Background()

	// [START pubsub_set_topic_policy]
	topic := c.Topic(topicName)
	policy, err := topic.IAM().Policy(ctx)
	if err != nil {
		log.Fatalf("GetPolicy: %v", err)
	}
	// Other valid prefixes are "serviceAccount:", "user:"
	// See the documentation for more values.
	policy.Add(iam.AllUsers, iam.Viewer)
	policy.Add("group:cloud-logs@google.com", iam.Editor)
	if err := topic.IAM().SetPolicy(ctx, policy); err != nil {
		log.Fatalf("SetPolicy: %v", err)
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END pubsub_set_topic_policy]
}

func testPermissions(c *pubsub.Client, topicName string) []string {
	ctx := context.Background()

	// [START pubsub_test_topic_permissions]
	topic := c.Topic(topicName)
	perms, err := topic.IAM().TestPermissions(ctx, []string{
		"pubsub.topics.publish",
		"pubsub.topics.update",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, perm := range perms {
		log.Printf("Allowed: %v", perm)
	}
	// [END pubsub_test_topic_permissions]
	return perms
}
