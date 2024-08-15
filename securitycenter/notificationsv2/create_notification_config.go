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

// [START securitycenter_create_notification_config_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

func createNotificationConfig(w io.Writer, orgID string, pubsubTopic string, notificationConfigID string) error {
	// orgID := "your-org-id"
	// pubsubTopic := "projects/{your-project}/topics/{your-topic}"
	// notificationConfigID := "your-config-id"

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.CreateNotificationConfigRequest{
		// Parent must be in one of the following formats:
		//		"organizations/{orgId}/locations/global"
		//		"projects/{projectId}/locations/global"
		//		"folders/{folderId}/locations/global"
		Parent:   fmt.Sprintf("organizations/%s/locations/global", orgID),
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

	notificationConfig, err := client.CreateNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to create notification config: %w", err)
	}
	fmt.Fprintln(w, "New NotificationConfig created: ", notificationConfig)

	return nil
}

// [END securitycenter_create_notification_config_v2]
