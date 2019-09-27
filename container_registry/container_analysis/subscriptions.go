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

package main

// [START containeranalysis_pubsub]

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	pubsub "cloud.google.com/go/pubsub"
)

// occurrencePubsub handles incoming Occurrences using a Cloud Pub/Sub subscription.
func occurrencePubsub(w io.Writer, subscriptionID string, timeout time.Duration, projectID string) (int, error) {
	// subscriptionID := fmt.Sprintf("my-occurrences-subscription")
	// timeout := time.Duration(20) * time.Second
	ctx := context.Background()

	var mu sync.Mutex
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return -1, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	// Subscribe to the requested Pub/Sub channel.
	sub := client.Subscription(subscriptionID)
	count := 0

	// Listen to messages for 'timeout' seconds.
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		count = count + 1
		fmt.Fprintf(w, "Message %d: %q\n", count, string(msg.Data))
		msg.Ack()
		mu.Unlock()
	})
	if err != nil {
		return -1, fmt.Errorf("sub.Receive: %v", err)
	}
	// Print and return the number of Pub/Sub messages received.
	fmt.Fprintln(w, count)
	return count, nil
}

// createOccurrenceSubscription creates a new Pub/Sub subscription object listening to the Occurrence topic.
func createOccurrenceSubscription(subscriptionID, projectID string) error {
	// subscriptionID := fmt.Sprintf("my-occurrences-subscription")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	// This topic id will automatically receive messages when Occurrences are added or modified
	topicID := "container-analysis-occurrences-v1"
	topic := client.Topic(topicID)
	config := pubsub.SubscriptionConfig{Topic: topic}
	_, err = client.CreateSubscription(ctx, subscriptionID, config)
	return fmt.Errorf("client.CreateSubscription: %v", err)
}

// [END containeranalysis_pubsub]
