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

// [START dlp_inspect_string_multiple_rules]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectStringMultipleRules inspects the provided text,
// avoiding matches specified in the exclusion list.
// this function implements the both exclusion and hot word rules.
func inspectStringMultipleRules(w io.Writer, projectID, textToInspect string) error {
	// projectID := "my-project-id"
	// textToInspect := "patient: Jane Doe"

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
	ContentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: dlppb.ByteContentItem_TEXT_UTF8,
				Data: []byte(textToInspect),
			},
		},
	}

	// Construct the rule that adjust the likelihood of findings
	// within a certain proximity of hot-word("patient").
	patientRule := &dlppb.InspectionRule_HotwordRule{
		HotwordRule: &dlppb.CustomInfoType_DetectionRule_HotwordRule{
			HotwordRegex: &dlppb.CustomInfoType_Regex{
				Pattern: "patient",
			},
			Proximity: &dlppb.CustomInfoType_DetectionRule_Proximity{
				WindowBefore: 10,
			},
			LikelihoodAdjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment{
				Adjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment_FixedLikelihood{
					FixedLikelihood: dlppb.Likelihood_VERY_LIKELY,
				},
			},
		},
	}

	// Construct the rule that adjust the likelihood of findings
	// within a certain proximity of hot-word("doctor").
	doctorRule := &dlppb.InspectionRule_HotwordRule{
		HotwordRule: &dlppb.CustomInfoType_DetectionRule_HotwordRule{
			HotwordRegex: &dlppb.CustomInfoType_Regex{
				Pattern: "doctor",
			},
			Proximity: &dlppb.CustomInfoType_DetectionRule_Proximity{
				WindowBefore: 10,
			},
			LikelihoodAdjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment{
				Adjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment_FixedLikelihood{
					FixedLikelihood: dlppb.Likelihood_UNLIKELY,
				},
			},
		},
	}

	// Construct the exclusion rule that avoids the matches for the specified word list.
	quasimodoRule := &dlppb.ExclusionRule{
		Type: &dlppb.ExclusionRule_Dictionary{
			Dictionary: &dlppb.CustomInfoType_Dictionary{
				Source: &dlppb.CustomInfoType_Dictionary_WordList_{
					WordList: &dlppb.CustomInfoType_Dictionary_WordList{
						Words: []string{"Quasimodo"},
					},
				},
			},
		},
		MatchingType: dlppb.MatchingType_MATCHING_TYPE_PARTIAL_MATCH,
	}

	// Construct the exclusion rule that avoids the matches for the specified regex pattern.
	redactedRule := &dlppb.ExclusionRule{
		Type: &dlppb.ExclusionRule_Regex{
			Regex: &dlppb.CustomInfoType_Regex{
				Pattern: "REDACTED",
			},
		},
		MatchingType: dlppb.MatchingType_MATCHING_TYPE_PARTIAL_MATCH,
	}

	// Construct a ruleSet that applies the rules to the PERSON_NAME infoType.
	ruleSet := &dlppb.InspectionRuleSet{
		InfoTypes: []*dlppb.InfoType{
			{Name: "PERSON_NAME"},
		},
		Rules: []*dlppb.InspectionRule{
			{Type: &dlppb.InspectionRule_HotwordRule{HotwordRule: doctorRule.HotwordRule}},
			{Type: &dlppb.InspectionRule_HotwordRule{HotwordRule: patientRule.HotwordRule}},
			{Type: &dlppb.InspectionRule_ExclusionRule{ExclusionRule: quasimodoRule}},
			{Type: &dlppb.InspectionRule_ExclusionRule{ExclusionRule: redactedRule}},
		},
	}

	// Construct the configuration for the Inspect request, including the ruleSet.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			{Name: "PERSON_NAME"},
		},
		IncludeQuote: true,
		RuleSet: []*dlppb.InspectionRuleSet{
			ruleSet,
		},
	}

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent:        fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:          ContentItem,
		InspectConfig: inspectConfig,
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

// [END dlp_inspect_string_multiple_rules]
