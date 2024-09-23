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

package snippets

// [START apikeys_authenticate_api_key]
import (
	"context"
	"fmt"
	"io"

	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
	"google.golang.org/api/option"
)

// authenticateWithAPIKey authenticates with an API key for Google Language
// service.
func authenticateWithAPIKey(w io.Writer, apiKey string) error {
	// apiKey := "api-key-string"

	ctx := context.Background()

	// Initialize the Language Service client and set the API key.
	client, err := language.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	text := "Hello, world!"
	// Make a request to analyze the sentiment of the text.
	res, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
	if err != nil {
		return fmt.Errorf("AnalyzeSentiment: %w", err)
	}

	fmt.Fprintf(w, "Text: %s\n", text)
	fmt.Fprintf(w, "Sentiment score: %v\n", res.DocumentSentiment.Score)
	fmt.Fprintln(w, "Successfully authenticated using the API key.")

	return nil
}

// [END apikeys_authenticate_api_key]
