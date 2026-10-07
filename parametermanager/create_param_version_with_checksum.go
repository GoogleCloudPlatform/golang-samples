// Copyright 2026 Google LLC
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

// [START parametermanager_create_param_version_with_checksum]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"hash/crc32"
)

// createParamVersionWithChecksum creates a parameter version with a client-computed CRC32C checksum. Parameter Manager verifies the checksum and rejects the request if it does not match the payload.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter.
// versionID: The ID of the template version.
// payload: The template payload containing {{.variableName}} placeholders.
//
// The function returns an error if the operation fails.
func createParamVersionWithChecksum(w io.Writer, projectID, parameterID, versionID, payload string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parent parameter.
	parent := fmt.Sprintf("projects/%s/locations/global/parameters/%s", projectID, parameterID)

	// Compute the CRC32C checksum (Castagnoli polynomial) of the payload.
	data := []byte(payload)
	crc32c := int64(crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli)))

	// Build the request to create a parameter version with the checksum.
	req := &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parent,
		ParameterVersionId: versionID,
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data:       data,
				DataCrc32C: &crc32c,
			},
		},
	}

	// Call the API to create the parameter version.
	version, err := client.CreateParameterVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create parameter version: %w", err)
	}

	fmt.Fprintf(w, "Created parameter version: %s with checksum source %s\n", version.Name, version.GetChecksumSource())
	return nil
}

// [END parametermanager_create_param_version_with_checksum]
