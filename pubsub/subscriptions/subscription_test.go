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

// package subscriptions is a tool to manage Google Cloud Pub/Sub subscriptions by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package subscriptions

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var topicID string
var subID string

// once guards cleanup related operations in setup. No need to set up and tear
// down every time, so this speeds things up.
var once sync.Once

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	topicID = "test-sub-topic"
	subID = "test-sub"
	var err error
	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Cleanup resources from the previous tests.
	once.Do(func() {
		topic := client.Topic(topicID)
		ok, err := topic.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if topic exists: %v", err)
		}
		if ok {
			if err := topic.Delete(ctx); err != nil {
				t.Fatalf("failed to cleanup the topic (%q): %v", topicID, err)
			}
		}
		sub := client.Subscription(subID)
		ok, err = sub.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if subscription exists: %v", err)
		}
		if ok {
			if err := sub.Delete(ctx); err != nil {
				t.Fatalf("failed to cleanup the subscription (%q): %v", subID, err)
			}
		}
	})

	return client
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	topic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	buf := new(bytes.Buffer)
	if err := create(buf, tc.ProjectID, subID, topic); err != nil {
		t.Fatalf("failed to create a subscription: %v", err)
	}
	ok, err := client.Subscription(subID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want sub = %q", subID)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		subs, err := list(tc.ProjectID)
		if err != nil {
			r.Errorf("failed to list subscriptions: %v", err)
			return
		}

		for _, sub := range subs {
			if sub.ID() == subID {
				return // PASS
			}
		}

		subIDs := make([]string, len(subs))
		for i, sub := range subs {
			subIDs[i] = sub.ID()
		}
		r.Errorf("got %+v; want a list with subscription %q", subIDs, subID)
	})
}

func TestIAM(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		perms, err := testPermissions(buf, tc.ProjectID, subID)
		if err != nil {
			r.Errorf("testPermissions: %v", err)
		}
		if len(perms) == 0 {
			r.Errorf("want non-zero perms")
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := addUsers(tc.ProjectID, subID); err != nil {
			r.Errorf("addUsers: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		policy, err := policy(buf, tc.ProjectID, subID)
		if err != nil {
			r.Errorf("policy: %v", err)
		}
		if role, member := iam.Editor, "group:cloud-logs@google.com"; !policy.HasRole(member, role) {
			r.Errorf("want %q as viewer, policy=%v", member, policy)
		}
		if role, member := iam.Viewer, iam.AllUsers; !policy.HasRole(member, role) {
			r.Errorf("want %q as viewer, policy=%v", member, policy)
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)

	topic := client.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}
	if !ok {
		topic, err := client.CreateTopic(ctx, topicID)
		if err != nil {
			t.Fatalf("CreateTopic: %v", err)
		}
		_, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			t.Fatalf("CreateSubscription: %v", err)
		}
	}

	buf := new(bytes.Buffer)
	if err := delete(buf, tc.ProjectID, subID); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", subID, err)
	}
	ok, err = client.Subscription(subID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("sub = %q; want none", subID)
	}
}

func TestPullMsgsSync(t *testing.T) {
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDSync := topicID + "-sync"
	subIDSync := subID + "-sync"

	topic, err := getOrCreateTopic(ctx, client, topicIDSync)
	if err != nil {
		t.Fatalf("getOrCreateTopic: %v", err)
	}
	defer topic.Delete(ctx)
	defer topic.Stop()

	sub, err := getOrCreateSub(ctx, client, topic, subIDSync)
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	// Publish 5 messages on the topic.
	const numMsgs = 5
	publishMsgs(ctx, topic, numMsgs)

	buf := new(bytes.Buffer)
	err = pullMsgsSync(buf, tc.ProjectID, subIDSync, topic)
	if err != nil {
		t.Fatalf("failed to pull messages: %v", err)
	}
	// Check for number of newlines, which should correspond with number of messages.
	if got := strings.Count(buf.String(), "\n"); got != numMsgs {
		t.Fatalf("pullMsgsSync got %d messages, want %d", got, numMsgs)
	}
}

func TestPullMsgsConcurrencyControl(t *testing.T) {
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDConc := topicID + "-conc"
	subIDConc := subID + "-conc"

	topic, err := getOrCreateTopic(ctx, client, topicIDConc)
	if err != nil {
		t.Fatalf("getOrCreateTopic: %v", err)
	}
	defer topic.Delete(ctx)
	defer topic.Stop()

	sub, err := getOrCreateSub(ctx, client, topic, subIDConc)
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	// Publish 5 message to test with.
	const numMsgs = 5
	publishMsgs(ctx, topic, 5)

	buf := new(bytes.Buffer)
	if err := pullMsgsConcurrenyControl(buf, tc.ProjectID, subIDConc); err != nil {
		t.Fatalf("failed to pull messages: %v", err)
	}
	// Check for number of newlines, which should correspond with number of messages.
	if got := strings.Count(buf.String(), "\n"); got != numMsgs {
		t.Fatalf("pullMsgsConcurrencyControl got %d messages, want %d", got, numMsgs)
	}
}

func publishMsgs(ctx context.Context, t *pubsub.Topic, numMsgs int) error {
	var results []*pubsub.PublishResult
	for i := 0; i < numMsgs; i++ {
		res := t.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("message#%d", i)),
		})
		results = append(results, res)
	}
	// Check that all messages were published.
	for _, r := range results {
		if _, err := r.Get(ctx); err != nil {
			return fmt.Errorf("Get publish result: %v", err)
		}
	}
	return nil
}

// getOrCreateTopic gets a topic or creates it if it doesn't exist.
func getOrCreateTopic(ctx context.Context, client *pubsub.Client, topicID string) (*pubsub.Topic, error) {
	topic := client.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if topic exists: %v", err)
	}
	if !ok {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic (%q): %v", topicID, err)
		}
	}
	return topic, nil
}

// getOrCreateSub gets a subscription or creates it if it doesn't exist.
func getOrCreateSub(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, subID string) (*pubsub.Subscription, error) {
	sub := client.Subscription(subID)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if subscription exists: %v", err)
	}
	if !ok {
		sub, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription (%q): %v", topicID, err)
		}
	}
	return sub, nil
}
