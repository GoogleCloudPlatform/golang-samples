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
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/go-cmp/cmp"
)

var topicID string
var subID string

const (
	topicPrefix = "topic"
	subPrefix   = "sub"
	expireAge   = 24 * time.Hour
)

// once guards cleanup related operations in setup. No need to set up and tear
// down every time, so this speeds things up.
var once sync.Once

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	var err error
	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	once.Do(func() {
		topicID = fmt.Sprintf("%s-%d", topicPrefix, time.Now().UnixNano())
		subID = fmt.Sprintf("%s-%d", subPrefix, time.Now().UnixNano())

		// Cleanup resources from the previous tests.
		it := client.Topics(ctx)
		for {
			t, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			tID := t.ID()
			p := strings.Split(tID, "-")

			// Only delete resources created from these tests.
			if p[0] == topicPrefix {
				tCreated := p[1]
				timestamp, err := strconv.ParseInt(tCreated, 10, 64)
				if err != nil {
					continue
				}
				timeTCreated := time.Unix(0, timestamp)
				if time.Since(timeTCreated) > expireAge {
					if err := t.Delete(ctx); err != nil {
						fmt.Printf("Delete topic err: %v: %v", t.String(), err)
					}
				}
			}
		}
		subIter := client.Subscriptions(ctx)
		for {
			s, err := subIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			sID := s.ID()
			p := strings.Split(sID, "-")

			// Only delete resources created from these tests.
			if p[0] == subPrefix {
				tCreated := p[1]
				timestamp, err := strconv.ParseInt(tCreated, 10, 64)
				if err != nil {
					continue
				}
				timeTCreated := time.Unix(0, timestamp)
				if time.Since(timeTCreated) > expireAge {
					if err := s.Delete(ctx); err != nil {
						fmt.Printf("Delete sub err: %v: %v", s.String(), err)
					}
				}
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

func TestPullMsgsAsync(t *testing.T) {
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	asyncTopicID := topicID + "-async"
	asyncSubID := subID + "-async"

	topic, err := getOrCreateTopic(ctx, client, asyncTopicID)
	if err != nil {
		t.Fatalf("getOrCreateTopic: %v", err)
	}
	defer topic.Delete(ctx)
	defer topic.Stop()

	cfg := &pubsub.SubscriptionConfig{
		Topic: topic,
	}
	sub, err := getOrCreateSub(ctx, client, asyncSubID, cfg)
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	// Publish 10 messages on the topic.
	const numMsgs = 10
	publishMsgs(ctx, topic, numMsgs)

	buf := new(bytes.Buffer)
	err = pullMsgs(buf, tc.ProjectID, asyncSubID)
	if err != nil {
		t.Fatalf("failed to pull messages: %v", err)
	}
	// Check for number of newlines, which should correspond with number of messages.
	if got := strings.Count(buf.String(), "\n"); got != numMsgs {
		t.Fatalf("pullMsgsSync got %d messages, want %d", got, numMsgs)
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

	cfg := &pubsub.SubscriptionConfig{
		Topic: topic,
	}
	sub, err := getOrCreateSub(ctx, client, subIDSync, cfg)
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	// Publish 5 messages on the topic.
	const numMsgs = 5
	publishMsgs(ctx, topic, numMsgs)

	buf := new(bytes.Buffer)
	err = pullMsgsSync(buf, tc.ProjectID, subIDSync)
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

	testutil.Retry(t, 3, time.Second, func(r *testutil.R) {
		topic, err := getOrCreateTopic(ctx, client, topicIDConc)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
		}
		defer topic.Delete(ctx)
		defer topic.Stop()

		cfg := &pubsub.SubscriptionConfig{
			Topic: topic,
		}
		sub, err := getOrCreateSub(ctx, client, subIDConc, cfg)
		if err != nil {
			r.Errorf("getOrCreateSub: %v", err)
		}
		defer sub.Delete(ctx)

		// Publish 5 message to test with.
		const numMsgs = 5
		publishMsgs(ctx, topic, numMsgs)

		buf := new(bytes.Buffer)
		if err := pullMsgsConcurrenyControl(buf, tc.ProjectID, subIDConc); err != nil {
			r.Errorf("failed to pull messages: %v", err)
		}
		got := buf.String()
		want := fmt.Sprintf("Received %d messages\n", numMsgs)
		if got != want {
			r.Errorf("pullMsgsConcurrencyControl got %s\nwant %s", got, want)
		}
	})
}

func TestPullMsgsCustomAttributes(t *testing.T) {
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDAttributes := topicID + "-attributes"
	subIDAttributes := subID + "-attributes"

	topic, err := getOrCreateTopic(ctx, client, topicIDAttributes)
	if err != nil {
		t.Fatalf("getOrCreateTopic: %v", err)
	}
	defer topic.Delete(ctx)
	defer topic.Stop()

	cfg := &pubsub.SubscriptionConfig{
		Topic: topic,
	}
	sub, err := getOrCreateSub(ctx, client, subIDAttributes, cfg)
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	res := topic.Publish(ctx, &pubsub.Message{
		Data:       []byte("message with custom attributes"),
		Attributes: map[string]string{"foo": "bar"},
	})
	if _, err := res.Get(ctx); err != nil {
		t.Fatalf("Get publish result: %v", err)
	}

	buf := new(bytes.Buffer)
	if err := pullMsgsCustomAttributes(buf, tc.ProjectID, subIDAttributes); err != nil {
		t.Fatalf("failed to pull messages: %v", err)
	}

	want := "foo = bar"
	if !strings.Contains(buf.String(), want) {
		t.Fatalf("pullMsgsCustomAttributes, got: %s, want %s", buf.String(), want)
	}
}

func TestCreateWithDeadLetterPolicy(t *testing.T) {
	client := setup(t)
	defer client.Close()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	deadLetterSourceID := topicID + "-dead-letter-source"
	deadLetterSubID := subID + "-dead-letter-sub"
	deadLetterSinkID := topicID + "-dead-letter-sink"

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		deadLetterSourceTopic, err := getOrCreateTopic(ctx, client, deadLetterSourceID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
			return
		}
		defer deadLetterSourceTopic.Delete(ctx)
		defer deadLetterSourceTopic.Stop()

		deadLetterSinkTopic, err := getOrCreateTopic(ctx, client, deadLetterSinkID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
			return
		}
		defer deadLetterSinkTopic.Delete(ctx)
		defer deadLetterSinkTopic.Stop()

		buf := new(bytes.Buffer)
		if err := createSubWithDeadLetter(buf, tc.ProjectID, deadLetterSubID, deadLetterSourceID, deadLetterSinkTopic.String()); err != nil {
			r.Errorf("createSubWithDeadLetter failed: %v", err)
			return
		}
		sub := client.Subscription(deadLetterSubID)
		ok, err := sub.Exists(context.Background())
		if err != nil {
			r.Errorf("sub.Exists failed: %v", err)
			return
		}
		if !ok {
			r.Errorf("got none; want sub = %q", deadLetterSubID)
			return
		}
		defer sub.Delete(ctx)

		cfg, err := sub.Config(ctx)
		if err != nil {
			r.Errorf("createSubWithDeadLetter config: %v", err)
			return
		}
		got := cfg.DeadLetterPolicy
		want := &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     deadLetterSinkTopic.String(),
			MaxDeliveryAttempts: 10,
		}
		if !cmp.Equal(got, want) {
			r.Errorf("got cfg: %+v; want cfg: %+v", got, want)
			return
		}
	})
}

func TestUpdateDeadLetterPolicy(t *testing.T) {
	client := setup(t)
	defer client.Close()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	deadLetterSourceID := topicID + "-update-source"
	deadLetterSubID := subID + "-update-sub"
	deadLetterSinkID := topicID + "-update-sink"

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		deadLetterSourceTopic, err := getOrCreateTopic(ctx, client, deadLetterSourceID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
			return
		}
		defer deadLetterSourceTopic.Delete(ctx)
		defer deadLetterSourceTopic.Stop()

		deadLetterSinkTopic, err := getOrCreateTopic(ctx, client, deadLetterSinkID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
			return
		}
		defer deadLetterSinkTopic.Delete(ctx)
		defer deadLetterSinkTopic.Stop()

		buf := new(bytes.Buffer)
		if err := createSubWithDeadLetter(buf, tc.ProjectID, deadLetterSubID, deadLetterSourceID, deadLetterSinkTopic.String()); err != nil {
			r.Errorf("createSubWithDeadLetter failed: %v", err)
			return
		}
		sub := client.Subscription(deadLetterSubID)
		ok, err := sub.Exists(context.Background())
		if err != nil {
			r.Errorf("sub.Exists failed: %v", err)
			return
		}
		if !ok {
			r.Errorf("got none; want sub = %q", deadLetterSubID)
			return
		}
		defer sub.Delete(ctx)

		if err := updateDeadLetter(buf, tc.ProjectID, deadLetterSubID, deadLetterSinkTopic.String()); err != nil {
			r.Errorf("updateDeadLetter failed: %v", err)
			return
		}

		cfg, err := sub.Config(ctx)
		if err != nil {
			r.Errorf("update dead letter policy config: %v", err)
			return
		}
		got := cfg.DeadLetterPolicy
		want := &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     deadLetterSinkTopic.String(),
			MaxDeliveryAttempts: 20,
		}
		if !cmp.Equal(got, want) {
			r.Errorf("got cfg: %+v; want cfg: %+v", got, want)
			return
		}

		if err := removeDeadLetterTopic(buf, tc.ProjectID, deadLetterSubID); err != nil {
			r.Errorf("removeDeadLetterTopic failed: %v", err)
			return
		}
		cfg, err = sub.Config(ctx)
		if err != nil {
			r.Errorf("update dead letter policy config: %v", err)
			return
		}
		got = cfg.DeadLetterPolicy
		if got != nil {
			r.Errorf("got dead letter policy: %+v, want nil", got)
			return
		}
	})
}

func TestPullMsgsDeadLetterDeliveryAttempts(t *testing.T) {
	client := setup(t)
	defer client.Close()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	deadLetterSourceID := topicID + "-delivery-source"
	deadLetterSinkID := topicID + "-delivery-sink"
	deadLetterSubID := subID + "-delivery-sub"

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		deadLetterSourceTopic, err := getOrCreateTopic(ctx, client, deadLetterSourceID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
			return
		}
		defer deadLetterSourceTopic.Delete(ctx)
		defer deadLetterSourceTopic.Stop()

		deadLetterSinkTopic, err := getOrCreateTopic(ctx, client, deadLetterSinkID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
			return
		}
		defer deadLetterSinkTopic.Delete(ctx)
		defer deadLetterSinkTopic.Stop()

		sub, err := getOrCreateSub(ctx, client, deadLetterSubID, &pubsub.SubscriptionConfig{
			Topic: deadLetterSourceTopic,
			DeadLetterPolicy: &pubsub.DeadLetterPolicy{
				DeadLetterTopic:     deadLetterSinkTopic.String(),
				MaxDeliveryAttempts: 10,
			},
		})
		if err != nil {
			r.Errorf("getOrCreateSub: %v", err)
			return
		}
		defer sub.Delete(ctx)

		if err = publishMsgs(ctx, deadLetterSourceTopic, 1); err != nil {
			r.Errorf("publishMsgs failed: %v", err)
			return
		}

		buf := new(bytes.Buffer)
		if err := pullMsgsDeadLetterDeliveryAttempt(buf, tc.ProjectID, deadLetterSubID); err != nil {
			r.Errorf("pullMsgsDeadLetterDeliveryAttempt failed: %v", err)
			return
		}
		got := buf.String()
		want := "delivery attempts: 1"
		if !strings.Contains(got, want) {
			r.Errorf("pullMsgsDeadLetterDeliveryAttempts got %s, want %s", got, want)
			return
		}
	})
}

func TestCreateWithOrdering(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	orderingSubID := subID + "-ordering"

	topic, err := getOrCreateTopic(ctx, client, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	buf := new(bytes.Buffer)
	if err := createWithOrdering(buf, tc.ProjectID, orderingSubID, topic); err != nil {
		t.Fatalf("failed to create a subscription: %v", err)
	}

	orderingSub := client.Subscription(orderingSubID)
	defer orderingSub.Delete(ctx)
	ok, err := orderingSub.Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want sub = %q", orderingSubID)
	}
	cfg, err := orderingSub.Config(ctx)
	if err != nil {
		t.Fatalf("failed to get config for ordering sub: %v", err)
	}
	if !cfg.EnableMessageOrdering {
		t.Fatalf("expected EnableMessageOrdering to be true for sub %s", orderingSubID)
	}
}

func TestDetachSubscription(t *testing.T) {
	client := setup(t)
	defer client.Close()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	detachTopicID := topicID + "-detach"
	detachSubID := "testdetachsubsxyz-" + subID

	topic, err := getOrCreateTopic(ctx, client, detachTopicID)
	if err != nil {
		t.Fatalf("getOrCreateTopic: %v", err)
	}
	defer topic.Delete(ctx)
	defer topic.Stop()

	sub, err := getOrCreateSub(ctx, client, detachSubID, &pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	buf := new(bytes.Buffer)
	if err = detachSubscription(buf, tc.ProjectID, sub.String()); err != nil {
		t.Fatalf("detachSubscription: %v", err)
	}
	got := buf.String()
	want := fmt.Sprintf("Detached subscription %s", sub.String())
	if got != want {
		t.Fatalf("detachSubscription got %s, want %s", got, want)
	}

	cfg, err := sub.Config(ctx)
	if err != nil {
		t.Fatalf("get sub config err: %v", err)
	}
	if !cfg.Detached {
		t.Fatalf("detached subscripion should have detached=true")
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
func getOrCreateSub(ctx context.Context, client *pubsub.Client, subID string, cfg *pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	sub := client.Subscription(subID)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if subscription exists: %v", err)
	}
	if !ok {
		sub, err = client.CreateSubscription(ctx, subID, *cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription (%q): %v", topicID, err)
		}
	}
	return sub, nil
}
