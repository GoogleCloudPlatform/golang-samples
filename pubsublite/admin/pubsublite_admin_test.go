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

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/GoogleCloudPlatform/golang-samples/pubsublite/internal/psltest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	resourcePrefix = "admin-test-"
	testRegion     = "us-west1"
)

var (
	supportedZones = []string{"us-west1-a", "us-west1-c"}

	once            sync.Once
	projNumber      string
	reservationID   string
	reservationPath string
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

		psltest.Cleanup(t, client, projNumber, testRegion, resourcePrefix, supportedZones)
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
	t.Run("CreateTopic", func(t *testing.T) {
		ctx := context.Background()
		reservationID = resourcePrefix + uuid.NewString()
		reservationPath = fmt.Sprintf("projects/%s/locations/%s/reservations/%s", projNumber, testRegion, reservationID)
		client.CreateReservation(ctx, pubsublite.ReservationConfig{
			Name:               reservationPath,
			ThroughputCapacity: 4,
		})

		buf := new(bytes.Buffer)
		err := createTopic(buf, tc.ProjectID, testRegion, testZone, topicID, reservationPath)
		if err != nil {
			t.Fatalf("createTopic: %v", err)
		}
		got := buf.String()
		want := "Created topic"
		if !strings.Contains(got, want) {
			t.Fatalf("createTopic() mismatch: got: %s\nwant: %s", got, want)
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
			want := "Got topic"
			if !strings.Contains(got, want) {
				r.Errorf("getTopic() mismatch: got: %s\nwant: %s", got, want)
			}
		})
	})

	t.Run("UpdateTopic", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := updateTopic(buf, projNumber, testRegion, testZone, topicID, reservationPath)
		if err != nil {
			t.Fatalf("updateTopic: %v", err)
		}

		got := buf.String()
		want := "Updated topic"
		if !strings.Contains(got, want) {
			t.Fatalf("updateTopic() mismatch: got: %s\nwant: %s", got, want)
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

	exportSubID := resourcePrefix + "-export-" + uuid.NewString()
	exportSubPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projNumber, testZone, exportSubID)

	// Destination Pub/Sub topic for testing export subscriptions.
	pubsubClient, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create pubsub client: %v", err)
	}
	defer pubsubClient.Close()
	pubsubTopic, err := pubsubClient.CreateTopic(ctx, topicID)
	if err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	defer pubsubTopic.Delete(ctx)

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

	t.Run("CreatePubsubExportSubscription", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := createPubsubExportSubscription(buf, tc.ProjectID, testRegion, testZone, topicID, exportSubID, topicID)
		if err != nil {
			t.Fatalf("createPubsubExportSubscription: %v", err)
		}
		got := buf.String()
		want := fmt.Sprintf("Created export subscription: %s\n", exportSubPath)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("createPubsubExportSubscription() mismatch: -want, +got:\n%s", diff)
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

	t.Run("SeekSubscription", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := seekSubscription(buf, projNumber, testRegion, testZone, subID, pubsublite.Beginning, false)
		if err != nil {
			t.Fatalf("seekSubscription: %v", err)
		}
		got := buf.String()
		want := "Seek operation initiated"
		if !strings.Contains(got, want) {
			t.Fatalf("got: %v, want %v", got, want)
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

	client.DeleteSubscription(ctx, exportSubPath)
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

func TestReservationsAdmin(t *testing.T) {
	t.Parallel()
	client := setupAdmin(t)
	defer client.Close()
	tc := testutil.SystemTest(t)

	reservationID := resourcePrefix + uuid.NewString()
	resPath := fmt.Sprintf("projects/%s/locations/%s/reservations/%s", projNumber, testRegion, reservationID)
	cap := 4
	t.Run("CreateReservation", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := createReservation(buf, tc.ProjectID, testRegion, reservationID, cap)
		if err != nil {
			t.Fatalf("createReservation: %v", err)
		}

		got := buf.String()
		want := "Created reservation"
		if !strings.Contains(got, want) {
			t.Fatalf("createReservation() mismatch: got: %s\nwant: %s", got, want)
		}
	})

	t.Run("GetReservation", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := getReservation(buf, tc.ProjectID, testRegion, reservationID)
		if err != nil {
			t.Fatalf("getReservation: %v", err)
		}

		got := buf.String()
		want := fmt.Sprintf("Got reservation: %#v\n", psltest.DefaultResConfig(resPath))
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("getReservation() mismatch: -want, +got:\n%s", diff)
		}
	})

	t.Run("UpdateReservation", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := updateReservation(buf, tc.ProjectID, testRegion, reservationID, cap)
		if err != nil {
			t.Fatalf("updateReservation: %v", err)
		}

		got := buf.String()
		want := "Updated reservation"
		if !strings.Contains(got, want) {
			t.Fatalf("updateReservation() mismatch: got: %s\nwant: %s", got, want)
		}
	})

	t.Run("DeleteReservation", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := deleteReservation(buf, tc.ProjectID, testRegion, reservationID)
		if err != nil {
			t.Fatalf("deleteReservation: %v", err)
		}

		got := buf.String()
		want := "Deleted reservation"
		if got != want {
			t.Fatalf("got: %v, want %v", got, want)
		}
	})
}

func TestListReservations(t *testing.T) {
	t.Parallel()
	client := setupAdmin(t)
	defer client.Close()
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	var resPaths []string
	for i := 0; i < 3; i++ {
		resID := resourcePrefix + uuid.NewString()
		resPath := fmt.Sprintf("projects/%s/locations/%s/reservations/%s", projNumber, testRegion, resID)
		resPaths = append(resPaths, resPath)
		psltest.MustCreateReservation(ctx, t, client, resPath)
	}

	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := listReservations(buf, tc.ProjectID, testRegion)
		if err != nil {
			r.Errorf("listReservations got err: %v", err)
		}
		got := buf.String()
		for _, rp := range resPaths {
			if !strings.Contains(got, rp) {
				r.Errorf("missing reservation from list: %s", rp)
			}
		}
	})

	for _, rp := range resPaths {
		client.DeleteReservation(ctx, rp)
	}
}

func randomZone() string {
	return supportedZones[rand.Intn(len(supportedZones))]
}
