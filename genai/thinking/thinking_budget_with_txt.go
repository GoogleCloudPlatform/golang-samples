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

// Package thinking shows how to use the GenAI SDK to include thoughts with txt.
package thinking

// [START googlegenaisdk_thinking_budget_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateThinkingBudgetContentWithText demonstrates how to generate text including the model's thought process.
func generateThinkingBudgetContentWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	thinkingBudget := int32(1024) //Use `0` to turn off thinking
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: "solve x^2 + 4x + 4 = 0"},
			},
			Role: "user",
		},
	}

	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinkingBudget,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("generate content failed: %w", err)
	}

	if resp.UsageMetadata != nil {
		fmt.Fprintf(w, "Thoughts token count: %d\n", resp.UsageMetadata.ThoughtsTokenCount)
		//Example response:
		//  908
		fmt.Fprintf(w, "Total token count: %d\n", resp.UsageMetadata.TotalTokenCount)
		//Example response:
		//  1364
	}

	fmt.Fprintln(w, resp.Text())

	// Example response:
	//    To solve the equation $x^2 + 4x + 4 = 0$, you can use several methods:
	//    **Method 1: Factoring**
	//    1.  Look for two numbers that multiply to the constant term (4) and add up to the coefficient of the $x$ term (4).
	//    2.  The numbers are 2 and 2 ($2 \times 2 = 4$ and $2 + 2 = 4$).
	//    ...
	//    ...
	//    Both methods yield the same result.
	//    The solution to the equation $x^2 + 4x + 4 = 0$ is **$x = -2$**.

	return nil
}

// [END googlegenaisdk_thinking_budget_with_txt]
