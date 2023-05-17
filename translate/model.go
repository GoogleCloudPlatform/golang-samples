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

// [START translate_text_with_model]
import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func translateTextWithModel(targetLanguage, text, model string) (string, error) {
	// targetLanguage := "ja"
	// text := "The Go Gopher is cute"
	// model := "nmt"

	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %w", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("translate.NewClient: %w", err)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, &translate.Options{
		Model: model, // Either "nmt" or "base".
	})
	if err != nil {
		return "", fmt.Errorf("Translate: %w", err)
	}
	if len(resp) == 0 {
		return "", nil
	}
	return resp[0].Text, nil
}

// [END translate_text_with_model]
