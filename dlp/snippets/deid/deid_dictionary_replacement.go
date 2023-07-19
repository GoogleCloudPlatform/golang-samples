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

package deid

// [START dlp_deidentify_dictionary_replacement]

import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyDataReplaceWithDictionary de-identifies sensitive data in a string by replacing
// each piece of detected sensitive data with a value that Cloud DLP randomly selects from
// a list of words that you provide.
func deidentifyDataReplaceWithDictionary(w io.Writer, projectID, textToDeIdentify string) error {
	// projectId := "my-project-id"
	// textToDeIdentify := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."

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
			Value: textToDeIdentify,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoType := &dlppb.InfoType{
		Name: "EMAIL_ADDRESS",
	}

	// Specify the infotypes to be inspected.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			infoType,
		},
	}

	// Specify list of values which will be randomly picked to replace identified email addresses.
	wordList := &dlppb.ReplaceDictionaryConfig_WordList{
		WordList: &dlppb.CustomInfoType_Dictionary_WordList{
			Words: []string{"izumi@example.com", "alex@example.com"},
		},
	}

	// Specify the Dictionary to use for selecting replacement values for the finding.
	replaceDictionaryConfig := &dlppb.ReplaceDictionaryConfig{
		Type: wordList,
	}

	// Define type of de-identification as replacement with items from dictionary.
	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_ReplaceDictionaryConfig{
			ReplaceDictionaryConfig: replaceDictionaryConfig,
		},
	}

	transformation := &dlppb.InfoTypeTransformations_InfoTypeTransformation{
		InfoTypes: []*dlppb.InfoType{
			infoType,
		},
		PrimitiveTransformation: primitiveTransformation,
	}

	// Combine configurations into a request for the service.
	deIdentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: &dlppb.InfoTypeTransformations{
				Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
					transformation,
				},
			},
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: deIdentifyConfig,
		InspectConfig:    inspectConfig,
		Item:             item,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprint(w, "output: ", resp.GetItem().GetValue())
	return nil

}

// [END dlp_deidentify_dictionary_replacement]
