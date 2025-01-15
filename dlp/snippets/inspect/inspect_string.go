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

// [START dlp_inspect_string]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectString inspects the a given string, and prints results.
func inspectString(w io.Writer, projectID, textToInspect string) error {
	// projectID := "my-project-id"
	// textToInspect := "My name is Gary and my email is gary@example.com"
	ctx := context.Background()

	// Initialize client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: textToInspect,
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
		return err
	}

	// Process the results.
	result := resp.Result
	fmt.Fprintf(w, "Findings: %d\n", len(result.Findings))
	for _, f := range result.Findings {
		fmt.Fprintf(w, "\tQuote: %s\n", f.Quote)
		fmt.Fprintf(w, "\tInfo type: %s\n", f.InfoType.Name)
		fmt.Fprintf(w, "\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_string]
