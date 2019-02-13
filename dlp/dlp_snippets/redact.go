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

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_redact_image]

// redactImage blacks out the identified portions of the input image (with type bytesType)
// and stores the result in outputPath.
func redactImage(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, infoTypes []string, bytesType dlppb.ByteContentItem_BytesType, inputPath, outputPath string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}

	// Convert the info type strings to a list of types to redact in the image.
	var ir []*dlppb.RedactImageRequest_ImageRedactionConfig
	for _, it := range infoTypes {
		ir = append(ir, &dlppb.RedactImageRequest_ImageRedactionConfig{
			Target: &dlppb.RedactImageRequest_ImageRedactionConfig_InfoType{
				InfoType: &dlppb.InfoType{Name: it},
			},
		})
	}

	// Read the input file.
	b, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	// Create a configured request.
	req := &dlppb.RedactImageRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     i,
			MinLikelihood: minLikelihood,
		},
		// The item to analyze.
		ByteItem: &dlppb.ByteContentItem{
			Type: bytesType,
			Data: b,
		},
		ImageRedactionConfigs: ir,
	}
	// Send the request.
	resp, err := client.RedactImage(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Write the output file.
	if err := ioutil.WriteFile(outputPath, resp.GetRedactedImage(), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Wrote output to %s", outputPath)
}

// [END dlp_redact_image]
