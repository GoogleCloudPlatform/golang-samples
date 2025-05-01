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

// Sample code for getting floor settings of a folder.

package modelarmor

// [START modelarmor_get_folder_floor_settings]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
)

// getFolderFloorSettings gets details of a single floor setting of a folder.
//
// This method retrieves the details of a single floor setting of a folder.
//
// w io.Writer: The writer to use for logging.
// folderID string: The ID of the folder.
func getFolderFloorSettings(w io.Writer, folderID string) error {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Prepare folder floor setting path/name
	floorSettingsName := fmt.Sprintf("folders/%s/locations/global/floorSetting", folderID)

	// Get the folder floor setting.
	req := &modelarmorpb.GetFloorSettingRequest{
		Name: floorSettingsName,
	}

	response, err := client.GetFloorSetting(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get floor setting: %w", err)
	}

	// Print the retrieved floor setting using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Retrieved folder floor setting: %v\n", response)

	return nil
}

// [END modelarmor_get_folder_floor_settings]
