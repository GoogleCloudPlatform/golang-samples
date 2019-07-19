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

// Command topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
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

func create(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START pubsub_create_topic]
	t, err := client.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	fmt.Printf("Topic created: %v\n", t)
	// [END pubsub_create_topic]
	return nil
}

func list(client *pubsub.Client) ([]*pubsub.Topic, error) {
	ctx := context.Background()

	// [START pubsub_list_topics]
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
	// [END pubsub_list_topics]
}

func listSubscriptions(client *pubsub.Client, topicID string) ([]*pubsub.Subscription, error) {
	ctx := context.Background()

	// [START pubsub_list_topic_subscriptions]
	var subs []*pubsub.Subscription

	it := client.Topic(topicID).Subscriptions(ctx)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	// [END pubsub_list_topic_subscriptions]
	return subs, nil
}

func delete(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START pubsub_delete_topic]
	t := client.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		return err
	}
	fmt.Printf("Deleted topic: %v\n", t)
	// [END pubsub_delete_topic]
	return nil
}

func publish(client *pubsub.Client, topic, msg string) error {
	ctx := context.Background()
	// [START pubsub_publish]
	// [START pubsub_quickstart_publisher]
	t := client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END pubsub_publish]
	// [END pubsub_quickstart_publisher]
	return nil
}

func publishThatScales(client *pubsub.Client, topic string, n int) error {
	ctx := context.Background()
	// [START pubsub_publish_with_error_handling_that_scales]
	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topic)

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				log.Output(1, fmt.Sprintf("Failed to publish: %v", err))
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Printf("Published message %d; msg ID: %v\n", i, id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return errors.New(
			fmt.Sprintf("%d of %d messages did not publish successfully",
				totalErrors, n))
	}
	return nil
	// [END pubsub_publish_with_error_handling_that_scales]
}

func publishCustomAttributes(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START pubsub_publish_custom_attributes]
	t := client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("Hello world!"),
		Attributes: map[string]string{
			"origin":   "golang",
			"username": "gcp",
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published message with custom attributes; msg ID: %v\n", id)
	// [END pubsub_publish_custom_attributes]
	return nil
}

func publishWithSettings(client *pubsub.Client, topic string, msg []byte) error {
	ctx := context.Background()
	// [START pubsub_publisher_batch_settings]
	t := client.Topic(topic)
	t.PublishSettings.ByteThreshold = 5000
	t.PublishSettings.CountThreshold = 10
	t.PublishSettings.DelayThreshold = 100 * time.Millisecond

	result := t.Publish(ctx, &pubsub.Message{Data: msg})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END pubsub_publisher_batch_settings]
	return nil
}

func publishSingleGoroutine(client *pubsub.Client, topic string, msg []byte) error {
	ctx := context.Background()
	// [START pubsub_publisher_concurrency_control]
	t := client.Topic(topic)
	t.PublishSettings.NumGoroutines = 1

	result := t.Publish(ctx, &pubsub.Message{Data: msg})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END pubsub_publisher_concurrency_control]
	return nil
}

func getPolicy(c *pubsub.Client, topicName string) (*iam.Policy, error) {
	ctx := context.Background()

	// [START pubsub_get_topic_policy]
	policy, err := c.Topic(topicName).IAM().Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, role := range policy.Roles() {
		log.Print(policy.Members(role))
	}
	// [END pubsub_get_topic_policy]
	return policy, nil
}

func addUsers(c *pubsub.Client, topicName string) error {
	ctx := context.Background()

	// [START pubsub_set_topic_policy]
	topic := c.Topic(topicName)
	policy, err := topic.IAM().Policy(ctx)
	if err != nil {
		return err
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
	return nil
}

func testPermissions(c *pubsub.Client, topicName string) ([]string, error) {
	ctx := context.Background()

	// [START pubsub_test_topic_permissions]
	topic := c.Topic(topicName)
	perms, err := topic.IAM().TestPermissions(ctx, []string{
		"pubsub.topics.publish",
		"pubsub.topics.update",
	})
	if err != nil {
		return nil, err
	}
	for _, perm := range perms {
		log.Printf("Allowed: %v", perm)
	}
	// [END pubsub_test_topic_permissions]
	return perms, nil
}
