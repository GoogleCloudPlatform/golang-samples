// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var topic *pubsub.Topic

const (
	subName   = "golang-samples-subscription"
	topicName = "golang-samples-topic"
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
		topic = client.Topic(topicName)
		ok, err := topic.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if topic exists: %v", err)
		}
		if !ok {
			if topic, err = client.CreateTopic(ctx, topicName); err != nil {
				t.Fatalf("failed to create the topic: %v", err)
			}
		}

		// delete the sub if already exists
		sub := client.Subscription(subName)
		ok, err = sub.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if sub exists: %v", err)
		}
		if ok {
			if err := client.Subscription(subName).Delete(ctx); err != nil {
				t.Fatalf("failed to cleanup the topic (%q): %v", subName, err)
			}
		}
	})
	return client
}

func TestCreate(t *testing.T) {
	c := setup(t)
	if err := create(c, subName, topic); err != nil {
		t.Fatalf("failed to create a subscription: %v", err)
	}
	ok, err := c.Subscription(subName).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want sub = %q", subName)
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
		s := c.Subscription(subName)
		for _, sub := range subs {
			if s.Name() == sub.Name() {
				ok = true
				break outer
			}
		}
		time.Sleep(2 * time.Second)
	}
	if !ok {
		subNames := make([]string, len(subs))
		for i, sub := range subs {
			subNames[i] = sub.Name()
		}
		t.Fatalf("got %+v; want a list with subscription %q", subNames, subName)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)
	if err := delete(c, subName); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", subName, err)
	}
	ok, err := c.Subscription(subName).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", subName)
	}
}
