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

// [START translate_v3_detect_language]
import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// detectLanguage detects the language of a text string.
func detectLanguage(w io.Writer, projectID string, text string) error {
	// projectID := "my-project-id"
	// text := "Hello, world!"

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %v", err)
	}
	defer client.Close()

	req := &translatepb.DetectLanguageRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/global", projectID),
		MimeType: "text/plain", // Mime types: "text/plain", "text/html"
		Source: &translatepb.DetectLanguageRequest_Content{
			Content: text,
		},
	}

	resp, err := client.DetectLanguage(ctx, req)
	if err != nil {
		return fmt.Errorf("DetectLanguage: %v", err)
	}

	// Display list of detected languages sorted by detection confidence.
	// The most probable language is first.
	for _, language := range resp.GetLanguages() {
		// The language detected.
		fmt.Fprintf(w, "Language code: %v\n", language.GetLanguageCode())
		// Confidence of detection result for this language.
		fmt.Fprintf(w, "Confidence: %v\n", language.GetConfidence())
	}

	return nil
}

// [END translate_v3_detect_language]
