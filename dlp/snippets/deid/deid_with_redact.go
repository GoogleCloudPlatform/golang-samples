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

// [START dlp_deidentify_redact]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyWithRedact de-identify the data by redacting with matched input values
func deidentifyWithRedact(w io.Writer, projectID, inputStr string, infoTypeNames []string) error {
	// projectID := "my-project-id"
	// inputStr := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."
	// infoTypeNames := []string{"EMAIL_ADDRESS"}

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the content to be inspected.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: inputStr,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: infoTypes,
	}

	// Define type of de-identification.
	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_RedactConfig{
			RedactConfig: &dlppb.RedactConfig{},
		},
	}

	// Associate de-identification type with info type.
	transformation := &dlppb.InfoTypeTransformations_InfoTypeTransformation{
		InfoTypes:               infoTypes,
		PrimitiveTransformation: primitiveTransformation,
	}

	// Construct the configuration for the Redact request and list all desired transformations.
	redactConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: &dlppb.InfoTypeTransformations{
				Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
					transformation,
				},
			},
		},
	}

	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: redactConfig,
		InspectConfig:    inspectConfig,
		Item:             contentItem,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the result.
	fmt.Fprintf(w, "output: %v", resp.GetItem().GetValue())
	return nil
}

// [END dlp_deidentify_redact]
