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

package redact

// [START dlp_redact_image]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// redactImage blacks out the identified portions of the input image (with type bytesType)
// and stores the result in outputPath.
func redactImage(w io.Writer, projectID string, infoTypeNames []string, bytesType dlppb.ByteContentItem_BytesType, inputPath, outputPath string) error {
	// projectID := "my-project-id"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	// bytesType := dlppb.ByteContentItem_IMAGE_PNG
	// inputPath := /tmp/input
	// outputPath := /tmp/output

	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}

	// Convert the info type strings to a list of InfoTypes.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	// Convert the info type strings to a list of types to redact in the image.
	var redactInfoTypes []*dlppb.RedactImageRequest_ImageRedactionConfig
	for _, it := range infoTypeNames {
		redactInfoTypes = append(redactInfoTypes, &dlppb.RedactImageRequest_ImageRedactionConfig{
			Target: &dlppb.RedactImageRequest_ImageRedactionConfig_InfoType{
				InfoType: &dlppb.InfoType{Name: it},
			},
		})
	}

	// Read the input file.
	b, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}

	// Create a configured request.
	req := &dlppb.RedactImageRequest{
		Parent: "projects/" + projectID,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     infoTypes,
			MinLikelihood: dlppb.Likelihood_POSSIBLE,
		},
		// The item to analyze.
		ByteItem: &dlppb.ByteContentItem{
			Type: bytesType,
			Data: b,
		},
		ImageRedactionConfigs: redactInfoTypes,
	}
	// Send the request.
	resp, err := client.RedactImage(ctx, req)
	if err != nil {
		return fmt.Errorf("RedactImage: %v", err)
	}
	// Write the output file.
	if err := ioutil.WriteFile(outputPath, resp.GetRedactedImage(), 0644); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %v", err)
	}
	fmt.Fprintf(w, "Wrote output to %s", outputPath)
	return nil
}

// [END dlp_redact_image]
