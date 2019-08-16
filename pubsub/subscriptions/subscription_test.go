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

// Package subscription is a tool to manage Google Cloud Pub/Sub subscriptions by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package subscription

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var topic *pubsub.Topic
var subID string
var client *pubsub.Client
var once sync.Once // guards cleanup related operations in setup.

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	var err error
	client, err = pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	subID = tc.ProjectID + "-test-sub"
	topicID := tc.ProjectID + "-test-sub-topic"

	// Cleanup resources from the previous failed tests.
	once.Do(func() {
		// Create a topic.
		topic = client.Topic(topicID)
		ok, err := topic.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if topic exists: %v", err)
		}
		if !ok {
			if topic, err = client.CreateTopic(ctx, topicID); err != nil {
				t.Fatalf("failed to create the topic: %v", err)
			}
		}

		// Delete the sub if already exists.
		sub := client.Subscription(subID)
		ok, err = sub.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if sub exists: %v", err)
		}
		if ok {
			if err := client.Subscription(subID).Delete(ctx); err != nil {
				t.Fatalf("failed to cleanup the topic (%q): %v", subID, err)
			}
		}
	})
	return client
}

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	c := setup(t)

	if err := create(tc.ProjectID, subID, topic); err != nil {
		t.Fatalf("failed to create a subscription: %v", err)
	}
	ok, err := c.Subscription(subID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want sub = %q", subID)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

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

		subNames := make([]string, len(subs))
		for i, sub := range subs {
			subNames[i] = sub.ID()
		}
		r.Errorf("got %+v; want a list with subscription %q", subNames, subID)
	})
}

func TestIAM(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

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

// Not sure about the best way to rewrite this test
func TestDelete(t *testing.T) {
	tc := testutil.SystemTest(t)
	setup(t)

	if err := delete(tc.ProjectID, subID); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", subID, err)
	}
	ok, err := client.Subscription(subID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", subID)
	}
}
