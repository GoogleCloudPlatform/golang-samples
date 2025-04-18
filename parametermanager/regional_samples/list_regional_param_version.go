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

// [START parametermanager_list_regional_param_versions]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listRegionalParamVersion lists all parameter versions regional using the Parameter Manager SDK for GCP.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The ID of the region where the parameter is located.
// parameterID: The ID of the parameter for which the version is to be listed.
//
// The function returns an error if the parameter version listing fails
func listRegionalParamVersion(w io.Writer, projectID, locationID, parameterID string) error {
	// Create a new context.
	ctx := context.Background()

	// Create a Parameter Manager client.
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parameter to list versions.
	parent := fmt.Sprintf("projects/%s/locations/%s/parameters/%s", projectID, locationID, parameterID)

	// Build the request to list all parameter versions.
	req := &parametermanagerpb.ListParameterVersionsRequest{
		Parent: parent,
	}

	// Call the API to list all parameter versions.
	parameterVersions := client.ListParameterVersions(ctx, req)
	for {
		version, err := parameterVersions.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list parameter versions: %w", err)
		}

		fmt.Fprintf(w, "Found regional parameter version %s with disabled state in %v\n", version.Name, version.Disabled)
	}

	return nil
}

// [END parametermanager_list_regional_param_versions]
