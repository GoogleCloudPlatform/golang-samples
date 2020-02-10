package notifications

// [START scc_delete_notification_config]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func deleteNotificationConfig(orgId string, notificationConfigId string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	// TODO orgId := "your-org-id"
	// TODO notificationConfigId := "config-to-delete"

	name := fmt.Sprintf("organizations/%s/notificationConfigs/%s", orgId, notificationConfigId)
	req := &securitycenterpb.DeleteNotificationConfigRequest{
		Name: name,
	}

	err = client.DeleteNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to retrieve notification config: %v", err)
	}
	fmt.Printf("Deleted config: %s\n", name)

	return nil
}

// [END scc_delete_notification_config]
