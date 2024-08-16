// Copyright 2024 Google LLC
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

package notificationsv2

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func orgID(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Fatal("GCLOUD_ORGANIZATION not set")
	}
	orgID = strings.TrimSpace(orgID)
	return orgID
}

func projectID(t *testing.T) string {
	projectID := os.Getenv("SCC_PUBSUB_PROJECT")
	if projectID == "" {
		t.Fatal("SCC_PUBSUB_PROJECT not set")
	}
	projectID = strings.TrimSpace(projectID)
	return projectID
}

func pubsubTopic(t *testing.T) string {
	pubsubTopic := os.Getenv("SCC_PUBSUB_TOPIC")
	if pubsubTopic == "" {
		t.Fatal("SCC_PUBSUB_TOPIC not set")
	}
	pubsubTopic = strings.TrimSpace(pubsubTopic)
	fmt.Printf("PubsubTopic: %v\n", pubsubTopic)
	return pubsubTopic
}

func pubsubSubscription(t *testing.T) string {
	pubsubSubscription := os.Getenv("SCC_PUBSUB_SUBSCRIPTION")
	if pubsubSubscription == "" {
		t.Fatal("SCC_PUBSUB_SUBSCRIPTION not set")
	}
	pubsubSubscription = strings.TrimSpace(pubsubSubscription)
	return pubsubSubscription
}

func createTestNotificationConfig(buf *bytes.Buffer, orgID string, pubsubTopic string, configID string) error {

	projectID := projectIDFromEnv()
	fullPubsubTopic := fmt.Sprintf("projects/%s/topics/%s", projectID, pubsubTopic)
	fmt.Printf("FullPubSubTopic: %v\n", fullPubsubTopic)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.CreateNotificationConfigRequest{
		Parent:   fmt.Sprintf("organizations/%s/locations/global", orgID),
		ConfigId: configID,
		NotificationConfig: &securitycenterpb.NotificationConfig{
			Description: "Go sample config",
			PubsubTopic: fullPubsubTopic,
			NotifyConfig: &securitycenterpb.NotificationConfig_StreamingConfig_{
				StreamingConfig: &securitycenterpb.NotificationConfig_StreamingConfig{
					Filter: `state = "ACTIVE"`,
				},
			},
		},
	}

	_, err = client.CreateNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to create notification config: %w", err)
	}

	buf.WriteString("New NotificationConfig created")
	return nil
}

func projectIDFromEnv() string {
	projectID := os.Getenv("SCC_PUBSUB_PROJECT")
	if projectID == "" {
		panic("SCC_PUBSUB_PROJECT not set")
	}
	return projectID
}

func cleanupNotificationConfig(t *testing.T, notificationConfigID string) error {
	orgID := orgID(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("organizations/%s/notificationConfigs/%s", orgID, notificationConfigID)
	req := &securitycenterpb.DeleteNotificationConfigRequest{
		Name: name,
	}

	if err = client.DeleteNotificationConfig(ctx, req); err != nil {
		return fmt.Errorf("Failed to delete notification config: %w", err)
	}

	return nil
}

func TestCreateNotificationConfig(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "go-test-create-config-id" + rand.String()

		if err := createTestNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("createTestNotificationConfig failed: %v", err)
			return
		}

		defer cleanupNotificationConfig(t, configID)

		if !strings.Contains(buf.String(), "New NotificationConfig created") {
			r.Errorf("createTestNotificationConfig did not create.")
		}
	})
}

func TestDeleteNotificationConfig(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "go-test-delete-config-id" + rand.String()

		if err := createTestNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		defer cleanupNotificationConfig(t, configID)

		if err := deleteNotificationConfig(buf, orgID(t), configID); err != nil {
			r.Errorf("deleteNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Deleted config:") {
			r.Errorf("deleteNotificationConfig did not delete.")
		}
	})
}

func TestGetNotificationConfig(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "go-test-get-config-id" + rand.String()

		if err := createTestNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		defer cleanupNotificationConfig(t, configID)

		if err := getNotificationConfig(buf, orgID(t), configID); err != nil {
			r.Errorf("getNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Received config:") {
			r.Errorf("getNotificationConfig did not retrieve.")
		}
	})
}

func TestListNotificationConfigs(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "go-test-list-config-id" + rand.String()

		if err := createTestNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		defer cleanupNotificationConfig(t, configID)

		if err := listNotificationConfigs(buf, orgID(t)); err != nil {
			r.Errorf("listNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "NotificationConfig") {
			r.Errorf("listNotificationConfigs did not list")
		}
	})
}

func TestUpdateNotificationConfig(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "go-test-update-config-id" + rand.String()
		projectID := projectIDFromEnv()

		if err := createTestNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		defer cleanupNotificationConfig(t, configID)

		if err := updateNotificationConfig(buf, orgID(t), configID, pubsubTopic(t), projectID); err != nil {
			r.Errorf("updateNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Updated NotificationConfig:") {
			r.Errorf("updateNotificationConfig did not update.")
		}
	})
}

func TestReceiveNotifications(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "go-test-receive-config-id" + rand.String()

		if err := createTestNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		defer cleanupNotificationConfig(t, configID)

		// Ensure a notification is sent before receiving
		if err := sendTestNotification(pubsubTopic(t)); err != nil {
			r.Errorf("sendTestNotification failed: %v", err)
			return
		}

		if err := receiveMessages(buf, projectID(t), pubsubSubscription(t)); err != nil {
			r.Errorf("receiveNotifications failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Got finding") {
			r.Errorf("Did not receive any notifications.")
			return
		}
	})
}

func sendTestNotification(pubsubTopic string) error {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, os.Getenv("SCC_PUBSUB_PROJECT"))
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	topic := client.Topic(pubsubTopic)

	msg := &pubsub.Message{
		Data: []byte("Test notification"),
	}

	result := topic.Publish(ctx, msg)

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("result.Get: %v", err)
	}

	return nil
}
