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

// [START parametermanager_enable_param_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// enableParamVersion enables a parameter version.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter for which the version is to be enabled.
// versionID: The ID of the version to be enabled.
//
// The function returns an error if the parameter version update fails.
func enableParamVersion(w io.Writer, projectID, parameterID, versionID string) error {
	// Create the client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create parametermanager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parameter version to enable.
	name := fmt.Sprintf("projects/%s/locations/global/parameters/%s/versions/%s", projectID, parameterID, versionID)

	// Build the request to enable the parameter version by updating the parameter version.
	req := &parametermanagerpb.UpdateParameterVersionRequest{
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"disabled"},
		},
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Name:     name,
			Disabled: false,
		},
	}

	// Call the API to enable the parameter version.
	if _, err := client.UpdateParameterVersion(ctx, req); err != nil {
		return fmt.Errorf("failed to enable parameter version: %w", err)
	}

	fmt.Fprintf(w, "Enabled parameter version %s for parameter %s\n", name, parameterID)
	return nil
}

// [END parametermanager_enable_param_version]
