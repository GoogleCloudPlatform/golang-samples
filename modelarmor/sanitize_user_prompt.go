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

// Sample code for sanitizing user prompt with model armor.

package modelarmor

// [START modelarmor_sanitize_user_prompt]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// sanitizeUserPrompt sanitizes a user prompt.
//
// This method sanitizes a user prompt based on the project, location, and template settings.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	projectID string: The ID of the project.
//	locationID string: The ID of the location.
//	templateID string: The ID of the template.
//	userPrompt string: The user prompt to sanitize.
//
// Returns:
//
//	*modelarmorpb.SanitizeUserPromptResponse: The sanitized user prompt.
//	error: Any error that occurred during sanitization.
//
// Example:
//
//	sanitizedPrompt, err := sanitizeUserPrompt(
//	    os.Stdout,
//	    "my-project",
//	    "my-location",
//	    "my-template",
//	    "user prompt",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(sanitizedPrompt)
func sanitizeUserPrompt(w io.Writer, projectID, locationID, templateID, userPrompt string) (*modelarmorpb.SanitizeUserPromptResponse, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize request argument(s)
	userPromptData := &modelarmorpb.DataItem{
		DataItem: &modelarmorpb.DataItem_Text{
			Text: userPrompt,
		},
	}

	// Prepare request for sanitizing user prompt.
	req := &modelarmorpb.SanitizeUserPromptRequest{
		Name:           fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		UserPromptData: userPromptData,
	}

	// Sanitize the user prompt.
	response, err := client.SanitizeUserPrompt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize user prompt: %v", err)
	}

	// Sanitization Result.
	fmt.Fprintf(w, "Sanitization Result: %v\n", response)

	// [END modelarmor_sanitize_user_prompt]

	return response, nil
}
