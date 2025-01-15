// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspect

// [START dlp_inspect_augment_infotypes]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectAugmentInfoTypes performs info type augmentation using Google Cloud DLP.
// It enhances data inspection by supplementing existing info types with custom-defined ones,
// expanding the ability to identify sensitive information in different contexts.
func inspectAugmentInfoTypes(w io.Writer, projectID, textToInspect string, wordList []string) error {
	// projectID := "your-project-id"
	// textToInspect := "The patient's name is quasimodo"
	// wordList := []string{"quasimodo"}

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the content to be inspected.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: textToInspect,
		},
	}

	// Construct the custom word list to be detected.
	dictionary := &dlppb.CustomInfoType_Dictionary{
		Source: &dlppb.CustomInfoType_Dictionary_WordList_{
			WordList: &dlppb.CustomInfoType_Dictionary_WordList{
				Words: wordList,
			},
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types.
	infoType := &dlppb.InfoType{
		Name: "PERSON_NAME",
	}

	// Construct a custom infoType detector by augmenting the PERSON_NAME detector with a word list.
	customInfoType := &dlppb.CustomInfoType{
		InfoType: infoType,
		Type: &dlppb.CustomInfoType_Dictionary_{
			Dictionary: dictionary,
		},
	}

	// Specify the inspect config for data inspection settings in DLP API, enabling rule
	// specification, custom info types, and actions on sensitive data. Crucial for tailored
	// data protection and privacy regulation compliance.
	inspectConfig := &dlppb.InspectConfig{
		CustomInfoTypes: []*dlppb.CustomInfoType{
			customInfoType,
		},
		IncludeQuote: true,
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent:        fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:          item,
		InspectConfig: inspectConfig,
	}

	// Create the request for the job configured above.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return err
	}

	// Process the results.
	result := resp.Result
	fmt.Fprintf(w, "Findings: %d\n", len(result.Findings))
	for _, f := range result.Findings {
		fmt.Fprintf(w, "\tQuote: %s\n", f.Quote)
		fmt.Fprintf(w, "\tInfo type: %s\n", f.InfoType.Name)
		fmt.Fprintf(w, "\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_augment_infotypes]
