// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     [https://www.apache.org/licenses/LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample code for updating the labels of the given model armor template.

package modelarmor

// [START modelarmor_update_template_with_labels]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateModelArmorTemplateLabels method updates
// the labels of a Model Armor template.
//
// w io.Writer: The writer to use for logging.
// projectID string: The ID of the project.
// locationID string: The ID of the location.
// templateID string: The ID of the template.
// labels map[string]string: The updated labels.
func updateModelArmorTemplateLabels(w io.Writer, projectID, locationID, templateID string, labels map[string]string) error {
	ctx := context.Background()

	// Create options for Model Armor client.
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID))
	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create client for project %s, location %s: %w", projectID, locationID, err)
	}
	defer client.Close()

	// Build the Model Armor template with your preferred filters.
	// For more details on filters, please refer to the following doc:
	// [https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters](https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters)
	template := &modelarmorpb.Template{
		Name:   fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		Labels: labels,
	}

	// Prepare the request to update the template.
	updateMask := &fieldmaskpb.FieldMask{
		Paths: []string{"labels"},
	}

	req := &modelarmorpb.UpdateTemplateRequest{
		Template:   template,
		UpdateMask: updateMask,
	}

	// Update the template.
	response, err := client.UpdateTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	fmt.Fprintf(w, "Updated Model Armor Template Labels: %s\n", response.Name)

	return nil
}

// [END modelarmor_update_template_with_labels]
