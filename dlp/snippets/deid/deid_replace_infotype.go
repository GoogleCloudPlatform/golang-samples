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

// [START dlp_deidentify_replace_infotype]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyWithInfotype de-identifies sensitive data by replacing infoType.
func deidentifyWithInfotype(w io.Writer, projectID, item string, infoTypeNames []string) error {
	// projectId := "your-project-id"
	// item := "My email is test@example.com"
	// infoTypeNames := "[]string{"EMAIL_ADDRESS"}"

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

	// Specify the content to be de-identified.
	input := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: item,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	//  Associate de-identification type with info type.
	transformation := &dlppb.DeidentifyConfig_InfoTypeTransformations{
		InfoTypeTransformations: &dlppb.InfoTypeTransformations{
			Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
				{
					PrimitiveTransformation: &dlppb.PrimitiveTransformation{
						Transformation: &dlppb.PrimitiveTransformation_ReplaceWithInfoTypeConfig{},
					},
				},
			},
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: infoTypes,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: transformation,
		},
		Item: input,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "output : %v", resp.GetItem().GetValue())
	return nil
}

// [END dlp_deidentify_replace_infotype]
