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
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

var topicID string

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

		// Cleanup resources from previous tests.
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
	})

	return client
}

func TestCreate(t *testing.T) {
	client := setup(t)
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := create(buf, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
	ok, err := client.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want topic = %q", topicID)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		topics, err := list(tc.ProjectID)
		if err != nil {
			r.Errorf("failed to list topics: %v", err)
		}

		for _, t := range topics {
			if t.ID() == topicID {
				return // PASS
			}
		}

		topicIDs := make([]string, len(topics))
		for i, t := range topics {
			topicIDs[i] = t.ID()
		}
		r.Errorf("got %+v; want a list with topic = %q", topicIDs, topicID)
	})
}

func TestPublish(t *testing.T) {
	// Nothing much to do here, unless we are consuming.
	// TODO(jbd): Merge topics and subscriptions programs maybe?
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	client.CreateTopic(ctx, topicID)
	buf := new(bytes.Buffer)
	if err := publish(buf, tc.ProjectID, topicID, "hello world"); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishThatScales(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	client.CreateTopic(ctx, topicID)
	buf := new(bytes.Buffer)
	if err := publishThatScales(buf, tc.ProjectID, topicID, 10); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishWithSettings(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	client.CreateTopic(ctx, topicID)
	if err := publishWithSettings(ioutil.Discard, tc.ProjectID, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishCustomAttributes(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	client.CreateTopic(ctx, topicID)
	buf := new(bytes.Buffer)
	if err := publishCustomAttributes(buf, tc.ProjectID, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestIAM(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	client.CreateTopic(ctx, topicID)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		perms, err := testPermissions(buf, tc.ProjectID, topicID)
		if err != nil {
			r.Errorf("testPermissions: %v", err)
		}
		if len(perms) == 0 {
			r.Errorf("want non-zero perms")
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := addUsers(tc.ProjectID, topicID); err != nil {
			r.Errorf("addUsers: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		policy, err := policy(buf, tc.ProjectID, topicID)
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

func TestPublishWithOrderingKey(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client := setup(t)
	client.CreateTopic(ctx, topicID)
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
	client.CreateTopic(ctx, topicID)
	buf := new(bytes.Buffer)
	resumePublishWithOrderingKey(buf, tc.ProjectID, topicID)

	got := buf.String()
	want := "Published a message with ordering key successfully\n"
	if got != want {
		t.Fatalf("failed to resume with ordering keys:\n got: %v", got)
	}
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
		_, err := client.CreateTopic(ctx, topicID)
		if err != nil {
			t.Fatalf("CreateTopic: %v", err)
		}
	}

	buf := new(bytes.Buffer)
	if err := delete(buf, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to delete topic (%q): %v", topicID, err)
	}
	ok, err = client.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}
	if ok {
		t.Fatalf("got topic = %q; want none", topicID)
	}
}
