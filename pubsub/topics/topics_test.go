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

// Package topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.package topics
package topics

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/pstest"
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	trace "cloud.google.com/go/trace/apiv1"
	"cloud.google.com/go/trace/apiv1/tracepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var topicID string
var topicName string

const (
	topicPrefix = "topic"
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
		topicName = fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID)

		// Cleanup resources from previous tests.
		it := client.TopicAdminClient.ListTopics(ctx, &pubsubpb.ListTopicsRequest{
			Project: tc.ProjectID,
		})
		for {
			t, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			tt := strings.Split(t.GetName(), "/")
			tID := tt[len(tt)-1]
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
					req := &pubsubpb.DeleteTopicRequest{
						Topic: t.GetName(),
					}
					if err := client.TopicAdminClient.DeleteTopic(ctx, req); err != nil {
						fmt.Printf("Delete topic err: %v: %v", t.String(), err)
					}
				}
			}
		}
	})

	return client
}

func TestCreate(t *testing.T) {
	setup(t)
	tc := testutil.SystemTest(t)
	if err := create(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := listTopics(io.Discard, tc.ProjectID); err != nil {
			r.Errorf("failed to list topics: %v", err)
		}
	})
}

func TestPublish(t *testing.T) {
	// Nothing much to do here, unless we are consuming.
	// TODO(jbd): Merge topics and subscriptions programs maybe?
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publish(io.Discard, tc.ProjectID, topicID, "hello world"); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishThatScales(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publishThatScales(io.Discard, tc.ProjectID, topicID, 10); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishWithSettings(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publishWithSettings(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishCustomAttributes(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publishCustomAttributes(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishWithRetrySettings(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publishWithRetrySettings(io.Discard, tc.ProjectID, topicID, "hello world"); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestIAM(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		perms, err := testPermissions(io.Discard, tc.ProjectID, topicID)
		if err != nil {
			r.Errorf("testPermissions: %v", err)
		}
		if len(perms) == 0 {
			r.Errorf("want non-zero perms")
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := addUsersToTopic(io.Discard, tc.ProjectID, topicID); err != nil {
			r.Errorf("addUsers: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if err := getIAMPolicy(buf, tc.ProjectID, topicID); err != nil {
			r.Errorf("getIAMPolicy: %v", err)
		}
		got := buf.String()

		if !strings.Contains(got, "role: roles/editor, member: group:cloud-logs@google.com") {
			r.Errorf("want %s as editor", "group:cloud-logs@google.com")
		}
	})
}

func TestPublishWithOrderingKey(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	buf := new(bytes.Buffer)
	publishWithOrderingKey(buf, tc.ProjectID, topicID)

	got := buf.String()
	want := "Published 4 messages with ordering keys successfully\n"
	if got != want {
		t.Fatalf("failed to publish with ordering keys:\n got: %v", got)
	}
}

func TestResumePublishWithOrderingKey(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	buf := new(bytes.Buffer)
	resumePublishWithOrderingKey(buf, tc.ProjectID, topicID)

	got := buf.String()
	want := "Published a message with ordering key successfully\n"
	if got != want {
		t.Fatalf("failed to resume with ordering keys:\n got: %v", got)
	}
}

func TestPublishWithFlowControl(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publishWithFlowControlSettings(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestDelete(t *testing.T) {
	tc := testutil.SystemTest(t)

	if err := delete(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to delete topic (%q): %v", topicID, err)
	}
}

func TestTopicKinesisIngestion(t *testing.T) {
	tc := testutil.SystemTest(t)

	// Use the pstest fake with emulator settings since Pub/Sub service expects real AWS Kinesis
	// resources, which we cannot provide in a samples test.
	srv := pstest.NewServer()
	t.Setenv("PUBSUB_EMULATOR_HOST", srv.Addr)

	if err := createTopicWithKinesisIngestion(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to create a topic with kinesis ingestion: %v", err)
	}

	// test updateTopicType
	if err := updateTopicType(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to update a topic type to kinesis ingestion: %v", err)
	}
}

func TestTopicCloudStorageIngestion(t *testing.T) {
	tc := testutil.SystemTest(t)

	srv := pstest.NewServer()
	t.Setenv("PUBSUB_EMULATOR_HOST", srv.Addr)

	// Test creating a cloud storage ingestion topic with Text input format.
	if err := createTopicWithCloudStorageIngestion(io.Discard, tc.ProjectID, topicID, "fake-bucket", "**.txt", "2006-01-02T15:04:05Z", ","); err != nil {
		t.Fatalf("failed to create a topic with cloud storage ingestion: %v", err)
	}
}

func TestPublishOpenTelemetryTracing(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	ctx := context.Background()

	// Use the pstest fake with emulator settings.
	srv := pstest.NewServer()
	t.Setenv("PUBSUB_EMULATOR_HOST", srv.Addr)
	setup(t)

	otelTopicID := topicID + "-otel"

	if err := create(buf, tc.ProjectID, otelTopicID); err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}
	defer delete(buf, tc.ProjectID, otelTopicID)

	if err := publishOpenTelemetryTracing(buf, tc.ProjectID, otelTopicID, 1.0); err != nil {
		t.Fatalf("failed to publish message with otel tracing: %v", err)
	}
	got := buf.String()
	want := "Published a traced message"
	if !strings.Contains(got, want) {
		t.Fatalf("failed to publish message:\n got: %v", got)
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
			Filter:    fmt.Sprintf("+messaging.destination.name:%v", otelTopicID),
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
		// Two traces are expected: create and (batch) publish traces.
		if want := 2; numTrace != want {
			r.Errorf("got %d traces, want %d", numTrace, want)
		}
	})
}

func TestPublishWithCompression(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	createTopic(ctx, client, topicName)
	if err := publishWithCompression(io.Discard, tc.ProjectID, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestCreateTopicWithSMT(t *testing.T) {
	setup(t)
	tc := testutil.SystemTest(t)
	smtTopicID := topicID + "-smt"
	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := createTopicWithSMT(buf, tc.ProjectID, smtTopicID)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.AlreadyExists {
				return // This is expected on a retry.
			}
			r.Errorf("failed to create topic with SMT: %v", err)
			return
		}

		got := buf.String()
		want := "Created topic with message transform"
		if !strings.Contains(got, want) {
			r.Errorf("got %q, want to contain %q", got, want)
		}
	})
}

func createTopic(ctx context.Context, client *pubsub.Client, topicName string) error {
	_, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: topicName,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}
