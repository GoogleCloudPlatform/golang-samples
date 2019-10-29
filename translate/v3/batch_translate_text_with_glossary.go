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

// [START translate_v3_batch_translate_text_with_glossary]
import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// batchTranslateTextWithGlossary translates a large volume of text in asynchronous batch mode.
func batchTranslateTextWithGlossary(w io.Writer, projectID string, location string, inputURI string, outputURI string, sourceLang string, targetLang string, glossaryID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://cloud-samples-data/text.txt"
	// outputURI := "gs://YOUR_BUCKET_ID/path_to_store_results/"
	// sourceLang := "en"
	// targetLang := "ja"
	// glossaryID := "your-glossary-id"

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %v", err)
	}
	defer client.Close()

	req := &translatepb.BatchTranslateTextRequest{
		Parent:              fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		SourceLanguageCode:  sourceLang,
		TargetLanguageCodes: []string{targetLang},
		InputConfigs: []*translatepb.InputConfig{
			{
				Source: &translatepb.InputConfig_GcsSource{
					GcsSource: &translatepb.GcsSource{InputUri: inputURI},
				},
				// Optional. Can be "text/plain" or "text/html".
				MimeType: "text/plain",
			},
		},
		Glossaries: map[string]*translatepb.TranslateTextGlossaryConfig{
			targetLang: {
				Glossary: fmt.Sprintf("projects/%s/locations/%s/glossaries/%s", projectID, location, glossaryID),
			},
		},
		OutputConfig: &translatepb.OutputConfig{
			Destination: &translatepb.OutputConfig_GcsDestination{
				GcsDestination: &translatepb.GcsDestination{
					OutputUriPrefix: outputURI,
				},
			},
		},
	}

	// The BatchTranslateText operation is async.
	op, err := client.BatchTranslateText(ctx, req)
	if err != nil {
		return fmt.Errorf("BatchTranslateText: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Total characters: %v\n", resp.GetTotalCharacters())
	fmt.Fprintf(w, "Translated characters: %v\n", resp.GetTranslatedCharacters())

	return nil
}

// [END translate_v3_batch_translate_text_with_glossary]
