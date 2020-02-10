package notifications

// [START scc_update_notification_config]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
	"google.golang.org/genproto/protobuf/field_mask"
)

func updateNotificationConfig(orgId string, notificationConfigId string, updatedPubsubTopic string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	// TODO orgId := "1081635000895"
	// TODO: notificationConfigId := "go-sample-config-id"
	// TODO: updatedPubsubTopic := "projects/{new-project}/topics/{new-topic}"

	updatedDescription := "Updated sample config"
	req := &securitycenterpb.UpdateNotificationConfigRequest{
		NotificationConfig: &securitycenterpb.NotificationConfig{
			Name:        fmt.Sprintf("organizations/%s/notificationConfigs/%s", orgId, notificationConfigId),
			Description: updatedDescription,
			PubsubTopic: updatedPubsubTopic,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"description", "pubsub_topic"},
		},
	}

	notificationConfig, err := client.UpdateNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to create notification config: %v", err)
	}

	fmt.Printf("Updated NotificationConfig: %s\n", notificationConfig)

	return nil
}

// [END scc_update_notification_config]
