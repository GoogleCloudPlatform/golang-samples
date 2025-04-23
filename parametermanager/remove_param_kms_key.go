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

// [START parametermanager_remove_param_kms_key]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// removeParamKmsKey removes a parameter kms_key using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter to be updated.
//
// The function returns an error if the parameter creation fails.
func removeParamKmsKey(w io.Writer, projectID, parameterID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()

	// Create a Parameter Manager client.
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the create parameter.
	name := fmt.Sprintf("projects/%s/locations/global/parameters/%s", projectID, parameterID)

	// Create a parameter with unformatted format.
	req := &parametermanagerpb.UpdateParameterRequest{
		Parameter: &parametermanagerpb.Parameter{
			Name:   name,
			Format: parametermanagerpb.ParameterFormat_UNFORMATTED,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"kms_key"},
		},
	}
	parameter, err := client.UpdateParameter(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update parameter: %w", err)
	}

	fmt.Fprintf(w, "Removed kms_key for parameter %s\n", parameter.Name)
	return nil
}

// [END parametermanager_remove_param_kms_key]
