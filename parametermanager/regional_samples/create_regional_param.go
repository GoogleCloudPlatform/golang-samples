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

// [START parametermanager_create_regional_param]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

// createRegionalParam creates a parameter regional of the format type "unformatted" using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The region where the parameter is to be created.
// parameterID: The ID of the parameter to be created.
//
// The function returns an error if the parameter creation fails.
func createRegionalParam(w io.Writer, projectID, locationID, parameterID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()

	// Create a Parameter Manager client.
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("Failed to create Parameter Manager client: %v\n", err)
	}
	defer client.Close()

	// Construct the name of the create parameter.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Create a parameter with unformatted format.
	req := &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: parametermanagerpb.ParameterFormat_UNFORMATTED,
		},
	}
	parameter, err := client.CreateParameter(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to create parameter: %v\n", err)
	}

	// Print the created parameter.
	fmt.Fprintf(w, "Created regional parameter: %s\n", parameter.Name)
	return nil
}

// [END parametermanager_create_regional_param]
