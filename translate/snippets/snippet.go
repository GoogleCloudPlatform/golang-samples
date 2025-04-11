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

package translate_snippets

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func createClientWithKey() {
	// import "cloud.google.com/go/translate"
	// import "google.golang.org/api/option"
	// import "golang.org/x/text/language"
	ctx := context.Background()

	const apiKey = "YOUR_TRANSLATE_API_KEY"
	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Translate(ctx, []string{"Hello, world!"}, language.Russian, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", resp)
}

// [START translate_translate_text]

func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", err
	}
	return resp[0].Text, nil
}

// [END translate_translate_text]
// [START translate_detect_language]

func detectLanguage(text string) (*translate.Detection, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	lang, err := client.DetectLanguage(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return &lang[0][0], nil
}

// [END translate_detect_language]
// [START translate_list_codes]
// [START translate_list_language_names]

func listSupportedLanguages(w io.Writer, targetLanguage string) error {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return err
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	langs, err := client.SupportedLanguages(ctx, lang)
	if err != nil {
		return err
	}

	for _, lang := range langs {
		fmt.Fprintf(w, "%q: %s\n", lang.Tag, lang.Name)
	}

	return nil
}

// [END translate_list_language_names]
// [END translate_list_codes]

// [START translate_text_with_model]

func translateTextWithModel(targetLanguage, text, model string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, &translate.Options{
		Model: model, // Either "mnt" or "base".
	})
	if err != nil {
		return "", err
	}
	return resp[0].Text, nil
}

// [END translate_text_with_model]
