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

// [START securitycenter_get_mute_config_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

// getMuteRule retrieves a mute configuration given its resource name.
func getMuteRule(w io.Writer, parent string, muteConfigId string) error {
	// Use any one of the following resource paths to get mute configuration:
	//         - organizations/{organization_id}
	//         - folders/{folder_id}
	//         - projects/{project_id}
	// parent := fmt.Sprintf("projects/%s", "your-google-cloud-project-id")
	//
	// Name of the mute config to retrieve.
	// muteConfigId := "mute-config-id"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.GetMuteConfigRequest{
		Name: fmt.Sprintf("%s/muteConfigs/%s", parent, muteConfigId),
	}

	muteconfig, err := client.GetMuteConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to retrieve Muteconfig: %w", err)
	}
	fmt.Fprintf(w, "Muteconfig Name: %s ", muteconfig.Name)
	return nil
}

// [END securitycenter_get_mute_config_v2]
