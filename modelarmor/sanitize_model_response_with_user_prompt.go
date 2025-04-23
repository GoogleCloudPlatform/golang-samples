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

// Sample code for sanitizing a model response using the model armor.

package modelarmor

// [START modelarmor_sanitize_model_response_with_user_prompt]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// sanitizeModelResponseWithUserPrompt sanitizes a model response with a user prompt.
//
// w io.Writer: The writer to use for logging.
// projectID string: The ID of the project.
// locationID string: The ID of the location.
// templateID string: The ID of the template.
// modelResponse string: The model response to sanitize.
// userPrompt string: The user prompt to use for sanitization.
//
// The function returns an error if sanitization fails.
func sanitizeModelResponseWithUserPrompt(w io.Writer, projectID, locationID, templateID, modelResponse, userPrompt string) error {
	ctx := context.Background()

	// Create options for Model Armor client.
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID))
	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Initialize request argument(s)
	modelResponseData := &modelarmorpb.DataItem{
		DataItem: &modelarmorpb.DataItem_Text{
			Text: modelResponse,
		},
	}
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID)
	// Prepare request for sanitizing model response.
	req := &modelarmorpb.SanitizeModelResponseRequest{
		Name:              templateName,
		ModelResponseData: modelResponseData,
		UserPrompt:        userPrompt,
	}

	// Call the API to sanitize the model response.
	response, err := client.SanitizeModelResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to sanitize model response: %w", err)
	}

	fmt.Fprintf(w, "Sanitized response: %s\n", response)

	return err
}

// [END modelarmor_sanitize_model_response_with_user_prompt]
