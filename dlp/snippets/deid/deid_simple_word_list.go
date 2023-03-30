// Copyright 2023 Google LLC
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
package deid

// [START dlp_deidentify_simple_word_list]
import(
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

//deidentifyWithWordList matches against a custom simple word list to de-identify sensitive 
//data based on the input
func deidentifyWithWordList(w io.Writer, projectID, input string, wordList []string) error {
	// projectID := "my-project-id"
	// input := "Patient was seen in RM-YELLOW then transferred to rm green."
	// wordList := []string{"RM-GREEN", "RM-YELLOW", "RM-ORANGE"}
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}
	defer client.Close()

	// Specify what content you want the service to DeIdentify.
	var item = &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: input,
		},
	}

	// Construct the word list to be detected
	words := wordList

	// Specify the word list custom info type the inspection will look for.
	var infoType = &dlppb.InfoType{
		Name: "CUSTOM_ROOM_ID",
	}

	var customInfoType = &dlppb.CustomInfoType{
		InfoType: infoType,
		Type: &dlppb.CustomInfoType_Dictionary_{
			Dictionary: &dlppb.CustomInfoType_Dictionary{
				Source: &dlppb.CustomInfoType_Dictionary_WordList_{
					WordList: &dlppb.CustomInfoType_Dictionary_WordList{
						Words: words, //[]string{"RM-GREEN", "RM-YELLOW", "RM-ORANGE"}
					},
				},
			},
		},
	}

	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			CustomInfoTypes: []*dlppb.CustomInfoType{
				customInfoType,
			},
		},
		// Construct the configuration for the Redact request and list all desired transformations.
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					// Associate deidentification type with info type.
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{infoType},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{ // Define type of deidentification as replacement.
								Transformation: &dlppb.PrimitiveTransformation_ReplaceWithInfoTypeConfig{
									ReplaceWithInfoTypeConfig: &dlppb.ReplaceWithInfoTypeConfig{},
								},
							},
						},
					},
				},
			},
		},
		// The item to analyze.
		Item: item,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return fmt.Errorf("DeidentifyContent: %v", err)
	}

	// Print the result.
	fmt.Fprintf(w, "output : %v", resp.GetItem().GetValue())
	return nil
}
// [END dlp_deidentify_simple_word_list]