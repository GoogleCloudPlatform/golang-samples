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

// [START dlp_deidentify_replace]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyWithReplacement de-identifies sensitive data by replacing matched input values
func deidentifyWithReplacement(w io.Writer, projectID, inputStr string, infoTypeNames []string, replaceVal string) error {
	// projectId := "your-project-id"
	// inputStr := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."
	// infoTypeNames := []string{"EMAIL_ADDRESS"}
	// replaceVal := "[email-address]"

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

	// item to be analyzed
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{Value: inputStr},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	var infoTypes []*dlppb.InfoType
	for _, v := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: v})
	}
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: infoTypes,
	}

	// Specify replacement string to be used for the finding.
	replaceValueConfig := &dlppb.ReplaceValueConfig{
		NewValue: &dlppb.Value{
			Type: &dlppb.Value_StringValue{
				StringValue: replaceVal,
			},
		},
	}

	// Define type of de-identification as replacement.
	primitiveTransformation := &dlppb.PrimitiveTransformation_ReplaceConfig{
		ReplaceConfig: replaceValueConfig,
	}

	// Associate de-identification type with info type.
	infoTypeTransformation := &dlppb.InfoTypeTransformations_InfoTypeTransformation{
		InfoTypes: infoTypes,
		PrimitiveTransformation: &dlppb.PrimitiveTransformation{
			Transformation: primitiveTransformation,
		},
	}

	deIdentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: &dlppb.InfoTypeTransformations{
				Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
					infoTypeTransformation,
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
	fmt.Fprintf(w, "output : %v", resp.GetItem().GetValue())
	return nil
}

// [END dlp_deidentify_replace]
