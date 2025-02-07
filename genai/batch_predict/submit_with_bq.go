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

// Package batch_predict shows examples of generating text using the GenAI SDK. [TODO: Update this]
package batch_predict

// [START googlegenaisdk_batchpredict_with_bq]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// update_me shows how to generate text using a ... [TODO: Add function description]
func update_me(w io.Writer) error {
  ctx := context.Background()

  client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
  if err != nil {
    return fmt.Errorf("failed to create genai client: %w", err)
  }


  // modelName := "gemini-2.0-flash-001"
  // contents := []*genai.Content{
  //   {Parts: []*genai.Part{
  //     {Text: "What's in this video?"},
  //     {FileData: &genai.FileData{
  //       FileURI:  "gs://cloud-samples-data/generative-ai/video/pixel8.mp4",
  //       MIMEType: "video/mp4",
  //     }},
  //   }},
  // }
  //
  // resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to generate content: %w", err)
	// }
  //
	// respText, err := resp.Text()
	// if err != nil {
	// 	return fmt.Errorf("failed to convert model response to text: %w", err)
	// }
	// fmt.Fprintln(w, respText)

  // Example response:
  // ...

  return nil
}

// [END googlegenaisdk_batchpredict_with_bq]
