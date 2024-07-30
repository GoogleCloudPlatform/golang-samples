// Copyright 2020 Google LLC
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
	"os"
	"strings"
	"testing"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func orgID(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	return orgID
}

func projectID(t *testing.T) string {
	projectID := os.Getenv("SCC_PUBSUB_PROJECT")
	if projectID == "" {
		t.Skip("SCC_PUBSUB_PROJECT not set")
	}
	return projectID
}

func pubsubTopic(t *testing.T) string {
	pubsubTopic := os.Getenv("SCC_PUBSUB_TOPIC")
	if pubsubTopic == "" {
		t.Skip("SCC_PUBSUB_TOPIC not set")
	}
	return pubsubTopic
}

func pubsubSubscription(t *testing.T) string {
	pubsubSubscription := os.Getenv("SCC_PUBSUB_SUBSCRIPTION")
	if pubsubSubscription == "" {
		t.Skip("SCC_PUBSUB_SUBSCRIPTION not set")
	}
	return pubsubSubscription
}

func addNotificationConfig(t *testing.T, notificationConfigID string) error {
	orgID := orgID(t)
	pubsubTopic := pubsubTopic(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.CreateNotificationConfigRequest{
		Parent:   fmt.Sprintf("organizations/%s", orgID),
		ConfigId: notificationConfigID,
		NotificationConfig: &securitycenterpb.NotificationConfig{
			Description: "Go sample config",
			PubsubTopic: pubsubTopic,
			NotifyConfig: &securitycenterpb.NotificationConfig_StreamingConfig_{
				StreamingConfig: &securitycenterpb.NotificationConfig_StreamingConfig{
					Filter: `state = "ACTIVE"`,
				},
			},
		},
	}

	_, err0 := client.CreateNotificationConfig(ctx, req)
	if err0 != nil {
		return fmt.Errorf("Failed to create notification config: %w", err0)
	}

	return nil
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
		return fmt.Errorf("Failed to retrieve notification config: %w", err)
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

		if err := createNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
			r.Errorf("createNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "New NotificationConfig created") {
			r.Errorf("createNotificationConfig did not create.")
		}

		cleanupNotificationConfig(t, configID)
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

		if err := addNotificationConfig(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

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

		if err := addNotificationConfig(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		if err := getNotificationConfig(buf, orgID(t), configID); err != nil {
			r.Errorf("getNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Received config:") {
			r.Errorf("getNotificationConfig did not delete.")
		}

		cleanupNotificationConfig(t, configID)
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

		if err := addNotificationConfig(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		if err := listNotificationConfigs(buf, orgID(t)); err != nil {
			r.Errorf("listNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "NotificationConfig") {
			r.Errorf("listNotificationConfigs did not list")
		}

		cleanupNotificationConfig(t, configID)
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

		if err := addNotificationConfig(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		if err := updateNotificationConfig(buf, orgID(t), configID, pubsubTopic(t)); err != nil {
			r.Errorf("updateNotificationConfig failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Updated NotificationConfig:") {
			r.Errorf("updateNotificationConfig did not update.")
		}
		cleanupNotificationConfig(t, configID)
	})
}

func TestReceiveNotifications(t *testing.T) {
	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if err := receiveMessages(buf, projectID(t), pubsubSubscription(t)); err != nil {
			r.Errorf("receiveNotifications failed: %v", err)
			return
		}

		if !strings.Contains(buf.String(), "Got finding") {
			r.Errorf("Did not receive any notifications.")
		}
	})
}
