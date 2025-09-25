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

// [START googlegenaisdk_thinking_includethoughts_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateContentWithThoughts demonstrates how to generate text including the model's thought process.
func generateContentWithThoughts(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-pro"
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
				IncludeThoughts: true,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no content was generated")
	}

	// The response may contain both the final answer and the model's thoughts.
	// Iterate through the parts to print them separately.
	fmt.Fprintln(w, "Answer:")
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" && !part.Thought {
			fmt.Fprintln(w, part.Text)
		}
	}
	fmt.Fprintln(w, "\nThoughts:")
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Thought {
			fmt.Fprintln(w, part.Text)
		}
	}

	// Example response:
	//  Answer:
	//	Of course! We can solve this quadratic equation in a couple of ways.
	//
	//### Method 1: Factoring (the easiest method for this problem)
	//
	//1.  **Recognize the pattern.** The expression `x² + 4x + 4` is a perfect square trinomial. It fits the pattern `a² + 2ab + b² = (a + b)²`. In this case, `a = x` and `b = 2`.
	//
	//2.  **Factor the equation.**
	//    `x² + 4x + 4 = (x + 2)(x + 2) = (x + 2)²`
	//
	//3.  **Solve for x.** Now set the factored expression to zero:
	//    `(x + 2)² = 0`
	//
	//    Take the square root of both sides:
	//    `x + 2 = 0`
	//
	//    Subtract 2 from both sides:
	//    `x = -2`
	//
	//This type of solution is called a "repeated root" or a "double root" because the factor `(x+2)` appears twice.
	//
	//---
	//
	//### Method 2: Using the Quadratic Formula
	//
	//You can use the quadratic formula for any equation in the form `ax² + bx + c = 0`.
	//
	//The formula is: `x = [-b ± sqrt(b² - 4ac)] / 2a`
	//
	//1.  **Identify a, b, and c.**
	//    *   a = 1
	//    *   b = 4
	//    *   c = 4
	//
	//2.  **Plug the values into the formula.**
	//    `x = [-4 ± sqrt(4² - 4 * 1 * 4)] / (2 * 1)`
	//
	//3.  **Simplify.**
	//    `x = [-4 ± sqrt(16 - 16)] / 2`
	//    `x = [-4 ± sqrt(0)] / 2`
	//    `x = -4 / 2`
	//
	//4.  **Solve for x.**
	//    `x = -2`
	//Alright, the user wants to solve the quadratic equation `x² + 4x + 4 = 0`. My first instinct is to see if I can factor it; that's often the fastest approach if it works.  Looking at the coefficients, I see `a = 1`, `b = 4`, and `c = 4`.  Factoring is clearly the most direct path here. I need to find two numbers that multiply to 4 (c) and add up to 4 (b). Hmm, let's see… 1 and 4? Nope, that adds to 5.  2 and 2? Perfect!  2 times 2 is 4, and 2 plus 2 is also 4.
	//
	//So, `x² + 4x + 4` factors nicely into `(x + 2)(x + 2)`.  Ah, a perfect square trinomial! That's useful to note. Now, I can write the equation as `(x + 2)² = 0`.  Taking the square root of both sides gives me `x + 2 = 0`.  And finally, subtracting 2 from both sides, I get `x = -2`.  That's the solution.
	//
	//Just to be thorough, and maybe to offer an alternative explanation, let's verify this using the quadratic formula. It's `x = [-b ± √(b² - 4ac)] / 2a`. Plugging in my values:  `x = [-4 ± √(4² - 4 * 1 * 4)] / (2 * 1)`.  That simplifies to `x = [-4 ± √(16 - 16)] / 2`, or `x = [-4 ± 0] / 2`.  Therefore, `x = -2`. The discriminant being zero tells me I have exactly one real, repeated root.  Great. So, whether I factor or use the quadratic formula, the answer is the same.
	return nil
}

// [END googlegenaisdk_thinking_includethoughts_with_txt]
