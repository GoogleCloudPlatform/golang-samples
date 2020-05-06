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

// [START scc_update_notification_config]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

func updateNotificationConfig(w io.Writer, orgID string, notificationConfigID string, updatedPubsubTopic string) error {
	// orgID := "your-org-id"
	// notificationConfigID := "your-config-id"
	// updatedPubsubTopic := "projects/{new-project}/topics/{new-topic}"

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close()

	updatedDescription := "Updated sample config"
	req := &securitycenterpb.UpdateNotificationConfigRequest{
		NotificationConfig: &securitycenterpb.NotificationConfig{
			Name:        fmt.Sprintf("organizations/%s/notificationConfigs/%s", orgID, notificationConfigID),
			Description: updatedDescription,
			PubsubTopic: updatedPubsubTopic,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"description", "pubsub_topic"},
		},
	}

	notificationConfig, err := client.UpdateNotificationConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to update notification config: %v", err)
	}

	fmt.Fprintln(w, "Updated NotificationConfig: ", notificationConfig)

	return nil
}

// [END scc_update_notification_config]
