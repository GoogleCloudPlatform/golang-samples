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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub/v2"
	pb "cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/v2/pstest"
	"cloud.google.com/go/storage"
	trace "cloud.google.com/go/trace/apiv1"
	"cloud.google.com/go/trace/apiv1/tracepb"
	"google.golang.org/api/iterator"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var topicID, topicName string
var subID, subName string

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
		topicName = fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID)
		subName = fmt.Sprintf("projects/%s/subscriptions/%s", tc.ProjectID, subID)

		// Cleanup subscription resources from the previous tests.
		req := &pb.ListSubscriptionsRequest{
			Project: tc.ProjectID,
		}
		subIter := client.SubscriptionAdminClient.ListSubscriptions(ctx, req)
		for {
			s, err := subIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			ss := strings.Split(s.GetName(), "/")
			sID := ss[len(ss)-1]
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
					deleteReq := &pb.DeleteSubscriptionRequest{
						Subscription: s.GetName(),
					}
					client.SubscriptionAdminClient.DeleteSubscription(ctx, deleteReq)
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

	var err error
	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		_, err = client.TopicAdminClient.CreateTopic(ctx, &pb.Topic{
			Name: topicName,
		})
		if err != nil {
			t.Fatalf("CreateTopic: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if err := create(buf, tc.ProjectID, topicName, subName); err != nil {
			t.Fatalf("failed to create a subscription: %v", err)
		}
		got := buf.String()
		want := "Created subscription"
		if !strings.Contains(got, want) {
			t.Fatalf("got: %s, want: %v", got, want)
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
			if sub.GetName() == subName {
				return // PASS
			}
		}

		subNames := make([]string, len(subs))
		for i, sub := range subs {
			subNames[i] = sub.GetName()
		}
		r.Errorf("got %+v; want a list with subscription %q", subNames, subName)

	})
}

func TestIAM(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client := setup(t)
	defer client.Close()

	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		createTopic(ctx, client, topicName)
		createSubscription(ctx, client, topicName, subName)

		perms, err := testPermissions(io.Discard, tc.ProjectID, subID)
		if err != nil {
			r.Errorf("testPermissions: %v", err)
		}
		if len(perms) == 0 {
			r.Errorf("want non-zero perms")
		}
	})

	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		if err := addUsersToSubscription(io.Discard, tc.ProjectID, subID); err != nil {
			r.Errorf("addUsers: %v", err)
		}
	})

	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := getIAMPolicy(buf, tc.ProjectID, subName)
		if err != nil {
			r.Errorf("policy: %v", err)
		}

		got := buf.String()
		if !strings.Contains(got, "role: roles/editor, member: group:cloud-logs@google.com") {
			r.Errorf("want %s as editor", "group:cloud-logs@google.com")
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)

	createTopic(ctx, client, topicName)
	createSubscription(ctx, client, topicName, subName)

	if err := delete(io.Discard, tc.ProjectID, subID); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", subID, err)
	}
}

func TestPullMsgsAsync(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	asyncTopic := topicName + "-async"
	asyncSub := subName + "-async"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		createTopic(ctx, client, asyncTopic)
		createSubscription(ctx, client, asyncTopic, asyncSub)

		// Publish 1 message. This avoids race conditions
		// when calling fmt.Fprintf from multiple receive
		// callbacks. This is sufficient for testing since
		// we're not testing client library functionality,
		// and makes the sample more readable.
		const numMsgs = 1
		publisher := client.Publisher(asyncTopic)
		publishMsgs(ctx, publisher, numMsgs)

		buf := new(bytes.Buffer)
		err := pullMsgs(buf, tc.ProjectID, asyncSub)
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

func TestPullMsgsConcurrencyControl(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	topicIDConc := topicName + "-conc"
	subIDConc := subName + "-conc"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		createTopic(ctx, client, topicIDConc)
		createSubscription(ctx, client, topicIDConc, subIDConc)

		// Publish 5 message to test with.
		const numMsgs = 5
		publisher := client.Publisher(topicIDConc)
		publishMsgs(ctx, publisher, numMsgs)

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
	topicIDAttributes := topicName + "-attributes"
	subIDAttributes := subName + "-attributes"

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		createTopic(ctx, client, topicIDAttributes)
		createSubscription(ctx, client, topicIDAttributes, subIDAttributes)

		p := client.Publisher(topicIDAttributes)
		res := p.Publish(ctx, &pubsub.Message{
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

func TestDeadLetterPolicy(t *testing.T) {
	t.Parallel()
	client := setup(t)
	defer client.Close()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	deadLetterSource := topicName + "-dlq-source"
	deadLetterSub := subName + "-dlq-sub"
	deadLetterSink := topicName + "-dlq-sink"

	testutil.Retry(t, 1, time.Second, func(r *testutil.R) {
		createTopic(ctx, client, deadLetterSource)
		createTopic(ctx, client, deadLetterSink)

		if err := createSubWithDeadLetter(io.Discard, tc.ProjectID, deadLetterSource, deadLetterSub, deadLetterSink); err != nil {
			r.Errorf("createSubWithDeadLetter failed: %v", err)
			return
		}
		if err := updateDeadLetter(io.Discard, tc.ProjectID, deadLetterSub, deadLetterSink); err != nil {
			r.Errorf("updateDeadLetter failed: %v", err)
			return
		}
		if err := removeDeadLetterTopic(io.Discard, tc.ProjectID, deadLetterSub); err != nil {
			r.Errorf("removeDeadLetterTopic failed: %v", err)
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
	orderingSub := subName + "-ordering"

	createTopic(ctx, client, topicName)
	if err := createWithOrdering(io.Discard, tc.ProjectID, topicName, orderingSub); err != nil {
		t.Fatalf("failed to create ordered subscription: %v", err)
	}
}

func TestDetachSubscription(t *testing.T) {
	t.Parallel()
	client := setup(t)
	defer client.Close()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	detachTopic := topicName + "-detach"
	detachSub := subName + "detatchtopic"

	createTopic(ctx, client, detachTopic)
	createSubscription(ctx, client, detachTopic, detachSub)

	buf := new(bytes.Buffer)
	if err := detachSubscription(buf, tc.ProjectID, detachSub); err != nil {
		t.Fatalf("detachSubscription: %v", err)
	}
	got := buf.String()
	want := fmt.Sprintf("Detached subscription %s", detachSub)
	if got != want {
		t.Fatalf("detachSubscription got %s, want %s", got, want)
	}
}

func TestCreateWithFilter(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	filterSub := subName + "-filter"

	createTopic(ctx, client, topicName)
	filter := "attributes.author=\"unknown\""
	if err := createWithFilter(io.Discard, tc.ProjectID, topicName, filterSub, filter); err != nil {
		t.Fatalf("failed to create subscription with filter: %v", err)
	}
}

func TestCreatePushSubscription(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()

	t.Run("default push subscription", func(t *testing.T) {
		pushTopic := topicName + "-default-push"
		pushSub := subName + "-default-push"

		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			createTopic(ctx, client, pushTopic)

			buf := new(bytes.Buffer)
			endpoint := "https://my-test-project.appspot.com/push"
			if err := createWithEndpoint(buf, tc.ProjectID, pushTopic, pushSub, endpoint); err != nil {

				r.Errorf("failed to create push subscription: %v", err)
			}

			got := buf.String()
			want := "Created push subscription"
			if !strings.Contains(got, want) {
				r.Errorf("got %s, want %s", got, want)
			}
		})
	})

	t.Run("no wrapper", func(t *testing.T) {
		pushTopic := topicName + "-no-wrapper"
		pushSub := subName + "-no-wrapper"

		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			createTopic(ctx, client, pushTopic)

			buf := new(bytes.Buffer)
			endpoint := "https://my-test-project.appspot.com/push"
			if err := createPushNoWrapperSubscription(buf, tc.ProjectID, pushTopic, pushSub, endpoint); err != nil {
				r.Errorf("failed to create push subscription: %v", err)
			}

			got := buf.String()
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
	bqSub := subName + "-bigquery"

	createTopic(ctx, client, topicName)

	datasetID := fmt.Sprintf("go_samples_dataset_%d", time.Now().UnixNano())
	tableID := fmt.Sprintf("go_samples_table_%d", time.Now().UnixNano())
	if err := createBigQueryTable(tc.ProjectID, datasetID, tableID); err != nil {
		t.Fatalf("failed to create bigquery table: %v", err)
	}

	bqTable := fmt.Sprintf("%s.%s.%s", tc.ProjectID, datasetID, tableID)

	if err := createBigQuerySubscription(io.Discard, tc.ProjectID, topicName, bqSub, bqTable); err != nil {
		t.Fatalf("failed to create bigquery subscription: %v", err)
	}

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
	storageSubName := subName + "-cloud-storage"

	createTopic(ctx, client, topicID)

	// Use the same bucket across test instances. This
	// is safe since we're not writing to the bucket
	// and this makes us not have to do bucket cleanups.
	bucketID := fmt.Sprintf("%s-%s", tc.ProjectID, "pubsub-storage-sub-sink")
	if err := createOrGetStorageBucket(tc.ProjectID, bucketID); err != nil {
		t.Fatalf("failed to get or create storage bucket (%v): %v", bucketID, err)
	}

	if err := createCloudStorageSubscription(io.Discard, tc.ProjectID, topicName, storageSubName, bucketID); err != nil {
		t.Fatalf("failed to create cloud storage subscription: %v", err)
	}
}

func TestCreateSubscriptionWithExactlyOnceDelivery(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	defer client.Close()
	eodSub := subName + "-create-eod"

	createTopic(ctx, client, topicID)
	if err := createSubscriptionWithExactlyOnceDelivery(io.Discard, tc.ProjectID, topicName, eodSub); err != nil {
		t.Fatalf("failed to create exactly once delivery subscription: %v", err)
	}
}

func TestReceiveMessagesWithExactlyOnceDelivery(t *testing.T) {
	t.Parallel()
	client := setup(t)
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	eodTopic := topicName + "-eod"
	eodSub := subName + "-eod"

	createTopic(ctx, client, eodTopic)
	_, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pb.Subscription{
		Name:                      eodSub,
		Topic:                     eodTopic,
		EnableExactlyOnceDelivery: true,
	})

	// Publish 1 message. This avoids race conditions
	// when calling fmt.Fprintf from multiple receive
	// callbacks. This is sufficient for testing since
	// we're not testing client library functionality,
	// and makes the sample more readable.
	const numMsgs = 1
	p := client.Publisher(eodTopic)
	publishMsgs(ctx, p, numMsgs)

	buf := new(bytes.Buffer)
	err = receiveMessagesWithExactlyOnceDeliveryEnabled(buf, tc.ProjectID, eodSub)
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
	optTopic := topicName + "-opt"
	optSub := subName + "-opt"

	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		createTopic(ctx, client, optTopic)

		buf := new(bytes.Buffer)
		err := optimisticSubscribe(buf, tc.ProjectID, optTopic, optSub)
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

	otelTopic := topicName + "-otel"
	otelSubID := subID + "-otel"
	otelSub := subName + "-otel"

	createTopic(ctx, client, otelTopic)
	if err := create(buf, tc.ProjectID, otelTopic, otelSub); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}

	p := client.Publisher(otelTopic)
	if err := publishMsgs(ctx, p, 1); err != nil {
		t.Fatalf("failed to publish setup message: %v", err)
	}

	if err := subscribeOpenTelemetryTracing(buf, tc.ProjectID, otelSub, 1.0); err != nil {
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

func publishMsgs(ctx context.Context, p *pubsub.Publisher, numMsgs int) error {
	var results []*pubsub.PublishResult
	for i := 0; i < numMsgs; i++ {
		res := p.Publish(ctx, &pubsub.Message{
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

// createTopic creates a topic if it doesn't exist.
func createTopic(ctx context.Context, client *pubsub.Client, topic string) error {
	_, err := client.TopicAdminClient.CreateTopic(ctx, &pb.Topic{
		Name: topic,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic (%q): %w", topic, err)
	}
	return nil
}

// createSubscription creates a subscription it if it doesn't exist.
func createSubscription(ctx context.Context, client *pubsub.Client, topic, subscription string) error {
	_, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pb.Subscription{
		Name:  subscription,
		Topic: topic,
	})
	if err != nil {
		return fmt.Errorf("failed to create subscription (%q): %w", subscription, err)
	}
	return nil
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
	if errors.Is(err, storage.ErrBucketNotExist) {
		if err := b.Create(ctx, projectID, nil); err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error retrieving existing bucket: %w", err)
	}

	return nil
}
