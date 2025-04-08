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

// Sample code for creating a new model armor template with basic SDP settings enabled.

package modelarmor

// [START modelarmor_create_template_with_basic_sdp]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// createModelArmorTemplateWithBasicSDP creates a new Model Armor template with basic SDP settings.
//
// This method creates a new Model Armor template with basic SDP settings.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	projectID string: The ID of the Google Cloud project.
//	locationID string: The ID of the Google Cloud location.
//	templateID string: The ID of the template to create.
//
// Returns:
//
//	*modelarmorpb.Template: The created template.
//	error: Any error that occurred during template creation.
//
// Example:
//
//	template, err := createModelArmorTemplateWithBasicSDP(
//	    os.Stdout,
//	    "my-project",
//	    "us-central1",
//	    "my-template",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(template)
func createModelArmorTemplateWithBasicSDP(w io.Writer, projectID, locationID, templateID string) (*modelarmorpb.Template, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
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
		return nil, fmt.Errorf("failed to create template: %v", err)
	}

	// Print the new template name using fmt.Fprintf with the io.Writer.
	fmt.Fprintf(w, "Created Template with basic SDP: %s\n", response.Name)

	// [END modelarmor_create_template_with_basic_sdp]

	return response, nil
}
