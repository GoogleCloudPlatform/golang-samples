// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const topicName = "golang-samples-subs-example-topic"

var (
	once    sync.Once // guards cleanup related operations that needs to be executed only for once.
	topic   *pubsub.Topic
	now     time.Time
	subname string
)

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	now = time.Now()

	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Cleanup resources from the previous failed tests.
	once.Do(func() {
		// create a topic to subscribe to if it doesn't exist.
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
		subname = fmt.Sprintf("sub-%d", now.Unix())

		// Cleanup subscriptions older than 12 hours.
		subit := client.Subscriptions(ctx)
		subre := regexp.MustCompile("sub-(\\d+)")
		for {
			s, err := subit.Next()
			if err == pubsub.Done {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			// TODO(jbd): Fix ugly string work once
			// https://github.com/GoogleCloudPlatform/gcloud-golang/issues/342 is resolved.
			name := strings.Replace(s.Name(), "projects/"+tc.ProjectID+"/subscriptions/", "", -1)
			groups := subre.FindStringSubmatch(name)
			if len(groups) == 0 {
				continue // not a timestamped subscription
			}
			name = groups[1]
			st, err := strconv.ParseInt(name, 10, 64)
			if err != nil {
				t.Errorf("cannot convert numeric string %q into int: %v", name, err)
			}
			if time.Now().Sub(time.Unix(st, 0)) > 12*time.Hour {
				// delete old garbage subscription
				if err := s.Delete(ctx); err != nil {
					t.Errorf("cannot delete garbage subscription (%q): %v", s.Name(), err)
				}
			}
		}
	})
	return client
}

func TestCreate(t *testing.T) {
	c := setup(t)
	if err := create(c, subname, topic); err != nil {
		t.Fatalf("failed to create a subscription: %v", err)
	}
	ok, err := c.Subscription(subname).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want sub = %q", subname)
	}
}

func TestList(t *testing.T) {
	c := setup(t)

	subs, err := list(c)
	if err != nil {
		t.Fatalf("failed to list subscriptions: %v", err)
	}
	var ok bool
	sub := c.Subscription(subname)
	for _, s := range subs {
		if sub.Name() == s.Name() {
			ok = true
			break
		}
	}
	if !ok {
		t.Fatalf("got a list without create sub; want a list with %v", sub)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)
	if err := delete(c, subname); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", subname, err)
	}
	ok, err := c.Subscription(subname).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", subname)
	}
}
