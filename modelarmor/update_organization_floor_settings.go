// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     [https://www.apache.org/licenses/LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample code for updating the model armor floor settings of an organization.

package modelarmor

// [START modelarmor_update_organization_floor_settings]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// updateOrganizationFloorSettings updates floor settings of an organization.
//
// This method updates the floor settings of an organization.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	organizationID string: The ID of the organization.
//	locationID string: The ID of the location.
//
// Returns:
//
//	*modelarmorpb.FloorSetting: The updated floor settings.
//	error: Any error that occurred during update.
//
// Example:
//
//	updatedSettings, err := updateOrganizationFloorSettings(
//	    os.Stdout,
//	    "my-organization",
//	    "my-location",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(updatedSettings)
func updateOrganizationFloorSettings(w io.Writer, organizationID, locationID string) (*modelarmorpb.FloorSetting, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Prepare organization floor setting path/name
	floorSettingsName := fmt.Sprintf("organizations/%s/locations/global/floorSetting", organizationID)

	// Update the organization floor setting
	// For more details on filters, please refer to the following doc:
	// [https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters](https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters)
	enableEnforcement := true
	req := &modelarmorpb.UpdateFloorSettingRequest{
		FloorSetting: &modelarmorpb.FloorSetting{
			Name: floorSettingsName,
			FilterConfig: &modelarmorpb.FilterConfig{
				RaiSettings: &modelarmorpb.RaiFilterSettings{
					RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
						{
							FilterType:      modelarmorpb.RaiFilterType_HATE_SPEECH,
							ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
						},
					},
				},
			},
			EnableFloorSettingEnforcement: &enableEnforcement,
		},
	}

	response, err := client.UpdateFloorSetting(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update floor setting: %v", err)
	}

	// Print the updated config
	fmt.Fprintf(w, "Updated org floor setting: %+v\n", response)

	// [END modelarmor_update_organization_floor_settings]

	return response, nil
}
