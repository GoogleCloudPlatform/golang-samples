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

// [START parametermanager_create_param_template]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// createParamTemplate creates a new parameter template in Parameter Manager.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// templateID: The ID of the template.
// format: The format of the template (YAML or JSON).
//
// The function returns an error if the operation fails.
func createParamTemplate(w io.Writer, projectID, templateID string, format parametermanagerpb.TemplateFormat) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the parent resource.
	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Build the request to create a new template.
	req := &parametermanagerpb.CreateTemplateRequest{
		Parent:     parent,
		TemplateId: templateID,
		Template: &parametermanagerpb.Template{
			Format: format,
		},
	}

	// Call the API to create the template.
	template, err := client.CreateTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	fmt.Fprintf(w, "Created parameter template: %s with format %s\n", template.Name, template.Format)
	return nil
}

// [END parametermanager_create_param_template]
