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

// [START translate_v3_get_supported_languages]
import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// getSupportedLanguages gets a list of supported language codes.
func getSupportedLanguages(w io.Writer, projectID string) error {
	// projectID := "my-project-id"

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %v", err)
	}
	defer client.Close()

	req := &translatepb.GetSupportedLanguagesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
	}

	resp, err := client.GetSupportedLanguages(ctx, req)
	if err != nil {
		return fmt.Errorf("GetSupportedLanguages: %v", err)
	}

	// List language codes of supported languages
	fmt.Fprintf(w, "Supported languages:\n")
	for _, language := range resp.GetLanguages() {
		fmt.Fprintf(w, "Language code: %v\n", language.GetLanguageCode())
	}

	return nil
}

// [END translate_v3_get_supported_languages]
