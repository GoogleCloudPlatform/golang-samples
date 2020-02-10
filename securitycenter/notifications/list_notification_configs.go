package notifications

// [START scc_list_notification_configs]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	"google.golang.org/api/iterator"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func listNotificationConfigs(orgId string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	// TODO orgId := "your-org-id"

	req := &securitycenterpb.ListNotificationConfigsRequest{
		Parent: fmt.Sprintf("organizations/%s", orgId),
	}
	it := client.ListNotificationConfigs(ctx, req)
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("it.Next: %v", err)
		}

		fmt.Printf("NotificationConfig: %s, \n", result)
	}

	return nil
}

// [END scc_list_notification_configs]
