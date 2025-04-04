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

// [START parametermanager_create_regional_param_version_with_secret]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

// createRegionalParamVersionWithSecret creates a new version of a regional parameter with a JSON payload that has a secret reference in it using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The ID of the region where the parameter is located.
// parameterID: The ID of the parameter for which the version is to be created.
// versionID: The ID of the version to be created.
// secretID: The ID of the secret to be referenced.
//
// The function returns an error if the parameter version creation fails.
func createRegionalParamVersionWithSecret(w io.Writer, projectID, locationID, parameterID, versionID, secretID string) error {
	// Create a context.
	ctx := context.Background()

	// Create a Parameter Manager client.
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create parametermanager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the create parameter version.
	parent := fmt.Sprintf("projects/%s/locations/%s/parameters/%s", projectID, locationID, parameterID)

	payload := []byte(fmt.Sprintf(`{"db_host": "localhost","db_password": "__REF__(//secretmanager.googleapis.com/%s)"}`, secretID))

	// Create a parameter version.
	req := &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parent,
		ParameterVersionId: versionID,
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data: payload,
			},
		},
	}
	version, err := client.CreateParameterVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create parameter version: %w", err)
	}
	fmt.Fprintf(w, "Created regional parameter version with secret reference: %s\n", version.Name)
	return nil
}

// [END parametermanager_create_regional_param_version_with_secret]
