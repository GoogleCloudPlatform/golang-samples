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

package parametermanager

// [START parametermanager_get_param_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// getParamVersion get parameter version using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter for which the version details are to be retrieved.
// versionID: The ID of the version to be retrieved.
//
// The function returns an error if the parameter version retrieval fails.
func getParamVersion(w io.Writer, projectID, parameterID, versionID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parameter to get the parameter version.
	name := fmt.Sprintf("projects/%s/locations/global/parameters/%s/versions/%s", projectID, parameterID, versionID)

	// Build the request to get parameter version.
	req := &parametermanagerpb.GetParameterVersionRequest{
		Name: name,
	}

	// Call the API to get parameter version.
	version, err := client.GetParameterVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get parameter version: %w", err)
	}

	// Find more details for the Parameter Version object here:
	// https://cloud.google.com/secret-manager/parameter-manager/docs/reference/rest/v1/projects.locations.parameters.versions#ParameterVersion
	fmt.Fprintf(w, "Found parameter version %s with disabled state in %v\n", version.Name, version.Disabled)
	if !version.Disabled {
		fmt.Fprintf(w, "Payload: %s\n", version.Payload.Data)
	}
	return nil
}

// [END parametermanager_get_param_version]
