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

// [START translate_v3_create_glossary]
import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

// createGlossary creates a glossary to use for other operations.
func createGlossary(w io.Writer, projectID string, location string, glossaryID string, glossaryInputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// glossaryID := "my-glossary-display-name"
	// glossaryInputURI := "gs://cloud-samples-data/translation/glossary.csv"

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %w", err)
	}
	defer client.Close()

	req := &translatepb.CreateGlossaryRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Glossary: &translatepb.Glossary{
			Name: fmt.Sprintf("projects/%s/locations/%s/glossaries/%s", projectID, location, glossaryID),
			Languages: &translatepb.Glossary_LanguageCodesSet_{
				LanguageCodesSet: &translatepb.Glossary_LanguageCodesSet{
					LanguageCodes: []string{"en", "ja"},
				},
			},
			InputConfig: &translatepb.GlossaryInputConfig{
				Source: &translatepb.GlossaryInputConfig_GcsSource{
					GcsSource: &translatepb.GcsSource{
						InputUri: glossaryInputURI,
					},
				},
			},
		},
	}

	// The CreateGlossary operation is async.
	op, err := client.CreateGlossary(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateGlossary: %w", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Created: %v\n", resp.GetName())
	fmt.Fprintf(w, "Input URI: %v\n", resp.InputConfig.GetGcsSource().GetInputUri())

	return nil
}

// [END translate_v3_create_glossary]
