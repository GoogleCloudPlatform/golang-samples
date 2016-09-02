// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/net/context"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const topicName = "golang-samples-topic-example"

func setup(t *testing.T) *pubsub.Client {
	netctx := context.Background()
	ctx := testutil.SystemTest(t)

	client, err := pubsub.NewClient(netctx, ctx.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestCreate(t *testing.T) {
	c := setup(t)
	if err := create(c, topicName); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
	ok, err := c.Topic(topicName).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want topic = %q", topicName)
	}
}

func TestList(t *testing.T) {
	c := setup(t)
	topics, err := list(c)
	if err != nil {
		t.Fatalf("failed to list topics: %v", err)
	}
	var ok bool
	topic := c.Topic(topicName)
	for _, t := range topics {
		// TODO(jbd): Fix HasSuffix when
		if t.Name() == topic.Name() {
			ok = true
			break
		}
	}
	if !ok {
		t.Errorf("got %+v; want the list with topic = %q", topics, topicName)
	}
}

func TestPublish(t *testing.T) {
	// Nothing much to do here, unless we are consuming.
	// TODO(jbd): Merge topics and subscriptions programs maybe?
	c := setup(t)
	if err := publish(c, topicName, "hello world"); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)
	if err := delete(c, topicName); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", topicName, err)
	}
	ok, err := c.Topic(topicName).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", topicName)
	}
}
