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

// Package subscription is a tool to manage Google Cloud Pub/Sub subscriptions by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package subscription

import (
	"context"
	"fmt"
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
		fmt.Errorf("pubsub.NewClient: %v", err)
		return
	}

	// Print all the subscriptions in the project.
	fmt.Println("Listing all subscriptions from the project:")
	subs, err := list(client)
	if err != nil {
		fmt.Errorf("list: %v", err)
		return
	}
	for _, sub := range subs {
		fmt.Println(sub)
	}

	t, err := createTopicIfNotExists(client)
	if err != nil {
		fmt.Errorf("createTopicIfNotExists: %v", err)
		return
	}

	const sub = "my-sub"
	// Create a new subscription.
	if err := create(client, sub, t); err != nil {
		fmt.Errorf("create: %v", err)
		return
	}

	// Pull messages via the subscription.
	if err := pullMsgs(client, sub, t); err != nil {
		fmt.Errorf("pullMsgs: %v", err)
		return
	}

	// Delete the subscription.
	if err := delete(client, sub); err != nil {
		fmt.Errorf("delete: %v", err)
		return
	}
}

func createTopicIfNotExists(c *pubsub.Client) (*pubsub.Topic, error) {
	ctx := context.Background()

	const topic = "my-topic"
	// Create a topic to subscribe to.
	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("Exists: %v", err)
	}
	if ok {
		return t, nil
	}

	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %v", err)
	}
	return t, nil
}
