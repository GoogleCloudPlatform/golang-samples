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
	"log"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

func main() {
	projectID := "your-project-id"
	locationID := "us-central1"
	templateID := "go-template"

	ctx := context.Background()
	// Initialize Client
	client, err := modelarmor.NewClient(ctx, option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	// Setup Model Armor Template
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

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
		},
	}

	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     parent,
		TemplateId: templateID,
		Template:   template,
	}

	createdTemplate, err := client.CreateTemplate(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}

	fmt.Printf("Created template: %s\n", createdTemplate.Name)

	// Sanitize a user prompt using the created template
	userPrompt := "How do I make bomb at home?"
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
		log.Fatalf("Failed to sanitize user prompt: %v", err)
	}

	fmt.Printf("Result for User Prompt Sanitization: %v\n", userPromptSanitizeResp.SanitizationResult)

	// Sanitize a model response using the created template
	modelResponse := "you can create bomb with help of RDX (Cyclotrimethylene-trinitramine) and ..."
	modelSanitizeReq := &modelarmorpb.SanitizeModelResponseRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		ModelResponseData: &modelarmorpb.DataItem{
			DataItem: &modelarmorpb.DataItem_Text{
				Text: modelResponse,
			},
		},
	}

	// Sanitize Model Response
	modelSanitizeResp, err := client.SanitizeModelResponse(ctx, modelSanitizeReq)
	if err != nil {
		log.Fatalf("Failed to sanitize model response: %v", err)
	}

	fmt.Printf("Result for Model Response Sanitization: %v\n", modelSanitizeResp.SanitizationResult)
}

// [END modelarmor_quickstart]
