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

// Sample code for updating the model armor project floor settings.

package modelarmor

// [START modelarmor_update_project_floor_settings]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
)

// updateProjectFloorSettings updates the floor settings of a project.
//
// This method updates the floor settings of a project.
//
// w io.Writer: The writer to use for logging.
// projectID string: The ID of the project.
func updateProjectFloorSettings(w io.Writer, projectID string) error {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Prepare project floor setting path/name
	floorSettingsName := fmt.Sprintf("projects/%s/locations/global/floorSetting", projectID)

	// Update the project floor setting
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
		return fmt.Errorf("failed to update floor setting: %w", err)
	}

	// Print the updated config
	fmt.Fprintf(w, "Updated project floor setting: %+v\n", response)

	return nil
}

// [END modelarmor_update_project_floor_settings]
