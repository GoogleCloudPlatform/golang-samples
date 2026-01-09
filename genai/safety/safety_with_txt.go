// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package safety shows how to use the GenAI SDK to generate safety content.
package safety

// [START googlegenaisdk_safety_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateTextWithSafety shows how to apply safety settings to a text generation request.
func generateTextWithSafety(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	systemInstruction := &genai.Content{
		Parts: []*genai.Part{
			{Text: "Be as mean as possible."},
		},
		Role: genai.RoleUser,
	}

	prompt := "Write a list of 5 disrespectful things that I might say to the universe after stubbing my toe in the dark."

	safetySettings := []*genai.SafetySetting{
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: systemInstruction,
		SafetySettings:    safetySettings,
	}
	modelName := "gemini-2.5-flash"
	resp, err := client.Models.GenerateContent(ctx, modelName,
		[]*genai.Content{{Parts: []*genai.Part{{Text: prompt}}, Role: genai.RoleUser}},
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Fprintln(w, resp.Text())

	if len(resp.Candidates) > 0 {
		fmt.Fprintln(w, "Finish Reason:", resp.Candidates[0].FinishReason)

		for _, rating := range resp.Candidates[0].SafetyRatings {
			fmt.Fprintf(w, "\nCategory: %v\nIs Blocked: %v\nProbability: %v\nProbability Score: %v\nSeverity: %v\nSeverity Score: %v\n",
				rating.Category,
				rating.Blocked,
				rating.Probability,
				rating.ProbabilityScore,
				rating.Severity,
				rating.SeverityScore,
			)
		}
	}

	// Example response:
	// Category: HARM_CATEGORY_HATE_SPEECH
	// Is Blocked: false
	// Probability: NEGLIGIBLE
	// Probability Score: 8.996795e-06
	// Severity: HARM_SEVERITY_NEGLIGIBLE
	// Severity Score: 0.04771039
	//
	// Category: HARM_CATEGORY_DANGEROUS_CONTENT
	// Is Blocked: false
	// Probability: NEGLIGIBLE
	// Probability Score: 2.2431707e-06
	// Severity: HARM_SEVERITY_NEGLIGIBLE
	// Severity Score: 0
	//
	// Category: HARM_CATEGORY_HARASSMENT
	// Is Blocked: false
	// Probability: NEGLIGIBLE
	// Probability Score: 0.00026123362
	// Severity: HARM_SEVERITY_NEGLIGIBLE
	// Severity Score: 0.022358216
	//
	// Category: HARM_CATEGORY_SEXUALLY_EXPLICIT
	// Is Blocked: false
	// Probability: NEGLIGIBLE
	// Probability Score: 6.1352006e-07
	// Severity: HARM_SEVERITY_NEGLIGIBLE
	// Severity Score: 0.020111412

	return nil
}

// [END googlegenaisdk_safety_with_txt]
