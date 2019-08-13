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

package inspect

// [START dlp_inspect_file]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// inspectTextFile inspects a text file at a given filePath, and prints results.
func inspectTextFile(w io.Writer, projectID, filePath string) error {
	// projectID := "my-project-id"
	// filePath := "path/to/image.png"
	ctx := context.Background()

	// Initialize client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Gather the resources for the request.
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + projectID,
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_ByteItem{
				ByteItem: &dlppb.ByteContentItem{
					Type: dlppb.ByteContentItem_TEXT_UTF8,
					Data: data,
				},
			},
		},
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{Name: "PHONE_NUMBER"},
				{Name: "EMAIL_ADDRESS"},
				{Name: "CREDIT_CARD_NUMBER"},
			},
			IncludeQuote: true,
		},
	}
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return fmt.Errorf("InspectContent: %v", err)
	}

	// Process the results.
	fmt.Fprintf(w, "Findings: %d\n", len(resp.Result.Findings))
	for _, f := range resp.Result.Findings {
		fmt.Fprintf(w, "\tQoute: %s\n", f.Quote)
		fmt.Fprintf(w, "\tInfo type: %s\n", f.InfoType.Name)
		fmt.Fprintf(w, "\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_file]
