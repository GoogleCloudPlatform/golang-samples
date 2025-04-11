// Copyright 2025 Google LLC
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

// Sample code for getting floor settings of a project.

package modelarmor

// [START modelarmor_get_project_floor_settings]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
)

// getProjectFloorSettings gets details of a single floor setting of a project.
//
// This method retrieves the details of a single floor setting of a project.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	projectID string: The ID of the project.
//
// Returns:
//
//	*modelarmorpb.FloorSetting: The retrieved floor setting details.
//	error: Any error that occurred during retrieval.
//
// Example:
//
//	floorSetting, err := getProjectFloorSettings(
//	    os.Stdout,
//	    "my-project",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(floorSetting)
func getProjectFloorSettings(w io.Writer, projectID string) (*modelarmorpb.FloorSetting, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	floorSettingsName := fmt.Sprintf("projects/%s/locations/global/floorSetting", projectID)

	// Get the project floor setting.
	req := &modelarmorpb.GetFloorSettingRequest{
		Name: floorSettingsName,
	}

	response, err := client.GetFloorSetting(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get floor setting: %w", err)
	}

	// Print the retrieved floor setting using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Retrieved floor setting: %v\n", response)

	// [END modelarmor_get_project_floor_settings]

	return response, nil
}
