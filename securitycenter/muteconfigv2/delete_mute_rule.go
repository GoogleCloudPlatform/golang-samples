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

package muteconfigv2

// [START securitycenter_delete_mute_config_v2]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

// deleteMuteRule deletes a mute configuration given its resource name.
// Note: Previously muted findings are not affected when a mute config is deleted.
func deleteMuteRule(w io.Writer, parent string, muteConfigId string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}
	//             - folders/{folder_id}
	//             - projects/{project_id}
	// parent := fmt.Sprintf("projects/%s", "your-google-cloud-project-id")
	//
	// muteConfigId: Specify the name of the mute config to delete.
	// muteConfigId := "mute-config-id"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.DeleteMuteConfigRequest{
		Name: fmt.Sprintf("%s/muteConfigs/%s", parent, muteConfigId),
	}

	if err := client.DeleteMuteConfig(ctx, req); err != nil {
		return fmt.Errorf("failed to delete Muteconfig: %w", err)
	}
	fmt.Fprintf(w, "Mute rule deleted successfully: %s", muteConfigId)
	return nil
}

// [END securitycenter_delete_mute_config_v2]
