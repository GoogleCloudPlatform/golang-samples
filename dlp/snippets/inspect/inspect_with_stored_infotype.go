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

// [START dlp_inspect_with_stored_infotype]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectWithStoredInfotype inspects the given text using the specified stored infoType detector.
func inspectWithStoredInfotype(w io.Writer, projectID, infoTypeId, textToDeidentify string) error {
	// projectId := "your-project-id"
	// infoTypeId := "your-info-type-id"
	// textToDeidentify := "This commit was made by kewin2010"

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

	// Specify the content to be inspected.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: textToDeidentify,
		},
	}

	// Specify the info type the inspection will look for.
	infoType := &dlppb.InfoType{
		Name: "GITHUB_LOGINS",
	}

	// Specify the stored info type the inspection will look for.
	storedType := &dlppb.StoredType{
		Name: infoTypeId,
	}

	customInfoType := &dlppb.CustomInfoType{
		InfoType: infoType,
		Type: &dlppb.CustomInfoType_StoredType{
			StoredType: storedType,
		},
	}

	// Specify how the content should be inspected.
	inspectConfig := &dlppb.InspectConfig{
		CustomInfoTypes: []*dlppb.CustomInfoType{
			customInfoType,
		},
		IncludeQuote: true,
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent:        fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: inspectConfig,
		Item:          contentItem,
	}

	// Use the client to send the API request.
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

// [END dlp_inspect_with_stored_infotype]
