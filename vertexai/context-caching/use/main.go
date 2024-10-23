// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// contextcaching shows an example of caching the tokens of a mulitple PDF prompt
package use

// [START generativeaionvertexai_gemini_use_context_cache]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

// useContextCache shows how to use an existing cached content, when prompting the model
// contentName is the ID of the cached content
func UseContextCache(w io.Writer, contentName string, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-pro-001"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.CachedContentName = contentName
	prompt := genai.Text("What are the papers about?")

	res, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])
	return nil
}

// [END generativeaionvertexai_gemini_use_context_cache]

func main() {
	err := UseContextCache(
		os.Stdout,
		// TODO(developer): Update the argument values
		"projects/[PROJECT_ID]/locations/us-central1/cachedContents/[CACHE_ID]",
		"acme-corp-dev",
		"us-central1",
		"gemini-1.5-pro-001",
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
