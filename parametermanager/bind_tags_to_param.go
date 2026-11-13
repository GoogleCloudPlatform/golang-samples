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

// [START parametermanager_bind_tags_to_param]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
)

// bindTagsToParam creates a parameter and binds a tag value to it.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter.
// tagValue: The tag value, in the form tagValues/{id}.
//
// The function returns an error if the operation fails.
func bindTagsToParam(w io.Writer, projectID, parameterID, tagValue string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parent resource.
	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Create the parameter without tags.
	parameter, err := client.CreateParameter(ctx, &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: parametermanagerpb.ParameterFormat_UNFORMATTED,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create parameter: %w", err)
	}
	fmt.Fprintf(w, "Created parameter: %s\n", parameter.Name)

	// Create a Resource Manager tag bindings client.
	tagBindingsClient, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create tag bindings client: %w", err)
	}
	defer tagBindingsClient.Close()

	// Bind the tag value to the parameter.
	op, err := tagBindingsClient.CreateTagBinding(ctx, &resourcemanagerpb.CreateTagBindingRequest{
		TagBinding: &resourcemanagerpb.TagBinding{
			Parent:   fmt.Sprintf("//parametermanager.googleapis.com/%s", parameter.Name),
			TagValue: tagValue,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create tag binding: %w", err)
	}

	// Wait for the operation to complete.
	binding, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for tag binding: %w", err)
	}

	fmt.Fprintf(w, "Bound tag value %s to parameter: %s\n", binding.TagValue, parameter.Name)
	return nil
}

// [END parametermanager_bind_tags_to_param]
