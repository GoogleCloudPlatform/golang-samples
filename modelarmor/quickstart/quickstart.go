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

// [START modelarmor_quickstart]

package main

import (
	"context"
	"fmt"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// Modelarmor quickstart demonstrates how to create a Model Armor template and use it to
// sanitize a user prompt and a model response.
func main() {
	// Google Project ID
	projectID := "your-project-id"
	// Google Cloud Location
	locationID := "us-central1"
	// ID For The Model Armor Template To Create
	templateID := "go-template"

	ctx := context.Background()
	// Initialize Client
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID))
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create client: %w", err)
		fmt.Println(wrappedErr)
	}
	defer client.Close()
	// Setup Template
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Template for Model Armor API requests.
	// This template defines a filter configuration that detects and filters out
	// sensitive content, including hate speech, harassment, sexually explicit content,
	// and dangerous content.
	template := &modelarmorpb.Template{
		FilterConfig: &modelarmorpb.FilterConfig{
			RaiSettings: &modelarmorpb.RaiFilterSettings{
				// Define individual filters for sensitive content detection
				RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
					// Filter for detecting dangerous content with high confidence level
					{
						FilterType:      modelarmorpb.RaiFilterType_DANGEROUS,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
					// Filter for detecting harassment with medium and above confidence level
					{
						FilterType:      modelarmorpb.RaiFilterType_HARASSMENT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
					},
					// Filter for detecting hate speech with high confidence level
					{
						FilterType:      modelarmorpb.RaiFilterType_HATE_SPEECH,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
					},
					// Filter for detecting sexually explicit content with high confidence level
					{
						FilterType:      modelarmorpb.RaiFilterType_SEXUALLY_EXPLICIT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
					},
				},
			},
		},
	}

	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     parent,
		TemplateId: templateID,
		Template:   template,
	}

	createdTemplate, err := client.CreateTemplate(ctx, req)
	if err != nil {
		wrappedErr := fmt.Errorf("Failed to create template: %w", err)
		fmt.Println(wrappedErr)
	}

	fmt.Printf("Created template: %s\n", createdTemplate.Name)

	// Sanitize a user prompt using the created template
	userPrompt := "Unsafe user prompt"
	userPromptSanitizeReq := &modelarmorpb.SanitizeUserPromptRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		UserPromptData: &modelarmorpb.DataItem{
			DataItem: &modelarmorpb.DataItem_Text{
				Text: userPrompt,
			},
		},
	}

	userPromptSanitizeResp, err := client.SanitizeUserPrompt(ctx, userPromptSanitizeReq)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to sanitize user prompt: %w", err)
		fmt.Println(wrappedErr)
	}

	fmt.Printf("Result for User Prompt Sanitization: %v\n", userPromptSanitizeResp.SanitizationResult)

	// Sanitize a model response using the created template
	modelResponse := "Unsanitized model output"
	sanitizeModelRespReq := &modelarmorpb.SanitizeModelResponseRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		ModelResponseData: &modelarmorpb.DataItem{
			DataItem: &modelarmorpb.DataItem_Text{
				Text: modelResponse,
			},
		},
	}

	// Sanitize Model Response
	sanitizeModelRespResp, err := client.SanitizeModelResponse(ctx, sanitizeModelRespReq)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to sanitize model response: %w", err)
		fmt.Println(wrappedErr)
	}

	fmt.Printf("Result for Model Response Sanitization: %v\n", sanitizeModelRespResp.SanitizationResult)
}

// [END modelarmor_quickstart]
