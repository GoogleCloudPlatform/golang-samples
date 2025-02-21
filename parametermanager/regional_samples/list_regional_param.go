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

// [START parametermanager_list_regional_param]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listRegionalParam lists all parameters regional using the Parameter Manager SDK for GCP.
//
// projectID: The ID of the project where the parameter is located.
// locationID: The ID of the region where the parameter is located.
// parameterID: The ID of the parameter to be listed.
//
// The function returns an error if the parameter listing fails
func listRegionalParam(w io.Writer, projectID, locationID string) error {
	// Create a new context.
	ctx := context.Background()

	// Create a Parameter Manager client.
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("Failed to create Parameter Manager client: %v\n", err)
	}
	defer client.Close()

	// Construct the name of the parent resource to list parameters.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Build the request to list all parameters.
	req := &parametermanagerpb.ListParametersRequest{
		Parent: parent,
	}

	// Call the API to list all parameters.
	it := client.ListParameters(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Failed to list parameters: %v\n", err)
		}

		// Print the name of the parameter.
		fmt.Fprintf(w, "Found regional parameter: %s with format %s \n", resp.Name, resp.Format)
	}

	return nil
}

// [END parametermanager_list_regional_param]
