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
package redact

// [START dlp_redact_image_listed_infotypes]
import (
	"context"
	"fmt"
	"io"
	"os"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// redactImageFileListedInfoTypes redacts only certain sensitive
// data from an image using infoTypes
func redactImageFileListedInfoTypes(w io.Writer, projectID, inputPath, outputPath string) error {
	// projectId := "my-project-id"
	// inputPath := "testdata/image.jpg"
	// outputPath := "testdata/test-output-image-file-listed-infoTypes-redacted.jpeg"

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

	// read the image file
	fileBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Specify the content to be redacted.
	byteItem := &dlppb.ByteContentItem{
		Type: dlppb.ByteContentItem_IMAGE_JPEG,
		Data: fileBytes,
	}

	// Specify the types of info necessary to redact.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoTypes := []*dlppb.InfoType{
		{Name: "US_SOCIAL_SECURITY_NUMBER"},
		{Name: "EMAIL_ADDRESS"},
		{Name: "PHONE_NUMBER"},
	}

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: infoTypes,
	}

	// Prepare redaction configs.
	var x []*dlppb.RedactImageRequest_ImageRedactionConfig
	for _, v := range infoTypes {
		x = append(x, &dlppb.RedactImageRequest_ImageRedactionConfig{Target: &dlppb.RedactImageRequest_ImageRedactionConfig_InfoType{InfoType: v}})
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.RedactImageRequest{
		Parent:                fmt.Sprintf("projects/%s/locations/global", projectID),
		ByteItem:              byteItem,
		ImageRedactionConfigs: x,
		InspectConfig:         inspectConfig,
	}

	// Send the request.
	resp, err := client.RedactImage(ctx, req)
	if err != nil {
		return err
	}

	// Write the output file.
	if err := os.WriteFile(outputPath, resp.GetRedactedImage(), 0644); err != nil {
		return err
	}
	fmt.Fprintf(w, "Wrote output to %s\n", outputPath)
	return nil
}

// [END dlp_redact_image_listed_infotypes]
