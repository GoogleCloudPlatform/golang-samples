// Copyright 2025 Google LLC
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

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"cloud.google.com/go/storage"
	trace "cloud.google.com/go/trace/apiv1"
	"cloud.google.com/go/trace/apiv1/tracepb"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
					// Topic deletion can be fire and forget
					t.Delete(ctx)
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
					// Subscription deletion can be fire and forget
					s.Delete(ctx)
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

	var topic *pubsub.Topic
	var err error
	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			t.Fatalf("CreateTopic: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if err := create(buf, tc.ProjectID, subID, topic); err != nil {
			t.Fatalf("failed to create a subscription: %v", err)
		}
		got := buf.String()
		want := "Created subscription"
		if !strings.Contains(got, want) {
			t.Fatalf("got: %s, want: %v", got, want)
		}
		ok, err := client.Subscription(subID).Exists(context.Background())
		if err != nil {
			t.Fatalf("failed to check if sub exists: %v", err)
		}
		if !ok {
			t.Fatalf("got none; want sub = %q", subID)
		}
	})
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
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	asyncTopicID := topicID + "-async"
	asyncSubID := subID + "-async"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		topic, err := getOrCreateTopic(ctx, client, asyncTopicID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
		}
		defer topic.Delete(ctx)
		defer topic.Stop()

		cfg := &pubsub.SubscriptionConfig{
			Topic: topic,
		}
		sub, err := getOrCreateSub(ctx, client, asyncSubID, cfg)
		if err != nil {
			r.Errorf("getOrCreateSub: %v", err)
		}
		defer sub.Delete(ctx)

		// Publish 1 message. This avoids race conditions
		// when calling fmt.Fprintf from multiple receive
		// callbacks. This is sufficient for testing since
		// we're not testing client library functionality,
		// and makes the sample more readable.
		const numMsgs = 1
		publishMsgs(ctx, topic, numMsgs)

		buf := new(bytes.Buffer)
		err = pullMsgs(buf, tc.ProjectID, asyncSubID)
		if err != nil {
			r.Errorf("failed to pull messages: %v", err)
		}
		got := buf.String()
		want := fmt.Sprintf("Received %d messages\n", numMsgs)
		if !strings.Contains(got, want) {
			r.Errorf("pullMsgs got %s\nwant %s", got, want)
		}
	})
}

func TestPullMsgsSync(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDSync := topicID + "-sync"
	subIDSync := subID + "-sync"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		topic, err := getOrCreateTopic(ctx, client, topicIDSync)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
		}
		defer topic.Delete(ctx)
		defer topic.Stop()

		cfg := &pubsub.SubscriptionConfig{
			Topic: topic,
		}
		sub, err := getOrCreateSub(ctx, client, subIDSync, cfg)
		if err != nil {
			r.Errorf("getOrCreateSub: %v", err)
		}
		defer sub.Delete(ctx)

		// Publish 1 message. This avoids race conditions
		// when calling fmt.Fprintf from multiple receive
		// callbacks. This is sufficient for testing since
		// we're not testing client library functionality,
		// and makes the sample more readable.
		const numMsgs = 1
		publishMsgs(ctx, topic, numMsgs)

		buf := new(bytes.Buffer)
		err = pullMsgsSync(buf, tc.ProjectID, subIDSync)
		if err != nil {
			r.Errorf("failed to pull messages: %v", err)
		}

		got := buf.String()
		want := fmt.Sprintf("Received %d messages\n", numMsgs)
		if !strings.Contains(got, want) {
			r.Errorf("pullMsgsSync got %s\nwant %s", got, want)
		}
	})
}

func TestPullMsgsConcurrencyControl(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDConc := topicID + "-conc"
	subIDConc := subID + "-conc"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
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
		if err := pullMsgsConcurrencyControl(buf, tc.ProjectID, subIDConc); err != nil {
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
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDAttributes := topicID + "-attributes"
	subIDAttributes := subID + "-attributes"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		topic, err := getOrCreateTopic(ctx, client, topicIDAttributes)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
		}
		defer topic.Delete(ctx)
		defer topic.Stop()

		cfg := &pubsub.SubscriptionConfig{
			Topic: topic,
		}
		sub, err := getOrCreateSub(ctx, client, subIDAttributes, cfg)
		if err != nil {
			r.Errorf("getOrCreateSub: %v", err)
		}
		defer sub.Delete(ctx)

		res := topic.Publish(ctx, &pubsub.Message{
			Data:       []byte("message with custom attributes"),
			Attributes: map[string]string{"foo": "bar"},
		})
		if _, err := res.Get(ctx); err != nil {
			r.Errorf("Get publish result: %v", err)
		}

		buf := new(bytes.Buffer)
		if err := pullMsgsCustomAttributes(buf, tc.ProjectID, subIDAttributes); err != nil {
			r.Errorf("failed to pull messages: %v", err)
		}

		want := "foo = bar"
		if !strings.Contains(buf.String(), want) {
			r.Errorf("pullMsgsCustomAttributes, got: %s, want %s", buf.String(), want)
		}
	})
}

func TestCreateWithDeadLetterPolicy(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
		t.Fatalf("detached subscription should have detached=true")
	}
}

func TestCreateWithFilter(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	filterSubID := subID + "-filter"

	topic, err := getOrCreateTopic(ctx, client, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	buf := new(bytes.Buffer)
	filter := "attributes.author=\"unknown\""
	if err := createWithFilter(buf, tc.ProjectID, filterSubID, filter, topic); err != nil {
		t.Fatalf("failed to create subscription with filter: %v", err)
	}

	filterSub := client.Subscription(filterSubID)
	defer filterSub.Delete(ctx)
	ok, err := filterSub.Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want sub = %q", filterSubID)
	}
	cfg, err := filterSub.Config(ctx)
	if err != nil {
		t.Fatalf("failed to get config for sub with filter: %v", err)
	}
	if cfg.Filter != filter {
		t.Fatalf("subscription filter got: %s\nwant: %s", cfg.Filter, filter)
	}
}

func TestCreatePushSubscription(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()

	t.Run("default push subscription", func(t *testing.T) {
		topicID := topicID + "-default-push"
		subID := subID + "-default-push"
		t.Cleanup(func() {
			// Don't check delete errors since if it doesn't exist
			// that's fine.
			topic := client.Topic(topicID)
			topic.Delete(ctx)

			sub := client.Subscription(subID)
			sub.Delete(ctx)
		})

		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			topic, err := getOrCreateTopic(ctx, client, topicID)
			if err != nil {
				r.Errorf("CreateTopic: %v", err)
			}

			var b bytes.Buffer
			endpoint := "https://my-test-project.appspot.com/push"
			if err := createWithEndpoint(&b, tc.ProjectID, subID, topic, endpoint); err != nil {
				r.Errorf("failed to create push subscription: %v", err)
			}

			got := b.String()
			want := "Created push subscription"
			if !strings.Contains(got, want) {
				r.Errorf("got %s, want %s", got, want)
			}
		})
	})

	t.Run("no wrapper", func(t *testing.T) {
		topicID := topicID + "-no-wrapper"
		subID := subID + "-no-wrapper"

		t.Cleanup(func() {
			// Don't check delete errors since if it doesn't exist
			// that's fine.
			topic := client.Topic(topicID)
			topic.Delete(ctx)

			sub := client.Subscription(subID)
			sub.Delete(ctx)
		})

		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			topic, err := getOrCreateTopic(ctx, client, topicID)
			if err != nil {
				r.Errorf("CreateTopic: %v", err)
			}

			var b bytes.Buffer
			endpoint := "https://my-test-project.appspot.com/push"
			if err := createPushNoWrapperSubscription(&b, tc.ProjectID, subID, topic, endpoint); err != nil {
				r.Errorf("failed to create push subscription: %v", err)
			}

			got := b.String()
			want := "Created push no wrapper subscription"
			if !strings.Contains(got, want) {
				r.Errorf("got %s, want %s", got, want)
			}
		})
	})
}

func TestCreateBigQuerySubscription(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	bqSubID := subID + "-bigquery"

	topic, err := getOrCreateTopic(ctx, client, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	buf := new(bytes.Buffer)

	datasetID := fmt.Sprintf("go_samples_dataset_%d", time.Now().UnixNano())
	tableID := fmt.Sprintf("go_samples_table_%d", time.Now().UnixNano())
	if err := createBigQueryTable(tc.ProjectID, datasetID, tableID); err != nil {
		t.Fatalf("failed to create bigquery table: %v", err)
	}

	bqTable := fmt.Sprintf("%s.%s.%s", tc.ProjectID, datasetID, tableID)

	if err := createBigQuerySubscription(buf, tc.ProjectID, bqSubID, topic, bqTable); err != nil {
		t.Fatalf("failed to create bigquery subscription: %v", err)
	}

	sub := client.Subscription(bqSubID)
	sub.Delete(ctx)
	if err := deleteBigQueryDataset(tc.ProjectID, datasetID); err != nil {
		t.Logf("failed to delete bigquery dataset: %v", err)
	}
}

func TestCreateCloudStorageSubscription(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	storageSubID := subID + "-cloud-storage"

	topic, err := getOrCreateTopic(ctx, client, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	var buf bytes.Buffer

	// Use the same bucket across test instances. This
	// is safe since we're not writing to the bucket
	// and this makes us not have to do bucket cleanups.
	bucketID := fmt.Sprintf("%s-%s", tc.ProjectID, "pubsub-storage-sub-sink")
	if err := createOrGetStorageBucket(tc.ProjectID, bucketID); err != nil {
		t.Fatalf("failed to get or create storage bucket: %v", err)
	}

	if err := createCloudStorageSubscription(&buf, tc.ProjectID, storageSubID, topic, bucketID); err != nil {
		t.Fatalf("failed to create cloud storage subscription: %v", err)
	}

	sub := client.Subscription(storageSubID)
	sub.Delete(ctx)
}

func TestCreateSubscriptionWithExactlyOnceDelivery(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	eodSub := subID + "-create-eod"

	topic, err := getOrCreateTopic(ctx, client, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	buf := new(bytes.Buffer)

	if err := createSubscriptionWithExactlyOnceDelivery(buf, tc.ProjectID, eodSub, topic); err != nil {
		t.Fatalf("failed to create exactly once delivery subscription: %v", err)
	}

	sub := client.Subscription(eodSub)
	sub.Delete(ctx)
}

func TestReceiveMessagesWithExactlyOnceDelivery(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	eodTopicID := topicID + "-eod"
	eodSubID := subID + "-eod"

	topic, err := getOrCreateTopic(ctx, client, eodTopicID)
	if err != nil {
		t.Fatalf("getOrCreateTopic: %v", err)
	}
	defer topic.Delete(ctx)
	defer topic.Stop()

	cfg := &pubsub.SubscriptionConfig{
		Topic:                     topic,
		EnableExactlyOnceDelivery: true,
	}
	sub, err := getOrCreateSub(ctx, client, eodSubID, cfg)
	if err != nil {
		t.Fatalf("getOrCreateSub: %v", err)
	}
	defer sub.Delete(ctx)

	// Publish 1 message. This avoids race conditions
	// when calling fmt.Fprintf from multiple receive
	// callbacks. This is sufficient for testing since
	// we're not testing client library functionality,
	// and makes the sample more readable.
	const numMsgs = 1
	publishMsgs(ctx, topic, numMsgs)

	buf := new(bytes.Buffer)
	err = receiveMessagesWithExactlyOnceDeliveryEnabled(buf, tc.ProjectID, eodSubID)
	if err != nil {
		t.Fatalf("failed to pull messages: %v", err)
	}
	got := buf.String()
	want := "Message successfully acked"
	if !strings.Contains(got, want) {
		t.Fatalf("receiveMessagesWithExactlyOnceDeliveryEnabled got %s\nwant %s", got, want)
	}
}

func TestOptimisticSubscribe(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	optTopicID := topicID + "-opt"
	optSubID := subID + "-opt"

	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		topic, err := getOrCreateTopic(ctx, client, optTopicID)
		if err != nil {
			r.Errorf("getOrCreateTopic: %v", err)
		}
		defer topic.Delete(ctx)
		defer topic.Stop()

		buf := new(bytes.Buffer)
		err = optimisticSubscribe(buf, tc.ProjectID, optTopicID, optSubID)
		if err != nil {
			r.Errorf("failed to pull messages: %v", err)
		}

		// Check that we created the subscription instead of using
		// an existing one. We can't test receiving a message
		// since a message published won't be delivered to a new
		// subscription.
		got := buf.String()
		want := "Created subscription"
		if !strings.Contains(got, want) {
			r.Errorf("optimisticSubscribe\ngot: %s\nwant: %s", got, want)
		}

		sub := client.Subscription(optSubID)
		sub.Delete(ctx)
	})
}

func TestSubscribeOpenTelemetryTracing(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	ctx := context.Background()

	// Use the pstest fake with emulator settings.
	srv := pstest.NewServer()
	t.Setenv("PUBSUB_EMULATOR_HOST", srv.Addr)
	client := setup(t)

	otelTopicID := topicID + "-otel"
	otelSubID := subID + "-otel"

	topic, err := client.CreateTopic(ctx, otelTopicID)
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}
	defer topic.Delete(ctx)

	if err := create(buf, tc.ProjectID, otelSubID, topic); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
	defer client.Subscription(otelSubID).Delete(ctx)

	if err := publishMsgs(ctx, topic, 1); err != nil {
		t.Fatalf("failed to publish setup message: %v", err)
	}

	if err := subscribeOpenTelemetryTracing(buf, tc.ProjectID, otelSubID, 1.0); err != nil {
		t.Fatalf("failed to subscribe message with otel tracing: %v", err)
	}
	got := buf.String()
	want := "Received 1 message"
	if !strings.Contains(got, want) {
		t.Fatalf("expected 1 message, got: %s", got)
	}

	traceClient, err := trace.NewClient(ctx)
	if err != nil {
		t.Fatalf("trace client instantiation: %v", err)
	}

	testutil.Retry(t, 3, time.Second, func(r *testutil.R) {
		// Wait some time for the spans to show up in Cloud Trace.
		time.Sleep(5 * time.Second)
		iter := traceClient.ListTraces(ctx, &tracepb.ListTracesRequest{
			ProjectId: tc.ProjectID,
			Filter:    fmt.Sprintf("+messaging.destination.name:%v", otelSubID),
		})
		numTrace := 0
		for {
			_, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				r.Errorf("got err in iter.Next: %v", err)
			}
			numTrace++
		}
		// Three traces are created from subscribe side: subscribe, ack, modack spans.
		if want := 3; numTrace != want {
			r.Errorf("got %d traces, want %d", numTrace, want)
		}
	})
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
			return fmt.Errorf("Get publish result: %w", err)
		}
	}
	return nil
}

// getOrCreateTopic gets a topic or creates it if it doesn't exist.
func getOrCreateTopic(ctx context.Context, client *pubsub.Client, topicID string) (*pubsub.Topic, error) {
	topic := client.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if topic exists: %w", err)
	}
	if !ok {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic (%q): %w", topicID, err)
		}
	}
	return topic, nil
}

// getOrCreateSub gets a subscription or creates it if it doesn't exist.
func getOrCreateSub(ctx context.Context, client *pubsub.Client, subID string, cfg *pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	sub := client.Subscription(subID)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if subscription exists: %w", err)
	}
	if !ok {
		sub, err = client.CreateSubscription(ctx, subID, *cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription (%q): %w", topicID, err)
		}
	}
	return sub, nil
}

func createBigQueryTable(projectID, datasetID, tableID string) error {
	ctx := context.Background()

	c, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error instantiating bigquery client: %w", err)
	}
	dataset := c.Dataset(datasetID)
	if err = dataset.Create(ctx, &bigquery.DatasetMetadata{Location: "US"}); err != nil {
		return fmt.Errorf("error creating dataset: %w", err)
	}

	table := dataset.Table(tableID)
	schema := []*bigquery.FieldSchema{
		{Name: "data", Type: bigquery.BytesFieldType, Required: true},
		{Name: "message_id", Type: bigquery.StringFieldType, Required: true},
		{Name: "attributes", Type: bigquery.StringFieldType, Required: true},
		{Name: "subscription_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "publish_time", Type: bigquery.TimestampFieldType, Required: true},
	}
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func deleteBigQueryDataset(projectID, datasetID string) error {
	ctx := context.Background()

	c, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error instantiating bigquery client: %w", err)
	}
	dataset := c.Dataset(datasetID)
	if err = dataset.DeleteWithContents(ctx); err != nil {
		return fmt.Errorf("error deleting dataset: %w", err)
	}
	return nil
}

func createOrGetStorageBucket(projectID, bucketID string) error {
	ctx := context.Background()

	c, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error instantiating storage client: %w", err)
	}
	b := c.Bucket(bucketID)
	_, err = b.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		if err := b.Create(ctx, projectID, nil); err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error retrieving existing bucket: %w", err)
	}

	return nil
}

func TestCreateSubscriptionWithSMT(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)

	smtSubID := subID + "-smt"
	var topic *pubsub.Topic
	var err error
	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				r.Errorf("CreateTopic failed with unknown err: %v", err)
			}
			// If the topic alredy exists, that's fine.
			if st.Code() != codes.AlreadyExists {
				r.Errorf("CreateTopic: %v", err)
			}
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if err := createSubscriptionWithSMT(buf, tc.ProjectID, smtSubID, topic); err != nil {
			r.Errorf("failed to create subscription with SMT: %v", err)
		}
		got := buf.String()
		want := "Created subscription with message transform"
		if !strings.Contains(got, want) {
			r.Errorf("got: %s, want: %v", got, want)
		}
	})
}
