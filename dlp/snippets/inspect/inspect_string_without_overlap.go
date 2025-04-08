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

// [START dlp_inspect_string_without_overlap]

import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectStringWithoutOverlap inspects a string for sensitive data
// and omit overlapping matches on domain and email
func inspectStringWithoutOverlap(w io.Writer, projectID, textToInspect string) error {
	// projectID := "my-project-id"
	// textToInspect := "example.com is a domain, james@example.org is an email."

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

	// Specify the type and content to be inspected.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: dlppb.ByteContentItem_TEXT_UTF8,
				Data: []byte(textToInspect),
			},
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types.
	infoTypes := []*dlppb.InfoType{
		{Name: "DOMAIN_NAME"},
		{Name: "EMAIL_ADDRESS"},
	}

	// Define a custom info type to exclude email addresses
	customInfotype := &dlppb.CustomInfoType{
		InfoType: &dlppb.InfoType{
			Name: "EMAIL_ADDRESS",
		},
		ExclusionType: dlppb.CustomInfoType_EXCLUSION_TYPE_EXCLUDE,
	}

	// Exclude EMAIL_ADDRESS matches
	exclusionRule := &dlppb.ExclusionRule{
		Type: &dlppb.ExclusionRule_ExcludeInfoTypes{
			ExcludeInfoTypes: &dlppb.ExcludeInfoTypes{
				InfoTypes: []*dlppb.InfoType{
					{Name: "EMAIL_ADDRESS"},
				},
			},
		},
		MatchingType: dlppb.MatchingType_MATCHING_TYPE_PARTIAL_MATCH,
	}

	// Construct a ruleSet that applies the exclusion rule to the DOMAIN_NAME infoType.
	// If a DOMAIN_NAME match is part of an EMAIL_ADDRESS match, the DOMAIN_NAME match will
	// be excluded.
	ruleSet := &dlppb.InspectionRuleSet{
		InfoTypes: []*dlppb.InfoType{
			{Name: "DOMAIN_NAME"},
		},
		Rules: []*dlppb.InspectionRule{
			{
				Type: &dlppb.InspectionRule_ExclusionRule{
					ExclusionRule: exclusionRule,
				},
			},
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:   contentItem,
		// Construct the Inspect request to be sent by the client.
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: infoTypes,
			CustomInfoTypes: []*dlppb.CustomInfoType{
				customInfotype,
			},
			IncludeQuote: true,
			RuleSet: []*dlppb.InspectionRuleSet{
				ruleSet,
			},
		},
	}

	// Send the request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return err
	}

	// Process the results.
	fmt.Fprintf(w, "Findings: %v\n", len(resp.Result.Findings))
	for _, v := range resp.GetResult().Findings {
		fmt.Fprintf(w, "Quote: %v\n", v.GetQuote())
		fmt.Fprintf(w, "Infotype Name: %v\n", v.GetInfoType().GetName())
		fmt.Fprintf(w, "Likelihood: %v\n", v.GetLikelihood())
	}
	return nil
}

// [END dlp_inspect_string_without_overlap]
