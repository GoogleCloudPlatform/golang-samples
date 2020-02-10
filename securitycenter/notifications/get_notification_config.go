package notifications

// [START scc_get_notification_config]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func getNotificationConfig(orgId string, notificationConfigId string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	// TODO orgId := "your-org-id"
	// TODO notificationConfigId := "your-config-id"

	req := &securitycenterpb.GetNotificationConfigRequest{
		Name: fmt.Sprintf("organizations/%s/notificationConfigs/%s", orgId, notificationConfigId),
	}

	notificationConfig, err := client.GetNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to retrieve notification config: %v", err)
	}
	fmt.Printf("Received config: %s\n", notificationConfig)

	return nil
}

// [END scc_get_notification_config]
