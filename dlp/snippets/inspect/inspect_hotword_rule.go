// Copyright 2023 Google LLC
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

// [START dlp_inspect_hotword_rule]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectWithHotWordRules inspects data with hot word rule, it uses custom
// regex with a hot word rule to increase the likelihood match
func inspectWithHotWordRules(w io.Writer, projectID, textToInspect string) error {
	// projectID := "my-project-id"
	// textToInspect := "Patient's MRN 444-5-22222 and just a number 333-2-33333"

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

	// Construct the custom regex detectors
	customInfoType := &dlppb.CustomInfoType{
		InfoType: &dlppb.InfoType{
			Name: "C_MRN",
		},
		Type: &dlppb.CustomInfoType_Regex_{
			Regex: &dlppb.CustomInfoType_Regex{
				Pattern: "[1-9]{3}-[1-9]{1}-[1-9]{5}",
			},
		},
		Likelihood: dlppb.Likelihood_POSSIBLE,
	}

	// Construct hotword rule.
	hotWordRule := &dlppb.CustomInfoType_DetectionRule_HotwordRule{
		HotwordRegex: &dlppb.CustomInfoType_Regex{
			Pattern: "(?i)(mrn|medical)(?-i)",
		},
		Proximity: &dlppb.CustomInfoType_DetectionRule_Proximity{
			WindowBefore: int32(10),
		},
		LikelihoodAdjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment{
			Adjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment_FixedLikelihood{
				FixedLikelihood: dlppb.Likelihood_VERY_LIKELY,
			},
		},
	}

	inspectionRuleSet := &dlppb.InspectionRuleSet{
		Rules: []*dlppb.InspectionRule{
			{
				Type: &dlppb.InspectionRule_HotwordRule{
					HotwordRule: hotWordRule,
				},
			},
		},
		InfoTypes: []*dlppb.InfoType{
			customInfoType.InfoType,
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:   contentItem,
		InspectConfig: &dlppb.InspectConfig{
			CustomInfoTypes: []*dlppb.CustomInfoType{
				customInfoType,
			},
			RuleSet: []*dlppb.InspectionRuleSet{
				inspectionRuleSet,
			},
			IncludeQuote: true,
		},
	}

	// Send the request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return err
	}

	// Parse the response and process results
	fmt.Fprintf(w, "Findings: %v\n", len(resp.Result.Findings))
	for _, v := range resp.GetResult().Findings {
		fmt.Fprintf(w, "Quote: %v\n", v.GetQuote())
		fmt.Fprintf(w, "InfoType Name: %v\n", v.GetInfoType().GetName())
		fmt.Fprintf(w, "Likelihood: %v\n", v.GetLikelihood())
	}
	return nil
}

// [END dlp_inspect_hotword_rule]
