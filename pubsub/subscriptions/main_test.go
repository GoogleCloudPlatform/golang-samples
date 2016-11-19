// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var topic *pubsub.Topic

const (
	subID   = "golang-samples-subscription"
	topicID = "golang-samples-topic"
)

var once sync.Once // guards cleanup related operations that needs to be executed only for once.

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Cleanup resources from the previous failed tests.
	once.Do(func() {
		// create a topic to subscribe to.
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

		// delete the sub if already exists
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
	c := setup(t)
	if err := create(c, subID, topic); err != nil {
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
	c := setup(t)

	var ok bool
	var subs []*pubsub.Subscription
outer:
	for attempts := 0; attempts < 5; attempts++ {
		var err error
		subs, err = list(c)
		if err != nil {
			t.Fatalf("failed to list subscriptions: %v", err)
		}
		for _, sub := range subs {
			if sub.ID() == subID {
				ok = true
				break outer
			}
		}
		time.Sleep(2 * time.Second)
	}
	if !ok {
		subNames := make([]string, len(subs))
		for i, sub := range subs {
			subNames[i] = sub.ID()
		}
		t.Fatalf("got %+v; want a list with subscription %q", subNames, subID)
	}
}

func TestIAM(t *testing.T) {
	c := setup(t)

	perms := testPermissions(c, subID)
	if len(perms) == 0 {
		t.Fatalf("want non-zero perms")
	}

	addUsers(c, subID)

	policy := getPolicy(c, subID)
	if role, member := iam.Editor, "group:cloud-logs@google.com"; !policy.HasRole(member, role) {
		t.Fatalf("want %q as viewer, got %v", member, policy)
	}
	if role, member := iam.Viewer, iam.AllUsers; !policy.HasRole(member, role) {
		t.Fatalf("want %q as viewer, got %v", member, policy)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)
	if err := delete(c, subID); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", subID, err)
	}
	ok, err := c.Subscription(subID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", subID)
	}
}
