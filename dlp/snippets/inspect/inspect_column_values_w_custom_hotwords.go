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

// [START dlp_inspect_column_values_w_custom_hotwords]

import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectTableWithCustomHotword Sets the match likelihood of a table column to customize data inspection results.
// This example showcases how you can adjust the match likelihood for an entire column of data, enabling the
// exclusion of specific columns from inspection if needed.
func inspectTableWithCustomHotword(w io.Writer, projectID, hotwordRegexPattern string) error {
	// projectID := "your-project-id"
	// hotwordRegexPattern := "(Fake Social Security Number)"

	tableToInspect := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "Fake Social Security Number"},
			{Name: "Real Social Security Number"},
		},
		Rows: []*dlppb.Table_Row{
			{
				Values: []*dlppb.Value{
					{
						Type: &dlppb.Value_StringValue{StringValue: "111-11-1111"},
					},
					{
						Type: &dlppb.Value_StringValue{StringValue: "222-22-2222"},
					},
				},
			},
		},
	}

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

	// Specify what content you want the service to de-identify.
	contentItem := &dlppb.ContentItem_Table{
		Table: tableToInspect,
	}

	// Specify the likelihood adjustment to adjust the match likelihood for your detection rule
	// based on your needs and desired level of sensitivity in data analysis.
	likelihoodAdjustment := &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment{
		Adjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment_FixedLikelihood{
			FixedLikelihood: dlppb.Likelihood_VERY_UNLIKELY,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types.
	infoTypes := []*dlppb.InfoType{
		{Name: "US_SOCIAL_SECURITY_NUMBER"},
	}

	// Specify the proximity so that It helps identify sensitive information
	// occurring near other data points, enabling more context-aware analysis.
	proximity := &dlppb.CustomInfoType_DetectionRule_Proximity{
		WindowBefore: 5,
	}

	// Construct hotWord rule.
	hotwordRule := &dlppb.CustomInfoType_DetectionRule_HotwordRule{
		HotwordRegex: &dlppb.CustomInfoType_Regex{
			Pattern: hotwordRegexPattern,
		},
		Proximity:            proximity,
		LikelihoodAdjustment: likelihoodAdjustment,
	}

	// Construct rule set for the inspect config.
	inspectionRuleSet := &dlppb.InspectionRuleSet{
		InfoTypes: infoTypes,
		Rules: []*dlppb.InspectionRule{
			{
				Type: &dlppb.InspectionRule_HotwordRule{
					HotwordRule: hotwordRule,
				},
			},
		},
	}

	// Construct the configuration for the Inspect request.
	config := &dlppb.InspectConfig{
		IncludeQuote:  true,
		InfoTypes:     infoTypes,
		MinLikelihood: dlppb.Likelihood_POSSIBLE,
		RuleSet: []*dlppb.InspectionRuleSet{
			inspectionRuleSet,
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item: &dlppb.ContentItem{
			DataItem: contentItem,
		},
		InspectConfig: config,
	}
	// Send the request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return err
	}

	// Parse the response and process results.
	fmt.Fprintf(w, "Findings: %v\n", len(resp.Result.Findings))
	for _, v := range resp.GetResult().Findings {
		fmt.Fprintf(w, "Quote: %v\n", v.GetQuote())
		fmt.Fprintf(w, "Infotype Name: %v\n", v.GetInfoType().GetName())
		fmt.Fprintf(w, "Likelihood: %v\n", v.GetLikelihood())
	}
	return nil
}

// [END dlp_inspect_column_values_w_custom_hotwords]
