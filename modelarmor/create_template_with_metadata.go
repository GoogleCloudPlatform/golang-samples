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

// Sample code for creating a new model armor template with template metadata.

package modelarmor

// [START modelarmor_create_template_with_metadata]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// createModelArmorTemplateWithMetadata method creates a
// new Model Armor template with template metadata.
//
// w io.Writer: The writer to use for logging.
// projectID string: The ID of the Google Cloud project.
// locationID string: The ID of the Google Cloud location.
// templateID string: The ID of the template to create.
// metadata *modelarmorpb.TemplateMetadata: The template metadata to apply.
func createModelArmorTemplateWithMetadata(w io.Writer, projectID, locationID, templateID string) error {
	ctx := context.Background()

	// Create options for Model Armor client.
	endpoint := fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)
	opts := option.WithEndpoint(endpoint)
	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,opts)
	if err != nil {
		return fmt.Errorf("failed to create client for project %s, location %s: %w", projectID, locationID, err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Build the Model Armor template with your preferred filters.
	// For more details on filters, please refer to the following doc:
	// [https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters](https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters)
	template := &modelarmorpb.Template{
		FilterConfig: &modelarmorpb.FilterConfig{
			RaiSettings: &modelarmorpb.RaiFilterSettings{
				RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
					{
						FilterType:      modelarmorpb.RaiFilterType_HATE_SPEECH,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_SEXUALLY_EXPLICIT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
					},
				},
			},
		},
		// Add template metadata to the template.
		// For more details on template metadata, please refer to the following doc:
		// [https://cloud.google.com/security-command-center/docs/reference/model-armor/rest/v1/projects.locations.templates#templatemetadata](https://cloud.google.com/security-command-center/docs/reference/model-armor/rest/v1/projects.locations.templates#templatemetadata)
		TemplateMetadata: &modelarmorpb.Template_TemplateMetadata{
			IgnorePartialInvocationFailures: true,
			LogSanitizeOperations:           true,
		},
	}

	// Prepare the request for creating the template.
	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     parent,
		TemplateId: templateID,
		Template:   template,
	}

	// Create the template.
	response, err := client.CreateTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	// Print the new template name using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Created Model Armor Template: %s\n", response.Name)

	return err
}

// [END modelarmor_create_template_with_metadata]
