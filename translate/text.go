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

// translateText translates the given text into the specified targetLanguage. sourceLanguage
// is optional. If empty, the API will attempt to detect the source language automatically.
// targetLanguage and sourceLanguage should follow ISO 639 language code of the input text
// (e.g., 'fr' for French)
//
// Find a list of supported languages and codes here:
// https://cloud.google.com/translate/docs/languages#nmt
func translateText(w io.Writer, targetLanguage, sourceLanguage, text string) error {
	ctx := context.Background()

	// Create new Translate client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient error: %w", err)
	}
	defer client.Close()

	// Get required tag by parsing the target language.
	targetLang, err := language.Parse(targetLanguage)
	if err != nil {
		return fmt.Errorf("language.Parse: %w", err)
	}

	options := &translate.Options{}

	if sourceLanguage != "" {
		sourceLang, err := language.Parse(sourceLanguage)
		if err != nil {
			return fmt.Errorf("language.Parse: %w", err)
		}
		options = &translate.Options{
			Source: sourceLang,
		}
	}

	// Find more information about translate function here:
	// https://pkg.go.dev/cloud.google.com/go/translate#Client.Translate
	resp, err := client.Translate(ctx, []string{text}, targetLang, options)
	if err != nil {
		return fmt.Errorf("client.Translate error: %w", err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("client.Translate returned empty response to text: %s", text)
	}

	// Print results to buffer.
	fmt.Fprintf(w, "Input Text: %s\n", resp[0].Text)
	fmt.Fprintf(w, "Translated Test: %s\n", text)

	return nil
}

// [END translate_translate_text]
