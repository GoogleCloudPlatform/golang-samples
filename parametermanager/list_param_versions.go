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

// [START parametermanager_list_param_versions]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/iterator"
)

// listParamVersions lists parameter versions using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter for which the versions are to be listed.
//
// The function returns an error if the parameter version listing fails.
func listParamVersions(w io.Writer, projectID, parameterID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the list parameter.
	parent := fmt.Sprintf("projects/%s/locations/global/parameters/%s", projectID, parameterID)

	// Build the request to list parameter versions.
	req := &parametermanagerpb.ListParameterVersionsRequest{
		Parent: parent,
	}

	// Call the API to list parameter versions.
	parameterVersions := client.ListParameterVersions(ctx, req)
	for {
		version, err := parameterVersions.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list parameter versions: %w", err)
		}

		fmt.Fprintf(w, "Found parameter version %s with disabled state in %v\n", version.Name, version.Disabled)
	}

	return nil
}

// [END parametermanager_list_param_versions]
