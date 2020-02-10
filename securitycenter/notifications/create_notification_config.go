package notifications

// [START scc_create_notification_config]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func createNotificationConfig(orgId string, pubsubTopic string, notificationConfigId string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	// TODO orgId := "your-org-id"
	// TODO pubsubTopic := "projects/{your-project}/topics/{your-topic}"
	// TODO notificationConfigId := "your-config-id"

	req := &securitycenterpb.CreateNotificationConfigRequest{
		Parent:   fmt.Sprintf("organizations/%s", orgId),
		ConfigId: notificationConfigId,
		NotificationConfig: &securitycenterpb.NotificationConfig{
			Description: "Go sample config",
			PubsubTopic: pubsubTopic,
			EventType:   securitycenterpb.NotificationConfig_FINDING,
			NotifyConfig: &securitycenterpb.NotificationConfig_StreamingConfig_{
				StreamingConfig: &securitycenterpb.NotificationConfig_StreamingConfig{
					Filter: "state = \"ACTIVE\"",
				},
			},
		},
	}

	notificationConfig, err := client.CreateNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to create notification config: %v", err)
	}
	fmt.Printf("New NotificationConfig created: %s\n", notificationConfig)

	return nil
}

// [END scc_create_notification_config]
