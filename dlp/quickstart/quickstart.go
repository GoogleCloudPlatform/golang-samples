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

// [START dlp_quickstart]

// The quickstart program is an example of using the Data Loss Prevention API.
package main

import (
	"context"
	"fmt"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

func main() {
	ctx := context.Background()

	projectID := "PROJECT_ID"

	// Creates a DLP client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating DLP client: %v", err)
	}
	defer client.Close()

	// The string to inspect.
	input := "Robert Frost"

	// The minimum likelihood required before returning a match.
	minLikelihood := dlppb.Likelihood_POSSIBLE

	// The maximum number of findings to report (0 = server maximum).
	maxFindings := int32(0)

	// Whether to include the matching string.
	includeQuote := true

	// The infoTypes of information to match.
	infoTypes := []*dlppb.InfoType{
		{
			Name: "PERSON_NAME",
		},
		{
			Name: "US_STATE",
		},
	}

	// Construct item to inspect.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: input,
		},
	}

	// Construct request.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     infoTypes,
			MinLikelihood: minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		Item: item,
	}

	// Run request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	findings := resp.GetResult().GetFindings()
	if len(findings) == 0 {
		fmt.Println("No findings.")
	}
	fmt.Println("Findings:")
	for _, f := range findings {
		if includeQuote {
			fmt.Println("\tQuote: ", f.GetQuote())
		}
		fmt.Println("\tInfo type: ", f.GetInfoType().GetName())
		fmt.Println("\tLikelihood: ", f.GetLikelihood())
	}
}

// [END dlp_quickstart]
