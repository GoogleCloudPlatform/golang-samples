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

// Sample code for updating the model armor template metadata.

package modelarmor

// [START modelarmor_update_template_metadata]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// updateModelArmorTemplateMetadata updates the metadata of a Model Armor template.
//
// This method updates the metadata of a Model Armor template.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	projectID string: The ID of the project.
//	locationID string: The ID of the location.
//	templateID string: The ID of the template.
//
// Returns:
//
//	*modelarmorpb.Template: The updated template with new metadata.
//	error: Any error that occurred during update.
//
// Example:
//
//	updatedTemplate, err := updateModelArmorTemplateMetadata(
//	    os.Stdout,
//	    "my-project",
//	    "my-location",
//	    "my-template",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(updatedTemplate)
func updateModelArmorTemplateMetadata(w io.Writer, projectID, locationID, templateID string) (*modelarmorpb.Template, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for project %s, location %s: %v", projectID, locationID, err)
	}
	defer client.Close()

	// Build the full resource path for the template.
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID)

	// Build the Model Armor template with your preferred filters.
	// For more details on filters, please refer to the following doc:
	// [https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters](https://cloud.google.com/security-command-center/docs/key-concepts-model-armor#ma-filters)
	template := &modelarmorpb.Template{
		Name: templateName,
		FilterConfig: &modelarmorpb.FilterConfig{
			RaiSettings: &modelarmorpb.RaiFilterSettings{
				RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
					{
						FilterType:      modelarmorpb.RaiFilterType_DANGEROUS,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_HARASSMENT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_HATE_SPEECH,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_SEXUALLY_EXPLICIT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
				},
			},
			SdpSettings: &modelarmorpb.SdpFilterSettings{
				SdpConfiguration: &modelarmorpb.SdpFilterSettings_BasicConfig{
					BasicConfig: &modelarmorpb.SdpBasicConfig{
						FilterEnforcement: modelarmorpb.SdpBasicConfig_ENABLED,
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

	// Prepare the request to update the template.
	req := &modelarmorpb.UpdateTemplateRequest{
		Template: template,
	}

	// Update the template.
	response, err := client.UpdateTemplate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %v", err)
	}

	fmt.Fprintf(w, "Updated Model Armor Template Metadata: %s\n", response.Name)

	// [END modelarmor_update_template_metadata]

	return response, nil
}
