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

// [START translate_v3_list_glossary]
import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	"google.golang.org/api/iterator"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// listGlossaries gets the specified glossary.
func listGlossaries(w io.Writer, projectID string, location string) error {
	// projectID := "my-project-id"
	// location := "us-central1"

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %v", err)
	}
	defer client.Close()

	req := &translatepb.ListGlossariesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.ListGlossaries(ctx, req)

	// Iterate over all results
	for {
		glossary, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListGlossaries.Next: %v", err)
		}
		fmt.Fprintf(w, "Name: %v\n", glossary.GetName())
		fmt.Fprintf(w, "Entry count: %v\n", glossary.GetEntryCount())
		fmt.Fprintf(w, "Input URI: %v\n", glossary.GetInputConfig().GetGcsSource().GetInputUri())
		for _, languageCode := range glossary.GetLanguageCodesSet().GetLanguageCodes() {
			fmt.Fprintf(w, "Language code: %v\n", languageCode)
		}
		if languagePair := glossary.GetLanguagePair(); languagePair != nil {
			fmt.Fprintf(w, "Language pair: %v, %v\n",
				languagePair.GetSourceLanguageCode(), languagePair.GetTargetLanguageCode())
		}
	}

	return nil
}

// [END translate_v3_list_glossary]
