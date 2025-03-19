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

// [START parametermanager_get_param]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// getParam get parameter using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter to retrieved.
//
// The function returns an error if the parameter retrieval fails.
func getParam(w io.Writer, projectID, parameterID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create Parameter Manager client: %v\n", err)
	}
	defer client.Close()

	// Construct the name of the parameter to get parameter.
	name := fmt.Sprintf("projects/%s/locations/global/parameters/%s", projectID, parameterID)

	// Build the request to get parameter.
	req := &parametermanagerpb.GetParameterRequest{
		Name: name,
	}

	// Call the API to get parameter.
	param, err := client.GetParameter(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to get parameter: %v\n", err)
	}

	// Find more details for the Parameter object here:
	// https://cloud.google.com/secret-manager/parameter-manager/docs/reference/rest/v1/projects.locations.parameters#Parameter
	fmt.Fprintf(w, "Found parameter %s with format %s\n", param.Name, param.Format.String())
	return nil
}

// [END parametermanager_get_param]
