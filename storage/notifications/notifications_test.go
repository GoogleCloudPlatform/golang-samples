// Copyright 2022 Google LLC
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

package notifications

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	testPrefix      = "test-gcs-go-notifications"
	bucketExpiryAge = time.Hour * 24
)

var (
	client       *storage.Client
	pubsubClient *pubsub.Client
	serviceAgent string
)

func TestMain(m *testing.M) {
	// Initialize global vars
	tc, _ := testutil.ContextMain(m)

	ctx := context.Background()

	// Init clients
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	c.SetRetry(storage.WithPolicy(storage.RetryAlways))
	client = c
	defer client.Close()

	pubsubClient, err = pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer pubsubClient.Close()

	// Get service agent for storage client
	serviceAgent, err = client.ServiceAccount(ctx, tc.ProjectID)
	if err != nil {
		log.Fatalf("ServiceAccount: %v", err)
	}

	// Run tests
	exit := m.Run()

	// Delete old buckets whose name begins with our test prefix
	if err := testutil.DeleteExpiredBuckets(client, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
		// Don't fail the test if cleanup fails
		log.Printf("Post-test cleanup failed: %v", err)
	}
	os.Exit(exit)
}

func TestNotifications(t *testing.T) {
	tc := testutil.SystemTest(t)
	topicName := testPrefix + "-topic"
	topicName2 := testPrefix + "-topic-2"

	// Set up resources for test.
	ctx := context.Background()
	createTestTopic(t, tc.ProjectID, topicName)
	createTestTopic(t, tc.ProjectID, topicName2)

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	var buf bytes.Buffer

	// Test Create.
	t.Run("create bucket notification", func(t *testing.T) {
		testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
			if err := createBucketNotification(&buf, tc.ProjectID, bucketName, topicName); err != nil {
				r.Errorf("createBucketNotification: %v", err)
				return
			}

			if got, want := buf.String(), "created notification with ID"; !strings.Contains(got, want) {
				r.Errorf("createBucketNotification: got %q; want to contain %q", got, want)
			}
		})
	})

	// Add another notification so we know its ID
	notification2, err := client.Bucket(bucketName).AddNotification(ctx, &storage.Notification{
		TopicID:        topicName2,
		TopicProjectID: tc.ProjectID,
		PayloadFormat:  storage.JSONPayload,
	})
	if err != nil {
		t.Fatalf("Bucket.AddNotification: %v", err)
	}

	// Test Get.
	t.Run("print bucket notification", func(t *testing.T) {
		buf.Reset()
		if err := printPubsubBucketNotification(&buf, bucketName, notification2.ID); err != nil {
			t.Fatalf("printPubsubBucketNotification: %v", err)
		}

		if got, want := buf.String(), "TopicID:"+topicName2; !strings.Contains(got, want) {
			t.Errorf("printPubsubBucketNotification: got %q; want to contain %q", got, want)
		}
	})

	// Test List.
	t.Run("list bucket notifications", func(t *testing.T) {
		buf.Reset()
		if err := listBucketNotifications(&buf, bucketName); err != nil {
			t.Fatalf("listBucketNotifications: %v", err)
		}

		got := buf.String()
		var want, want2 bytes.Buffer
		fmt.Fprintf(&want, "Notification topic %s with ID", topicName) // we don't know its ID
		fmt.Fprintf(&want2, "Notification topic %s with ID %s", topicName2, notification2.ID)

		if !strings.Contains(got, want.String()) || !strings.Contains(got, want2.String()) {
			t.Errorf("listBucketNotifications: got %q; want to contain the following:\n\t%q\n\t%q", got, want, want2)
		}
	})

	// Test Delete.
	t.Run("delete bucket notification", func(t *testing.T) {
		buf.Reset()
		if err := deleteBucketNotification(&buf, bucketName, notification2.ID); err != nil {
			t.Fatalf("deleteBucketNotification: %v", err)
		}

		if got, want := buf.String(), "deleted notification with ID "+notification2.ID; !strings.Contains(got, want) {
			t.Errorf("deleteBucketNotification: got %q; want to contain %q", got, want)
		}
	})
}

// Creates a pubsub topic for testing and registers a cleanup func to remove it
// once the test finishes
func createTestTopic(t *testing.T, projectID, topicName string) {
	t.Helper()
	ctx := context.Background()
	topic := pubsubClient.Topic(topicName)

	// Create the topic if it doesn't exist.
	testutil.Retry(t, 10, time.Millisecond, func(r *testutil.R) {
		exists, err := topic.Exists(ctx)
		if err != nil {
			r.Errorf("topic.Exists: %v", err)
		}
		if !exists {
			_, err = pubsubClient.CreateTopic(ctx, topicName)
			if err != nil {
				r.Errorf("topic.CreateTopic: %v", err)
			}
		}
	})

	// Add the service agent to the topic's permissions so we can access it from storage.
	policy, err := topic.IAM().Policy(ctx)
	if err != nil {
		t.Errorf("Policy: %v", err)
	}
	policy.Add("serviceAccount:"+serviceAgent, iam.Editor)

	if err := topic.IAM().SetPolicy(ctx, policy); err != nil {
		t.Errorf("SetPolicy: %v", err)
	}

	t.Cleanup(func() {
		if err := topic.Delete(ctx); err != nil {
			t.Errorf("Delete: %v", err)
		}
	})
}
