// Copyright 2023 Google LLC
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

package muteconfig

// [START securitycenter_update_mute_config]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateMuteRule Updates an existing mute configuration.
// The following can be updated in a mute config: description and filter.
func updateMuteRule(w io.Writer, muteConfigName string) error {
	// Specify the name of the mute config to delete.
	// muteConfigName: Use any one of the following formats:
	//                 - organizations/{organization}/muteConfigs/{config_id}
	//                 - folders/{folder}/muteConfigs/{config_id}
	//                 - projects/{project}/muteConfigs/{config_id}
	// muteConfigName := fmt.Sprintf("projects/%s/muteConfigs/%s", "project-id", "mute-config")
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	updateMuteConfig := &securitycenterpb.MuteConfig{
		Name:        muteConfigName,
		Description: "Updated mute config description",
	}

	req := &securitycenterpb.UpdateMuteConfigRequest{
		MuteConfig: updateMuteConfig,
		// Set the update mask to specify which properties of the mute config should be
		// updated.
		// If empty, all mutable fields will be updated.
		// Make sure that the mask fields match the properties changed in 'updateMuteConfig'.
		// For more info on constructing update mask path, see the proto or:
		// https://cloud.google.com/security-command-center/docs/reference/rest/v1/folders.muteConfigs/patch?hl=en#query-parameters
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"description",
			},
		},
	}

	response, err := client.UpdateMuteConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("mute rule update failed! %w", err)
	}
	fmt.Fprintf(w, "Mute rule updated %s", response.Name)
	return nil
}

// [END securitycenter_update_mute_config]
