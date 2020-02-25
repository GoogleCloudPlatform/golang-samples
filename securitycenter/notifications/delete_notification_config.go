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

// [START scc_delete_notification_config]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1p1beta1"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func deleteNotificationConfig(w io.Writer, orgID string, notificationConfigID string) error {
	// orgID := "your-org-id"
	// notificationConfigID := "config-to-delete"

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
	fmt.Fprintln(w, "Deleted config: ", name)

	return nil
}

// [END scc_delete_notification_config]
