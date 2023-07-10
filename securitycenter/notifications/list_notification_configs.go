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

// [START securitycenter_list_notification_configs]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"google.golang.org/api/iterator"
)

func listNotificationConfigs(w io.Writer, orgID string) error {
	// orgId := "your-org-id"

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.ListNotificationConfigsRequest{
		// Parent must be in one of the following formats:
		//		"organizations/{orgId}"
		//		"projects/{projectId}"
		//		"folders/{folderId}"
		Parent: fmt.Sprintf("organizations/%s", orgID),
	}
	it := client.ListNotificationConfigs(ctx, req)
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("it.Next: %w", err)
		}

		fmt.Fprintln(w, "NotificationConfig: ", result)
	}

	return nil
}

// [END securitycenter_list_notification_configs]
