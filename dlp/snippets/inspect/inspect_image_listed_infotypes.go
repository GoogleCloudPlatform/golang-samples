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

// [START dlp_inspect_image_listed_infotypes]
import (
	"context"
	"fmt"
	"io"
	"os"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectImageFileListedInfoTypes inspects an image for sensitive data listed infoTypes.
func inspectImageFileListedInfoTypes(w io.Writer, projectID, filePath string) error {
	// projectId := "my-project-id"
	// filePath := "testdata/image.jpg"

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

	// Read a image file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Specify the type and content to be inspected.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: dlppb.ByteContentItem_IMAGE_JPEG,
				Data: data,
			},
		},
	}

	// Construct the configuration for the Inspect request.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			{Name: "PHONE_NUMBER"},
			{Name: "EMAIL_ADDRESS"},
			{Name: "US_SOCIAL_SECURITY_NUMBER"},
		},
		IncludeQuote: true,
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent:        fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:          item,
		InspectConfig: inspectConfig,
	}

	// Send the request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return err
	}

	// Process the results.
	fmt.Fprintf(w, "Findings: %d\n", len(resp.Result.Findings))
	for _, f := range resp.Result.Findings {
		fmt.Fprintf(w, "\tQuote: %s\n", f.Quote)
		fmt.Fprintf(w, "\tInfo type: %s\n", f.InfoType.Name)
		fmt.Fprintf(w, "\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_image_listed_infotypes]
