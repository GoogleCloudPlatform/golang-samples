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

package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var topicID string

var once sync.Once // guards cleanup related operations in setup.

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	topicID = tc.ProjectID + "-test-topic"

	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Cleanup resources from the previous failed tests.
	once.Do(func() {
		topic := client.Topic(topicID)
		ok, err := topic.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if topic exists: %v", err)
		}
		if !ok {
			return
		}
		if err := topic.Delete(ctx); err != nil {
			t.Fatalf("failed to cleanup the topic (%q): %v", topicID, err)
		}
	})
	return client
}

func TestCreate(t *testing.T) {
	c := setup(t)
	if err := create(c, topicID); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
	ok, err := c.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want topic = %q", topicID)
	}
}

func TestList(t *testing.T) {
	c := setup(t)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		topics, err := list(c)
		if err != nil {
			r.Errorf("failed to list topics: %v", err)
		}

		for _, t := range topics {
			if t.ID() == topicID {
				return // PASS
			}
		}

		topicNames := make([]string, len(topics))
		for i, t := range topics {
			topicNames[i] = t.ID()
		}
		r.Errorf("got %+v; want a list with topic = %q", topicNames, topicID)
	})
}

func TestPublish(t *testing.T) {
	// Nothing much to do here, unless we are consuming.
	// TODO(jbd): Merge topics and subscriptions programs maybe?
	c := setup(t)
	if err := publish(c, topicID, "hello world"); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishThatScales(t *testing.T) {
	c := setup(t)
	if err := publishThatScales(c, topicID, 10); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestPublishCustomAttributes(t *testing.T) {
	c := setup(t)
	if err := publishCustomAttributes(c, topicID); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestIAM(t *testing.T) {
	c := setup(t)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		perms, err := testPermissions(c, topicID)
		if err != nil {
			r.Errorf("testPermissions: %v", err)
		}
		if len(perms) == 0 {
			r.Errorf("want non-zero perms")
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := addUsers(c, topicID); err != nil {
			r.Errorf("addUsers: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		policy, err := getPolicy(c, topicID)
		if err != nil {
			r.Errorf("getPolicy: %v", err)
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
	c := setup(t)
	if err := delete(c, topicID); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", topicID, err)
	}
	ok, err := c.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", topicID)
	}
}
