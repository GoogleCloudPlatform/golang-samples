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

package template

// [START dlp_create_inspect_template]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// createInspectTemplate creates a template with the given configuration.
func createInspectTemplate(w io.Writer, projectID string, templateID, displayName, description string, infoTypeNames []string) error {
	// projectID := "my-project-id"
	// templateID := "my-template"
	// displayName := "My Template"
	// description := "My template description"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}

	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Convert the info type strings to a list of InfoTypes.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	// Create a configured request.
	req := &dlppb.CreateInspectTemplateRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/global", projectID),
		TemplateId: templateID,
		InspectTemplate: &dlppb.InspectTemplate{
			DisplayName: displayName,
			Description: description,
			InspectConfig: &dlppb.InspectConfig{
				InfoTypes:     infoTypes,
				MinLikelihood: dlppb.Likelihood_POSSIBLE,
				Limits: &dlppb.InspectConfig_FindingLimits{
					MaxFindingsPerRequest: 10,
				},
			},
		},
	}
	// Send the request.
	resp, err := client.CreateInspectTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateInspectTemplate: %w", err)
	}
	// Print the result.
	fmt.Fprintf(w, "Successfully created inspect template: %v", resp.GetName())
	return nil
}

// [END dlp_create_inspect_template]
