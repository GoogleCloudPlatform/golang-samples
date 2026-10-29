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

// [START parametermanager_delete_param_template]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// deleteParamTemplate deletes a parameter template.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// templateID: The ID of the template.
//
// The function returns an error if the operation fails.
func deleteParamTemplate(w io.Writer, projectID, templateID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the template.
	name := fmt.Sprintf("projects/%s/locations/global/templates/%s", projectID, templateID)

	// Build the request to delete the template.
	req := &parametermanagerpb.DeleteTemplateRequest{
		Name: name,
	}

	// Call the API to delete the template.
	if err := client.DeleteTemplate(ctx, req); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	fmt.Fprintf(w, "Deleted parameter template: %s\n", name)
	return nil
}

// [END parametermanager_delete_param_template]
