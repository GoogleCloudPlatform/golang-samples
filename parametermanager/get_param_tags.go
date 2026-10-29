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

// [START parametermanager_get_param_tags]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
)

// getParamTags lists the tag bindings of a parameter.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter.
//
// The function returns an error if the operation fails.
func getParamTags(w io.Writer, projectID, parameterID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Get the parameter.
	parameter, err := client.GetParameter(ctx, &parametermanagerpb.GetParameterRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/parameters/%s", projectID, parameterID),
	})
	if err != nil {
		return fmt.Errorf("failed to get parameter: %w", err)
	}

	// Create a Resource Manager tag bindings client.
	tagBindingsClient, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create tag bindings client: %w", err)
	}
	defer tagBindingsClient.Close()

	// List the tag bindings of the parameter.
	it := tagBindingsClient.ListTagBindings(ctx, &resourcemanagerpb.ListTagBindingsRequest{
		Parent: fmt.Sprintf("//parametermanager.googleapis.com/%s", parameter.Name),
	})
	for {
		binding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list tag bindings: %w", err)
		}
		fmt.Fprintf(w, "Found tag binding on parameter %s: %s (%s)\n", parameter.Name, binding.TagValue, binding.TagValueNamespacedName)
	}
	return nil
}

// [END parametermanager_get_param_tags]
