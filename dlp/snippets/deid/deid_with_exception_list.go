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

// [START dlp_deidentify_exception_list]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyExceptionList creates an exception list for a regular custom dictionary detector.
func deidentifyExceptionList(w io.Writer, projectID, input, infoType string, dictWordList []string) error {
	// projectID := "my-project-id"
	// input := "jack@example.org accessed customer record of user5@example.com"
	// infoType := "EMAIL_ADDRESS"
	// dictWordList := []string{"jack@example.org", "jill@example.org"}
	ctx := context.Background()
	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
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

	// Specify the word list custom info type and build-in info type the inspection will look for.
	var infoTypes = []*dlppb.InfoType{
		{Name: infoType},
	}
	inspectRuleSet := &dlppb.InspectionRuleSet{
		InfoTypes: infoTypes,
		Rules: []*dlppb.InspectionRule{
			{
				Type: &dlppb.InspectionRule_ExclusionRule{
					ExclusionRule: &dlppb.ExclusionRule{
						MatchingType: dlppb.MatchingType_MATCHING_TYPE_FULL_MATCH,
						Type: &dlppb.ExclusionRule_Dictionary{
							Dictionary: &dlppb.CustomInfoType_Dictionary{
								Source: &dlppb.CustomInfoType_Dictionary_WordList_{
									WordList: &dlppb.CustomInfoType_Dictionary_WordList{
										Words: dictWordList,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		// Construct the configuration for the de-id request and list all desired transformations.
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_ReplaceWithInfoTypeConfig{},
							},
						},
					},
				},
			},
		},
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:       infoTypes,
			CustomInfoTypes: []*dlppb.CustomInfoType{},
			RuleSet:         []*dlppb.InspectionRuleSet{inspectRuleSet},
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

// [END dlp_deidentify_exception_list]
