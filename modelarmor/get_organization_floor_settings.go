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

// Sample code for getting floor settings of an organization.

package modelarmor

// [START modelarmor_get_organization_floor_settings]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
)

// getOrganizationFloorSettings gets details of a single floor setting of an organization.
//
// This method retrieves the details of a single floor setting of an organization.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	organizationID string: The ID of the organization.
//
// Returns:
//
//	*modelarmorpb.FloorSetting: The retrieved floor setting.
//	error: Any error that occurred during retrieval.
//
// Example:
//
//	floorSetting, err := getOrganizationFloorSettings(
//	    os.Stdout,
//	    "my-organization",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(floorSetting)
//

func getOrganizationFloorSettings(w io.Writer, organizationID string) (*modelarmorpb.FloorSetting, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	floorSettingsName := fmt.Sprintf("organizations/%s/locations/global/floorSetting", organizationID)

	// Get the organization floor setting.
	req := &modelarmorpb.GetFloorSettingRequest{
		Name: floorSettingsName,
	}

	response, err := client.GetFloorSetting(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get floor setting: %v", err)
	}

	// Print the retrieved floor setting using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Retrieved org floor setting: %v\n", response)

	// [END modelarmor_get_organization_floor_settings]

	return response, nil
}
