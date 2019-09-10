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

	"cloud.google.com/go/pubsub"
)

// setupPubSub creates a subscription to the given topic.
func setupPubSub(projectID, topic, sub string) (*pubsub.Subscription, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	// Create the Topic if it doesn't exist.
	t := client.Topic(topic)
	if exists, err := t.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking PubSub topic: %v", err)
	} else if !exists {
		if t, err = client.CreateTopic(ctx, topic); err != nil {
			return nil, fmt.Errorf("error creating PubSub topic: %v", err)
		}
	}

	// Create the Subscription if it doesn't exist.
	s := client.Subscription(sub)
	if exists, err := s.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking for subscription: %v", err)
	} else if !exists {
		if s, err = client.CreateSubscription(ctx, sub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return nil, fmt.Errorf("failed to create subscription: %v", err)
		}
	}

	return s, nil
}
