// Copyright 2024 Google LLC
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

package streamtextbasic

// [START generativeaionvertexai_stream_text_basic]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/iterator"
)

// generateContent shows how to	send a basic streaming text prompt, writing
// the response to the provided io.Writer.
func generateContent(w io.Writer, projectID, modelName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, "us-central1")
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	iter := model.GenerateContentStream(
		ctx,
		genai.Text("Write a story about a magic backpack."),
	)
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return errors.New("empty response from model")
		}
		if err != nil {
			return err
		}
		fmt.Fprint(w, "generated response: ")
		for _, c := range resp.Candidates {
			for _, p := range c.Content.Parts {
				fmt.Fprintf(w, "%s ", p)
			}
		}
	}
}

// [END generativeaionvertexai_stream_text_basic]
