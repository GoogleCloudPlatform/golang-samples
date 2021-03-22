// Copyright 2021 Google LLC
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

// Package psltest contains utilities for pubsublite tests.
package psltest

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/pubsublite"
	"google.golang.org/api/iterator"
)

// Cleanup deletes all previous test topics/subscriptions from previous test
// runs. This prevents previous test failures from building up resources that
// count against quota.
func Cleanup(t *testing.T, client *pubsublite.AdminClient, proj string, zones []string) {
	ctx := context.Background()

	for _, zone := range zones {
		parent := fmt.Sprintf("projects/%s/locations/%s", proj, zone)
		topicIter := client.Topics(ctx, parent)
		for {
			topic, err := topicIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("topicIter.Next got err: %v", err)
			}
			if err := client.DeleteTopic(ctx, topic.Name); err != nil {
				t.Fatalf("AdminClient.DeleteTopic got err: %v", err)
			}
		}

		subIter := client.Subscriptions(ctx, parent)
		for {
			sub, err := subIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("subIter.Next() got err: %v", err)
			}
			if err := client.DeleteSubscription(ctx, sub.Name); err != nil {
				t.Fatalf("AdminClient.DeleteSubscription got err: %v", err)
			}
		}
	}
}

// MustCreateTopic creates a Pub/Sub Lite topic and fails the test if
// unsuccessful.
func MustCreateTopic(ctx context.Context, t *testing.T, client *pubsublite.AdminClient, topicPath string) *pubsublite.TopicConfig {
	cfg := defaultTopicConfig(topicPath)
	topicConfig, err := client.CreateTopic(ctx, *cfg)
	if err != nil {
		t.Fatalf("AdminClient.CreateTopic got err: %v", err)
	}
	return topicConfig
}

func defaultTopicConfig(topicPath string) *pubsublite.TopicConfig {
	cfg := &pubsublite.TopicConfig{
		Name:                       topicPath,
		PartitionCount:             2,
		PublishCapacityMiBPerSec:   4,
		SubscribeCapacityMiBPerSec: 4,
		PerPartitionBytes:          30 * 1024 * 1024 * 1024, // 30 GiB
		RetentionDuration:          pubsublite.InfiniteRetention,
	}
	return cfg
}

// MustCreateSubscription creates a Pub/Sub Lite subscription and fails the test
// if unsuccessful.
func MustCreateSubscription(ctx context.Context, t *testing.T, client *pubsublite.AdminClient, topicPath, subPath string) *pubsublite.SubscriptionConfig {
	cfg := defaultSubConfig(topicPath, subPath)
	subConfig, err := client.CreateSubscription(ctx, *cfg)
	if err != nil {
		t.Fatalf("AdminClient.CreateSubscription got err: %v", err)
	}
	return subConfig
}

func defaultSubConfig(topicPath, subPath string) *pubsublite.SubscriptionConfig {
	cfg := &pubsublite.SubscriptionConfig{
		Name:                subPath,
		Topic:               topicPath,
		DeliveryRequirement: pubsublite.DeliverImmediately,
	}
	return cfg
}
