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

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTemplateSamples(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	fullID := "projects/" + tc.ProjectID + "/locations/global/inspectTemplates/golang-samples-test-template"
	// Delete template before trying to create it since the test uses the same name every time.
	if err := listInspectTemplates(buf, tc.ProjectID); err != nil {
		t.Errorf("listInspectTemplates: %v", err)
	}

	if got := buf.String(); strings.Contains(got, fullID) {
		buf.Reset()
		if err := deleteInspectTemplate(buf, fullID); err != nil {
			t.Errorf("deleteInspectTemplate: %v", err)
		}
		if got, want := buf.String(), "Successfully deleted inspect template"; !strings.Contains(got, want) {
			t.Errorf("deleteInspectTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
		}
	}

	buf.Reset()

	buf.Reset()

}

func createInspectTemplateForTest(t *testing.T, projectID string, templateID, displayName, description string, infoTypeNames []string) error {
	t.Helper()
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
	log.Printf("Successfully created inspect template: %v", resp.GetName())
	return nil
}

func cleeanUpTemplates(t *testing.T, projectID, templateID string) error {
	t.Helper()
	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Delete template
	req := &dlppb.DeleteInspectTemplateRequest{
		Name: templateID,
	}

	if err := client.DeleteInspectTemplate(ctx, req); err != nil {
		return fmt.Errorf("DeleteInspectTemplate: %w", err)
	}
	log.Printf("Successfully deleted inspect template %v", templateID)

	return nil
}
