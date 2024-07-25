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

// [START securitycenter_get_notification_config_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

func getNotificationConfig(w io.Writer, orgID string, notificationConfigID string) error {
	// orgID := "your-org-id"
	// notificationConfigID := "your-config-id"

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	// Parent must be in one of the following formats:
	//		"organizations/{orgId}/locations/global"
	//		"projects/{projectId}/locations/global"
	//		"folders/{folderId}/locations/global"
	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	req := &securitycenterpb.GetNotificationConfigRequest{
		Name: fmt.Sprintf("%s/notificationConfigs/%s", parent, notificationConfigID),
	}

	notificationConfig, err := client.GetNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to retrieve notification config: %w", err)
	}
	fmt.Fprintln(w, "Received config: ", notificationConfig)

	return nil
}

// [END securitycenter_get_notification_config_v2]
