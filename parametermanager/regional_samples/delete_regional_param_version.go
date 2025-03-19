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

package regional_parametermanager

// [START parametermanager_delete_regional_param_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

// deleteRegionalParamVersion deletes a regional parameter version using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The ID of the region where the parameter is located.
// parameterID: The ID of the parameter for which the version is to be deleted.
// versionID: The ID of the version to be deleted.
//
// The function returns an error if the parameter version deletion fails.
func deleteRegionalParamVersion(w io.Writer, projectID, locationID, parameterID, versionID string) error {
	// Create a new context.
	ctx := context.Background()

	// Create a Parameter Manager client.
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("Failed to create Parameter Manager client: %v\n", err)
	}
	defer client.Close()

	// Construct the name of the parameter version to delete.
	name := fmt.Sprintf("projects/%s/locations/%s/parameters/%s/versions/%s", projectID, locationID, parameterID, versionID)

	// Build the request to delete the parameter version.
	req := &parametermanagerpb.DeleteParameterVersionRequest{
		Name: name,
	}

	// Call the API to delete the parameter version.
	if err := client.DeleteParameterVersion(ctx, req); err != nil {
		return fmt.Errorf("Failed to delete parameter version: %v\n", err)
	}

	// Output a success message.
	fmt.Fprintf(w, "Deleted regional parameter version: %s\n", name)
	return nil
}

// [END parametermanager_delete_regional_param_version]
