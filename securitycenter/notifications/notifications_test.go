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

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/google/uuid"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
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
		return fmt.Errorf("securitycenter.NewClient: %v", err)
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
		return fmt.Errorf("Failed to create notification config: %v", err0)
	}

	return nil
}

func cleanupNotificationConfig(t *testing.T, notificationConfigID string) error {
	orgID := orgID(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	name := fmt.Sprintf("organizations/%s/notificationConfigs/%s", orgID, notificationConfigID)
	req := &securitycenterpb.DeleteNotificationConfigRequest{
		Name: name,
	}

	if err = client.DeleteNotificationConfig(ctx, req); err != nil {
		return fmt.Errorf("Failed to retrieve notification config: %v", err)
	}

	return nil
}

func TestCreateNotificationConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Issue generating id.")
	}
	configID := "go-test-create-config-id" + rand.String()

	if err := createNotificationConfig(buf, orgID(t), pubsubTopic(t), configID); err != nil {
		t.Fatalf("createNotificationConfig failed: %v", err)
	}

	if !strings.Contains(buf.String(), "New NotificationConfig created") {
		t.Errorf("createNotificationConfig did not create.")
	}

	cleanupNotificationConfig(t, configID)
}

func TestDeleteNotificationConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Issue generating id.")
	}
	configID := "go-test-delete-config-id" + rand.String()

	if err := addNotificationConfig(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
	}

	if err := deleteNotificationConfig(buf, orgID(t), configID); err != nil {
		t.Fatalf("deleteNotificationConfig failed: %v", err)
	}

	if !strings.Contains(buf.String(), "Deleted config:") {
		t.Errorf("deleteNotificationConfig did not delete.")
	}
}

func TestGetNotificationConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Issue generating id.")
	}
	configID := "go-test-get-config-id" + rand.String()

	if err := addNotificationConfig(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
	}

	if err := getNotificationConfig(buf, orgID(t), configID); err != nil {
		t.Fatalf("getNotificationConfig failed: %v", err)
	}

	if !strings.Contains(buf.String(), "Received config:") {
		t.Errorf("getNotificationConfig did not delete.")
	}

	cleanupNotificationConfig(t, configID)
}

func TestListNotificationConfigs(t *testing.T) {
	buf := new(bytes.Buffer)
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Issue generating id.")
	}
	configID := "go-test-list-config-id" + rand.String()

	if err := addNotificationConfig(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
	}

	if err := listNotificationConfigs(buf, orgID(t)); err != nil {
		t.Fatalf("listNotificationConfig failed: %v", err)
	}

	if !strings.Contains(buf.String(), "NotificationConfig") {
		t.Errorf("listNotificationConfigs did not list")
	}

	cleanupNotificationConfig(t, configID)
}

func TestUpdateNotificationConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Issue generating id.")
	}
	configID := "go-test-update-config-id" + rand.String()

	if err := addNotificationConfig(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
	}

	if err := updateNotificationConfig(buf, orgID(t), configID, pubsubTopic(t)); err != nil {
		t.Fatalf("updateNotificationConfig failed: %v", err)
	}

	if !strings.Contains(buf.String(), "Updated NotificationConfig:") {
		t.Errorf("updateNotificationConfig did not update.")
	}

	cleanupNotificationConfig(t, configID)
}

func TestReceiveNotifications(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := receiveMessages(buf, projectID(t), pubsubSubscription(t)); err != nil {
		t.Fatalf("receiveNotifications failed: %v", err)
	}

	if !strings.Contains(buf.String(), "Got finding") {
		t.Errorf("Did not receive any notifications.")
	}
}
