// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command subscriptions is a tool to manage Google Cloud Pub/Sub subscriptions by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	// Print all the subscriptions in the project.
	fmt.Println("Listing all subscriptions from the project:")
	for _, sub := range list(client) {
		fmt.Printf("%v\n", sub.Name())
	}

	const topic = "example-topic"
	// Create a topic to subscribe to.
	t, err := client.NewTopic(ctx, topic)
	if err != nil {
		log.Fatalf("Failed to create the topic: %v", err)
	}
	defer t.Delete(ctx) // cleanup when finished using t.

	const sub = "example-subscription"
	// Create a new subscription.
	create(client, sub, t)

	// Pull messages via the subscription.
	pullMsgs(client, sub, t)

	// Delete the subscription.
	delete(client, sub)
}

func list(c *pubsub.Client) []*pubsub.Subscription {
	ctx := context.Background()
	// [START get_all_subscriptions]
	var subs []*pubsub.Subscription
	it := c.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == pubsub.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list subscriptions: %v", err)
		}
		subs = append(subs, s)
	}
	return subs
	// [END get_all_subscriptions]
}

func pullMsgs(c *pubsub.Client, name string, topic *pubsub.Topic) {
	ctx := context.Background()

	go func() {
		// publish 10 messages on the topic.
		for i := 0; i < 10; i++ {
			_, err := topic.Publish(ctx, &pubsub.Message{
				Data: []byte(fmt.Sprintf("hello world #%d", i)),
			})
			if err != nil {
				log.Fatalf("Failed to publish message #%d: %v", i, err)
			}
		}
	}()

	// [START pull_messages]
	sub := c.Subscription(name)
	it, err := sub.Pull(ctx)
	if err != nil {
		log.Fatalf("Failed to pull from subsription: %v", err)
	}
	defer it.Stop()

	// Consume 10 messages.
	for i := 0; i < 10; i++ {
		msg, err := it.Next()
		if err == pubsub.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed when iterating on messages: %v", err)
		}
		fmt.Printf("Got message: %q\n", string(msg.Data))
		msg.Done(true)
	}
	// [END pull_messages]
}

func create(c *pubsub.Client, name string, topic *pubsub.Topic) {
	ctx := context.Background()
	// [START create_subscription]
	sub, err := c.NewSubscription(ctx, name, topic, 10*time.Second, nil)
	if err != nil {
		log.Fatalf("Failed to create a new subscription: %v", err)
	}
	fmt.Printf("Created subscription: %v\n", sub)
	// [END create_subscription]
}

func delete(c *pubsub.Client, name string) {
	ctx := context.Background()
	// [START delete_subscription]
	sub := c.Subscription(name)
	if err := sub.Delete(ctx); err != nil {
		log.Fatalf("Failed to delete the subscription: %v\n", err)
	}
	fmt.Println("Subscription deleted.")
	// [END delete_subscription]
}
