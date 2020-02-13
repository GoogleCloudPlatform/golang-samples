package notifications

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func orgID() string {
	return os.Getenv("GCLOUD_ORGANIZATION")
}

func projectID() string {
	return os.Getenv("GCLOUD_PROJECT")
}

func pubsubTopic() string {
	return os.Getenv("GCLOUD_PUBSUB_TOPIC")
}

func pubsubSubscription() string {
	return os.Getenv("GCLOUD_PUBSUB_SUBSCRIPTION")
}

func addNotificationConfig(notificationConfigID string) error {
	orgID := orgID()
	pubsubTopic := pubsubTopic()

	ctx := context.Background()
	client, err0 := securitycenter.NewClient(ctx)

	if err0 != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err0)
	}
	defer client.Close()

	req := &securitycenterpb.CreateNotificationConfigRequest{
		Parent:   fmt.Sprintf("organizations/%s", orgID),
		ConfigId: notificationConfigID,
		NotificationConfig: &securitycenterpb.NotificationConfig{
			Description: "Go sample config",
			PubsubTopic: pubsubTopic,
			EventType:   securitycenterpb.NotificationConfig_FINDING,
			NotifyConfig: &securitycenterpb.NotificationConfig_StreamingConfig_{
				StreamingConfig: &securitycenterpb.NotificationConfig_StreamingConfig{
					Filter: `state = "ACTIVE"`,
				},
			},
		},
	}

	_, err1 := client.CreateNotificationConfig(ctx, req)
	if err1 != nil {
		return fmt.Errorf("Failed to create notification config: %v", err1)
	}

	return nil
}

func cleanupNotificationConfig(notificationConfigID string) error {
	orgID := orgID()

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
	configID := "go-test-create-config-id"
	
	if err := createNotificationConfig(buf, orgID(), pubsubTopic(), configID); err != nil {
		t.Fatalf("createNotificationConfig failed: %v", err)
	}

	if !strings.Contains(buf.String(), "New NotificationConfig created") {
		t.Errorf("createNotificationConfig did not create.")
	}

	cleanupNotificationConfig(configID)
}

func TestDeleteNotificationConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	configID := "go-test-delete-config-id"
	
	if err0 := addNotificationConfig(configID); err0 != nil {
		t.Fatalf("Could not setup test environment: %v", err0)
	}
	
	if err1 := deleteNotificationConfig(buf, orgID(), configID); err1 != nil {
		t.Fatalf("deleteNotificationConfig failed: %v", err1)
	}

	if !strings.Contains(buf.String(), "Deleted config:") {
		t.Errorf("deleteNotificationConfig did not delete.")
	}
}

func TestGetNotificationConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	configID := "go-test-get-config-id"

	if err0 := addNotificationConfig(configID); err0 != nil {
		t.Fatalf("Could not setup test environment: %v", err0)
	}
	
	if err1 := getNotificationConfig(buf, orgID(), configID); err1 != nil {
		t.Fatalf("getNotificationConfig failed: %v", err1)
	}

	if !strings.Contains(buf.String(), "Deleted config:") {
		t.Errorf("getNotificationConfig did not delete.")
	}

	cleanupNotificationConfig(configID)
}

func TestListNotificationConfigs(t *testing.T) {
	buf := new(bytes.Buffer)
	configID := "go-test-list-config-id"

	if err0 := addNotificationConfig(configID); err0 != nil {
		t.Fatalf("Could not setup test environment: %v", err0)
	}

	if err1 := listNotificationConfigs(buf, orgID()); err1 != nil {
		t.Fatalf("listNotificationConfig failed: %v", err1)
	}

	if !strings.Contains(buf.String(), "NotificationConfig") {
		t.Errorf("listNotificationConfigs did not list")
	}

	cleanupNotificationConfig(configID)
}

func TestUpdateNotificationConfigs(t *testing.T) {
	buf := new(bytes.Buffer)
	configID := "go-test-update-config-id"

	if err0 := addNotificationConfig(configID); err0 != nil {
		t.Fatalf("Could not setup test environment: %v", err0)
	}

	if err1 := updateNotificationConfig(buf, orgID(), pubsubTopic(), configID); err1 != nil {
		t.Fatalf("updateNotificationConfig failed: %v", err1)
	}

	if !strings.Contains(buf.String(), "Updated NotificationConfig:") {
		t.Errorf("updateNotificationConfig did not update.")
	}

	cleanupNotificationConfig(configID)
}

func TestReceiveNotifications(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := receiveMessages(buf, projectID(), pubsubSubscription()); err != nil {
		t.Fatalf("receiveNotifications failed: %v", err)
	}

	if !strings.Contains(buf.String(), "Got finding") {
		t.Errorf("Did not receive any notifications.")
	}
}
