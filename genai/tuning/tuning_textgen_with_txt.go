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

// Package tuning shows how to use the GenAI SDK for tuning jobs.
package tuning

// [START googlegenaisdk_tuning_textgen_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// predictWithTunedEndpoint demonstrates how to send a text generation request
// to a tuned endpoint created from a tuning job.
func predictWithTunedEndpoint(w io.Writer, tuningJobName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Retrieve the tuning job and its tuned model endpoint.
	tuningJob, err := client.Tunings.Get(ctx, tuningJobName, nil)
	if err != nil {
		return fmt.Errorf("failed to get tuning job: %w", err)
	}

	contents := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: "Why is the sky blue?"},
			},
		},
	}

	// Send prediction request to the tuned endpoint.
	resp, err := client.Models.GenerateContent(ctx,
		tuningJob.TunedModel.Endpoint,
		contents,
		nil,
	)
	if err != nil {
		return fmt.Errorf("generate content failed: %w", err)
	}

	fmt.Fprintln(w, resp.Text())
	// Example response:
	//   The sky is blue because ...

	return nil
}

// [END googlegenaisdk_tuning_textgen_with_txt]
