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

package regional_parametermanager

// [START parametermanager_update_regional_param_template_labels]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateRegionalParamTemplateLabels adds or updates a label on a parameter template.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The region where the resources are located.
// templateID: The ID of the template.
//
// The function returns an error if the operation fails.
func updateRegionalParamTemplateLabels(w io.Writer, projectID, locationID, templateID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the template.
	name := fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID)

	// Build the request to update the labels of the template.
	req := &parametermanagerpb.UpdateTemplateRequest{
		Template: &parametermanagerpb.Template{
			Name:   name,
			Labels: map[string]string{"environment": "test"},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"labels"}},
	}

	// Call the API to update the template.
	if _, err := client.UpdateTemplate(ctx, req); err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	// Get the template to read the labels.
	template, err := client.GetTemplate(ctx, &parametermanagerpb.GetTemplateRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	fmt.Fprintf(w, "Updated regional parameter template %s with labels %v\n", template.Name, template.Labels)
	return nil
}

// [END parametermanager_update_regional_param_template_labels]
