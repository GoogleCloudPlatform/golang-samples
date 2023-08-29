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

package risk

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

// setupPubSub creates a subscription to the given topic.
func setupPubSub(ctx context.Context, projectID, topic, sub string) (*pubsub.Subscription, error) {
	// ctx := context.Background()
	log.Println("into setupPubSub ------------------------ !!!!!")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()
	// Create the Topic if it doesn't exist.
	t := client.Topic(topic)
	if exists, err := t.Exists(ctx); err != nil {
		log.Println("pub sub topic exists !!!!")
		return nil, fmt.Errorf("error checking PubSub topic: %w", err)
	} else if !exists {
		log.Println("pub sub topic creating-------------- !!!!")
		if t, err = client.CreateTopic(ctx, topic); err != nil {
			log.Println("pub sub topic creating-------------- !!!!")
			return nil, fmt.Errorf("error creating PubSub topic: %w", err)
		}
	}

	// Create the Subscription if it doesn't exist.
	s := client.Subscription(sub)
	if exists, err := s.Exists(ctx); err != nil {
		log.Println("pub sub subscription exists-------------- !!!!")
		return nil, fmt.Errorf("error checking for subscription: %w", err)
	} else if !exists {
		if s, err = client.CreateSubscription(ctx, sub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			log.Println("pub sub subscription creating-------------- !!!!")
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}
	}

	return s, nil
}
