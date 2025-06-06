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

// Package v3 contains samples for Google Cloud Translation API v3.
package v3

// [START translate_v3_translate_text]
// Imports the Google Cloud Translation library
import (
	"context"
	"fmt"
	"io"

	// [START translate_v3_import_client_library]
	translate "cloud.google.com/go/translate/apiv3"
	// [END translate_v3_import_client_library]
	"cloud.google.com/go/translate/apiv3/translatepb"
)

func translateText(w io.Writer, projectID string, sourceLang string, targetLang string, text string) error {
	// projectID := "your-project-id"
	// sourceLang := "en-US"
	// targetLang := "fr"
	// text := "Text you wish to translate"

	// Instantiates a client
	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %w", err)
	}
	defer client.Close()

	// Construct request
	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		SourceLanguageCode: sourceLang,
		TargetLanguageCode: targetLang,
		MimeType:           "text/plain", // Mime types: "text/plain", "text/html"
		Contents:           []string{text},
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return fmt.Errorf("TranslateText: %w", err)
	}

	// Display the translation for each input text provided
	for _, translation := range resp.GetTranslations() {
		fmt.Fprintf(w, "Translated text: %v\n", translation.GetTranslatedText())
	}

	return nil
}

// [END translate_v3_translate_text]
