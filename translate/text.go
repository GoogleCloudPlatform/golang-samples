// Copyright 2019 Google LLC
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

// [START translate_translate_text]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// translateText translates the given text into the specified targetLanguage
//
// Find a list of supported languages and codes here:
// https://cloud.google.com/translate/docs/languages#nmt
func translateText(w io.Writer, targetLanguage, text string) error {
	ctx := context.Background()

	// Get required tag by parsing the target language.
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return fmt.Errorf("language.Parse: %w", err)
	}

	// Create new Translate client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient error: %w", err)
	}
	defer client.Close()

	// Find more information about translate function here:
	// https://cloud.google.com/python/docs/reference/translate/latest/google.cloud.translate_v2.client.Client#google_cloud_translate_v2_client_Client_translate
	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return fmt.Errorf("client.Translate error: %w", err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("client.Translate returned empty response to text: %s", text)
	}

	// Print results to buffer.
	fmt.Fprintf(w, "Text: %s\n", resp[0].Text)
	fmt.Fprintf(w, "Translation: %s\n", text)
	fmt.Fprintf(w, "Detected source language: %s\n", resp[0].Source.String())

	return nil
}

// [END translate_translate_text]
