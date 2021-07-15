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

package admin

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
	"github.com/GoogleCloudPlatform/golang-samples/pubsublite/internal/psltest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"google.golang.org/api/cloudresourcemanager/v1"
)

const (
	resourcePrefix = "admin-test-"
	testRegion     = "us-central1"
)

var (
	supportedZones = []string{"us-central1-a", "us-central1-b", "us-central1-c"}

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
		// Pub/Sub Lite returns project numbers in resource paths, so we need to convert from project id
		// to numbers for tests.
		crm, err := cloudresourcemanager.NewService(context.Background())
		if err != nil {
			t.Fatalf("cloudresourcemanager.NewService: %v", err)
		}

		project, err := crm.Projects.Get(tc.ProjectID).Do()
		if err != nil {
			t.Fatalf("crm.Projects.Get project: %v", err)
		}

		projNumber = strconv.FormatInt(project.ProjectNumber, 10)

		psltest.Cleanup(t, client, projNumber, resourcePrefix, supportedZones)
	})

	return client
}

func TestTopicAdmin(t *testing.T) {
	t.Parallel()
	client := setupAdmin(t)
	defer client.Close()
	tc := testutil.SystemTest(t)
	testZone := randomZone()

	topicID := resourcePrefix + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	t.Run("CreateTopic", func(t *testing.T) {
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
	})

	t.Run("GetTopic", func(t *testing.T) {
		testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := getTopic(buf, tc.ProjectID, testRegion, testZone, topicID)
			if err != nil {
				r.Errorf("getTopic: %v", err)
			}
			got := buf.String()
			want := fmt.Sprintf("Got topic: %#v\n", *psltest.DefaultTopicConfig(topicPath))
			if diff := cmp.Diff(want, got); diff != "" {
				r.Errorf("getTopic() mismatch: -want, +got:\n%s", diff)
			}
		})
	})

	t.Run("UpdateTopic", func(t *testing.T) {
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
	})

	t.Run("DeleteTopic", func(t *testing.T) {
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
	})
}

func TestListTopics(t *testing.T) {
	t.Parallel()
	client := setupAdmin(t)
	defer client.Close()
	tc := testutil.SystemTest(t)
	testZone := randomZone()
	ctx := context.Background()

	var topicPaths []string
	for i := 0; i < 3; i++ {
		topicID := resourcePrefix + uuid.NewString()
		topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
		topicPaths = append(topicPaths, topicPath)
		psltest.MustCreateTopic(ctx, t, client, topicPath)
	}

	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := listTopics(buf, tc.ProjectID, testRegion, testZone)
		if err != nil {
			r.Errorf("listTopics got err: %v", err)
		}
		got := buf.String()
		for _, tp := range topicPaths {
			if !strings.Contains(got, tp) {
				r.Errorf("missing topic path from list: %s", tp)
			}
		}
	})

	for _, tp := range topicPaths {
		client.DeleteTopic(ctx, tp)
	}
}

func TestSubscriptionAdmin(t *testing.T) {
	t.Parallel()
	client := setupAdmin(t)
	defer client.Close()
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	topicID := resourcePrefix + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)

	psltest.MustCreateTopic(ctx, t, client, topicPath)

	subID := resourcePrefix + uuid.NewString()
	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)

	t.Run("CreateSubscription", func(t *testing.T) {
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
	})

	t.Run("GetSubscription", func(t *testing.T) {
		testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := getSubscription(buf, projNumber, testRegion, testZone, subID)
			if err != nil {
				r.Errorf("getSubscription: %v", err)
			}
			got := buf.String()
			want := fmt.Sprintf("Got subscription: %#v\n", psltest.DefaultSubConfig(topicPath, subPath))
			if diff := cmp.Diff(want, got); diff != "" {
				r.Errorf("getSubscription mismatch: -want, +got:\n%s", diff)
			}
		})
	})

	t.Run("UpdateSubscription", func(t *testing.T) {
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
	})

	t.Run("DeleteSubscription", func(t *testing.T) {
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
	})

	client.DeleteTopic(ctx, topicPath)
}

func TestListSubscriptions(t *testing.T) {
	t.Parallel()
	client := setupAdmin(t)
	defer client.Close()
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	testZone := randomZone()

	var subPaths []string
	topicID := resourcePrefix + uuid.NewString()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projNumber, testZone, topicID)
	psltest.MustCreateTopic(ctx, t, client, topicPath)

	for i := 0; i < 3; i++ {
		subID := resourcePrefix + uuid.NewString()
		subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, subID)
		psltest.MustCreateSubscription(ctx, t, client, topicPath, subPath)
		subPaths = append(subPaths, subPath)
	}

	t.Run("ListSubscriptionsInProject", func(t *testing.T) {
		testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := listSubscriptionsInProject(buf, tc.ProjectID, testRegion, testZone)
			if err != nil {
				r.Errorf("listSubscriptionsInProject got err: %v", err)
			}
			got := buf.String()
			for _, sp := range subPaths {
				if !strings.Contains(got, sp) {
					r.Errorf("missing sub path from list: %s", sp)
				}
			}
		})
	})

	// Test listSubscriptionsInTopic with same list of subscriptions.
	t.Run("ListSubscriptionsInTopic", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := listSubscriptionsInTopic(buf, tc.ProjectID, testRegion, testZone, topicID)
		if err != nil {
			t.Fatalf("listSubscriptionsInTopic got err: %v", err)
		}
		got := buf.String()
		for _, sp := range subPaths {
			if !strings.Contains(got, sp) {
				t.Fatalf("missing sub path from list: %s", sp)
			}
		}
	})

	client.DeleteTopic(ctx, topicPath)
	for _, sp := range subPaths {
		client.DeleteSubscription(ctx, sp)
	}
}

func randomZone() string {
	return supportedZones[rand.Intn(len(supportedZones))]
}
