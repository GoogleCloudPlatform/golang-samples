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

// [START parametermanager_get_param_version_verify_checksum]
import (
	"context"
	"fmt"
	"io"

	"hash/crc32"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// getParamVersionVerifyChecksum retrieves a parameter version and verifies its checksum.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter.
// versionID: The ID of the template version.
//
// The function returns an error if the operation fails.
func getParamVersionVerifyChecksum(w io.Writer, projectID, parameterID, versionID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parameter version.
	name := fmt.Sprintf("projects/%s/locations/global/parameters/%s/versions/%s", projectID, parameterID, versionID)

	// Build the request to get the parameter version.
	req := &parametermanagerpb.GetParameterVersionRequest{
		Name: name,
		View: parametermanagerpb.View_FULL,
	}

	// Call the API to get the parameter version.
	version, err := client.GetParameterVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get parameter version: %w", err)
	}

	// Verify the checksum.
	data := version.Payload.Data
	computed := int64(crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli)))
	if version.Payload.DataCrc32C == nil || *version.Payload.DataCrc32C != computed {
		return fmt.Errorf("checksum mismatch for parameter version %s", version.Name)
	}

	fmt.Fprintf(w, "Verified checksum of parameter version %s (checksum source: %s)\n", version.Name, version.GetChecksumSource())
	return nil
}

// [END parametermanager_get_param_version_verify_checksum]
