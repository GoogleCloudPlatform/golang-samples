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

// [START dlp_inspect_phone_number]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectPhoneNumber demonstrates a simple scan request to the Cloud DLP API.
// Notice that the PHONE_NUMBER detector is specified in inspectConfig,
// which instructs Cloud DLP to scan the given string for a phone number.
func inspectPhoneNumber(w io.Writer, projectID, textToInspect string) error {
	// projectID := "my-project-id"
	// textToInspect := "My phone number is (123) 555-6789"

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

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: textToInspect,
			},
		},
		InspectConfig: &dlppb.InspectConfig{
			// Specify the type of info the inspection will look for.
			// See https://cloud.google.com/dlp/docs/infotypes-reference
			// for complete list of info types
			InfoTypes: []*dlppb.InfoType{
				{Name: "PHONE_NUMBER"},
			},
			IncludeQuote: true,
		},
	}

	// Send the request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "receive: %v", err)
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

// [END dlp_inspect_phone_number]
