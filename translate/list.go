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

// [START translate_list_codes]
// [START translate_list_language_names]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func listSupportedLanguages(w io.Writer, targetLanguage string) error {
	// targetLanguage := "th"
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return fmt.Errorf("language.Parse: %w", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %w", err)
	}
	defer client.Close()

	langs, err := client.SupportedLanguages(ctx, lang)
	if err != nil {
		return fmt.Errorf("SupportedLanguages: %w", err)
	}

	for _, lang := range langs {
		fmt.Fprintf(w, "%q: %s\n", lang.Tag, lang.Name)
	}

	return nil
}

// [END translate_list_language_names]
// [END translate_list_codes]
