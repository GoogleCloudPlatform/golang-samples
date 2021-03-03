// Copyright 2021 Google LLC
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
package pslite

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsublite"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"
)

const (
	topic      = "test-topic-"
	sub        = "test-sub-"
	testRegion = "us-central1"
)

var (
	supportedZoneIDs = []string{"a", "b", "c"}

	once       sync.Once
	projNumber string
)

func setupAdmin(t *testing.T) *pubsublite.AdminClient {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := pubsublite.NewAdminClient(ctx, testRegion)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		// PubSub Lite returns project numbers in resource paths, so we need to convert from project id
		// to numbers for tests.
		cloudresourcemanagerService, err := cloudresourcemanager.NewService(context.Background())
		if err != nil {
			t.Fatalf("cloudresourcemanager.NewService: %v", err)
		}

		project, err := cloudresourcemanagerService.Projects.Get(tc.ProjectID).Do()
		if err != nil {
			t.Fatalf("cloudresourcemanagerService.Projets.Get project: %v", err)
		}

		projNumber = strconv.FormatInt(project.ProjectNumber, 10)

		cleanup(t, client, projNumber)
	})

	return client
}

// cleanup deletes all previous test topics/subscriptions from
// previous test runs. This prevents previous test failures
// from building up resources that count against quota.
func cleanup(t *testing.T, client *pubsublite.AdminClient, proj string) {
	ctx := context.Background()

	for _, zoneID := range supportedZoneIDs {
		zone := fmt.Sprintf("%s-%s", testRegion, zoneID)
		parent := fmt.Sprintf("projects/%s/locations/%s", proj, zone)
		topicIter := client.Topics(ctx, parent)
		for {
			topic, err := topicIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("topicIter.Next got err: %v", err)
			}
			if err := client.DeleteTopic(ctx, topic.Name); err != nil {
				t.Fatalf("client.DeleteTopic got err: %v", err)
			}
		}

		subIter := client.Subscriptions(ctx, parent)
		for {
			sub, err := subIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("subIter.Next() got err: %v", err)
			}
			if err := client.DeleteSubscription(ctx, sub.Name); err != nil {
				t.Fatalf("client.DeleteSubscription got err: %v", err)
			}
		}
	}
}

func TestCreateTopic(t *testing.T) {
	client := setupAdmin(t)
	tc := testutil.SystemTest(t)
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	buf := new(bytes.Buffer)
	err := createTopic(buf, tc.ProjectID, testRegion, testZone, topicID)
	if err != nil {
		t.Fatalf("createTopic: %v", err)
	}
	got := buf.String()
	want := fmt.Sprintf("Created topic: %s\n", topicPath)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("createTopic() mismatch: -want, +got:\n%s", diff)
	}

	t.Cleanup(func() {
		client.DeleteTopic(context.Background(), topicPath)
	})
}

func TestGetTopic(t *testing.T) {
	client := setupAdmin(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	buf := new(bytes.Buffer)
	err := getTopic(buf, tc.ProjectID, testRegion, testZone, topicID)
	if err != nil {
		t.Fatalf("getTopic: %v", err)
	}
	got := buf.String()
	want := fmt.Sprintf("Got topic: %#v\n", *defaultTopicConfig(topicPath))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("getTopic() mismatch: -want, +got:\n%s", diff)
	}

	t.Cleanup(func() {
		client.DeleteTopic(ctx, topicPath)
	})
}

func TestListTopics(t *testing.T) {
	client := setupAdmin(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	var topicPaths []string
	for i := 0; i < 3; i++ {
		topicID := topic + uuid.NewString()
		topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
		topicPaths = append(topicPaths, topicPath)
		mustCreateTopic(ctx, t, client, topicPath)
	}

	buf := new(bytes.Buffer)
	err := listTopics(buf, tc.ProjectID, testRegion, testZone)
	if err != nil {
		t.Fatalf("listTopics got err: %v", err)
	}
	got := buf.String()
	for _, tp := range topicPaths {
		if !strings.Contains(got, tp) {
			t.Fatalf("missing topic path from list: %s", tp)
		}
	}

	t.Cleanup(func() {
		for _, tp := range topicPaths {
			client.DeleteTopic(ctx, tp)
		}
	})
}

func TestUpdateTopic(t *testing.T) {
	client := setupAdmin(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	buf := new(bytes.Buffer)
	err := updateTopic(buf, projNumber, testRegion, testZone, topicID)
	if err != nil {
		t.Fatalf("updateTopic: %v", err)
	}

	got := buf.String()
	// This is hard coded into the pubsublite/update_topic.go sample.
	// If the sample value changes, this value needs to change as well.
	wantCfg := &pubsublite.TopicConfig{
		Name:                       topicPath,
		PartitionCount:             3,
		PublishCapacityMiBPerSec:   8,
		SubscribeCapacityMiBPerSec: 16,
		PerPartitionBytes:          60 * 1024 * 1024 * 1024,
		RetentionDuration:          24 * time.Hour,
	}
	want := fmt.Sprintf("Updated topic: %#v\n", *wantCfg)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("updateTopic() mismatch: -want, +got:\n%s", diff)
	}

	t.Cleanup(func() {
		client.DeleteTopic(context.Background(), topicPath)
	})
}

func TestDeleteTopic(t *testing.T) {
	client := setupAdmin(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	buf := new(bytes.Buffer)
	err := deleteTopic(buf, tc.ProjectID, testRegion, testZone, topicID)
	if err != nil {
		t.Fatalf("deleteTopic: %v", err)
	}

	got := buf.String()
	want := "Deleted topic\n"
	if got != want {
		t.Fatalf("got: %v, want %v", got, want)
	}
}

func mustCreateTopic(ctx context.Context, t *testing.T, client *pubsublite.AdminClient, topicPath string) *pubsublite.TopicConfig {
	cfg := defaultTopicConfig(topicPath)
	topicConfig, err := client.CreateTopic(ctx, *cfg)
	if err != nil {
		t.Fatalf("AdminClient.CreateTopic got err: %v", err)
	}
	return topicConfig
}

func defaultTopicConfig(topicPath string) *pubsublite.TopicConfig {
	cfg := &pubsublite.TopicConfig{
		Name:                       topicPath,
		PartitionCount:             2,
		PublishCapacityMiBPerSec:   4,
		SubscribeCapacityMiBPerSec: 4,
		PerPartitionBytes:          30 * 1024 * 1024 * 1024, // 30 GiB
		RetentionDuration:          pubsublite.InfiniteRetention,
	}
	return cfg
}

func TestCreateSubscription(t *testing.T) {
	client := setupAdmin(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	subID := sub + uuid.NewString()
	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)

	buf := new(bytes.Buffer)
	err := createSubscription(buf, tc.ProjectID, testRegion, testZone, topicID, subID)
	if err != nil {
		t.Fatalf("createSubscription: %v", err)
	}
	got := buf.String()
	want := fmt.Sprintf("Created subscription: %s\n", subPath)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("createSubscription() mismatch: -want, +got:\n%s", diff)
	}

	t.Cleanup(func() {
		client.DeleteTopic(ctx, topicPath)
		client.DeleteSubscription(ctx, subPath)
	})
}

func TestGetSubscription(t *testing.T) {
	client := setupAdmin(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	subID := sub + uuid.NewString()
	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)
	mustCreateSubscription(ctx, t, client, topicPath, subPath)

	buf := new(bytes.Buffer)
	err := getSubscription(buf, projNumber, testRegion, testZone, subID)
	if err != nil {
		t.Fatalf("getSubscription: %v", err)
	}
	got := buf.String()
	want := fmt.Sprintf("Got subscription: %#v\n", defaultSubConfig(topicPath, subPath))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("getSubscription mismatch: -want, +got:\n%s", diff)
	}

	t.Cleanup(func() {
		client.DeleteTopic(ctx, topicPath)
		client.DeleteSubscription(ctx, subPath)
	})
}

func TestListSubscriptions(t *testing.T) {
	client := setupAdmin(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	var subPaths []string
	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	for i := 0; i < 3; i++ {
		subID := sub + uuid.NewString()
		subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)
		mustCreateSubscription(ctx, t, client, topicPath, subPath)
		subPaths = append(subPaths, subPath)
	}

	// Test listSubscriptionsInProject.
	buf := new(bytes.Buffer)
	err := listSubscriptionsInProject(buf, tc.ProjectID, testRegion, testZone)
	if err != nil {
		t.Fatalf("listSubscriptionsInProject got err: %v", err)
	}
	got := buf.String()
	for _, sp := range subPaths {
		if !strings.Contains(got, sp) {
			t.Fatalf("missing sub path from list: %s", sp)
		}
	}

	// Test listSubscriptionsInTopic with same list of subscriptions.
	buf = new(bytes.Buffer)
	err = listSubscriptionsInTopic(buf, tc.ProjectID, testRegion, testZone, topicID)
	if err != nil {
		t.Fatalf("listSubscriptionsInTopic got err: %v", err)
	}
	got = buf.String()
	for _, sp := range subPaths {
		if !strings.Contains(got, sp) {
			t.Fatalf("missing sub path from list: %s", sp)
		}
	}

	t.Cleanup(func() {
		for _, sp := range subPaths {
			client.DeleteTopic(ctx, sp)
		}
	})
}

func TestUpdateSubscription(t *testing.T) {
	client := setupAdmin(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	subID := sub + uuid.NewString()
	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)
	mustCreateSubscription(ctx, t, client, topicPath, subPath)

	buf := new(bytes.Buffer)
	err := updateSubscription(buf, projNumber, testRegion, testZone, subID)
	if err != nil {
		t.Fatalf("updateSubscription: %v", err)
	}
	got := buf.String()
	// This is hard coded into the pubsublite/update_subscription.go sample.
	// If the sample value changes, this value needs to change as well.
	wantCfg := &pubsublite.SubscriptionConfig{
		Name:                subPath,
		Topic:               topicPath,
		DeliveryRequirement: pubsublite.DeliverAfterStored,
	}
	want := fmt.Sprintf("Updated subscription: %#v\n", wantCfg)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("updateSubscription() mismatch: -want, +got:\n%s", diff)
	}

	t.Cleanup(func() {
		client.DeleteTopic(ctx, topicPath)
		client.DeleteSubscription(ctx, subPath)
	})
}

func TestDeleteSubscription(t *testing.T) {
	client := setupAdmin(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := topic + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	mustCreateTopic(ctx, t, client, topicPath)

	subID := sub + uuid.NewString()
	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)
	mustCreateSubscription(ctx, t, client, topicPath, subPath)

	buf := new(bytes.Buffer)
	err := deleteSubscription(buf, projNumber, testRegion, testZone, subID)
	if err != nil {
		t.Fatalf("deleteSubscription: %v", err)
	}
	got := buf.String()
	want := "Deleted subscription\n"
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	t.Cleanup(func() {
		client.DeleteTopic(ctx, topicPath)
		client.DeleteSubscription(ctx, subPath)
	})
}

func mustCreateSubscription(ctx context.Context, t *testing.T, client *pubsublite.AdminClient, topicPath, subPath string) *pubsublite.SubscriptionConfig {
	cfg := defaultSubConfig(topicPath, subPath)
	subConfig, err := client.CreateSubscription(ctx, *cfg)
	if err != nil {
		t.Fatalf("AdminClient.CreateSubscription got err: %v", err)
	}
	return subConfig
}

func defaultSubConfig(topicPath, subPath string) *pubsublite.SubscriptionConfig {
	cfg := &pubsublite.SubscriptionConfig{
		Name:                subPath,
		Topic:               topicPath,
		DeliveryRequirement: pubsublite.DeliverImmediately,
	}
	return cfg
}

func randomZone() string {
	zoneID := supportedZoneIDs[rand.Intn(len(supportedZoneIDs))]
	return fmt.Sprintf("%s-%s", testRegion, zoneID)
}
