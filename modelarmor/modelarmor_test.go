// Copyright 2025 Google LLC
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

package modelarmor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	// modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
)

func testLocation(t *testing.T) string {
	t.Helper()

	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.fatal(err)
	}

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

func testClient(t *testing.T) (*modelarmor.Client, context.Context) {
	t.Helper()

	ctx := context.Background()

	locationId := testLocation(t)

	//Endpoint to send the request to regional server
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationId)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client, ctx
}

func testTemplate(t *testing.T) *modelarmorpb.Template {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	template, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID)
	if err != nil {
		t.Fatal(err)
	}
	return template
}

func testCleanupTemplate(t *testing.T, templateName string) {
	t.Helper()

	client, ctx := testClient(t)
	if err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
		}
	}

}

func testSDPTemplate(t *testing.T, projectID string, locationID string) (string, string) {
	inspectTemplateID := fmt.Sprintf("model-armour-inspect-template-%s", uuid.New().String())
	deidentifyTemplateID := fmt.Sprintf("model-armour-deidentify-template-%s", uuid.New().String())
	apiEndpoint := fmt.Sprintf("dlp.%s.rep.googleapis.com:443", locationID)
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)
	infoTypes := []*dlppb.InfoType{
		{Name: "EMAIL_ADDRESS"},
		{Name: "PHONE_NUMBER"},
		{Name: "US_INDIVIDUAL_TAXPAYER_IDENTIFICATION_NUMBER"},
	}

	// Create the DLP client.
	ctx := context.Background()
	fmt.Println("Before Client")
	dlpClient, err := dlp.NewClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		fmt.Println("Getting error while creating the client")
		t.Fatal(err)
	}
	fmt.Println("After Client")
	defer dlpClient.Close()

	Create the inspect template.
	inspectRequest := &dlppb.CreateInspectTemplateRequest{
		Parent:     parent,
		TemplateId: inspectTemplateID,
		InspectTemplate: &dlppb.InspectTemplate{
			InspectConfig: &dlppb.InspectConfig{
				InfoTypes: infoTypes,
			},
		},
	}
	inspectResponse, err := dlpClient.CreateInspectTemplate(ctx, inspectRequest)
	if err != nil {
		t.Fatal(err)
	}

	Create the deidentify template.
	deidentifyRequest := &dlppb.CreateDeidentifyTemplateRequest{
		Parent:     parent,
		TemplateId: deidentifyTemplateID,
		DeidentifyTemplate: &dlppb.DeidentifyTemplate{
			DeidentifyConfig: &dlppb.DeidentifyConfig{
				Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
					InfoTypeTransformations: &dlppb.InfoTypeTransformations{
						Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
							{
								InfoTypes: []*dlppb.InfoType{},
								PrimitiveTransformation: &dlppb.PrimitiveTransformation{
									Transformation: &dlppb.PrimitiveTransformation_ReplaceConfig{
										ReplaceConfig: &dlppb.ReplaceValueConfig{
											NewValue: &dlppb.Value{
												Type: &dlppb.Value_StringValue{StringValue: "REDACTED"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	deidentifyResponse, err := dlpClient.CreateDeidentifyTemplate(ctx, deidentifyRequest)
	if err != nil {
		t.Fatal(err)
	}

	inspectTemplateName, deidentifyTemplateName := inspectResponse.Name, deidentifyResponse.Name

	Clean up the templates.
	defer func() {
		time.Sleep(5 * time.Second)
		err := dlpClient.DeleteInspectTemplate(ctx, &dlppb.DeleteInspectTemplateRequest{Name: inspectResponse.Name})
		if err != nil {
			t.Errorf("failed to delete inspect template: %v", err)
		}
		err = dlpClient.DeleteDeidentifyTemplate(ctx, &dlppb.DeleteDeidentifyTemplateRequest{Name: deidentifyResponse.Name})
		if err != nil {
			t.Errorf("failed to delete deidentify template: %v", err)
		}
	}()

	return inspectTemplateName, deidentifyTemplateName
}

func TestCreateModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if got, want := b.String(), "Created template:"; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

func TestCreateModelArmorTemplateWithMetadata(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplateWithMetadata(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if got, want := b.String(), "Created Model Armor Template:"; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithMetadata: expected %q to contain %q", got, want)
	}
}

func TestCreateModelArmorTemplateWithLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplateWithLabels(&b, tc.ProjectID, "us-central1", templateID, map[string]string{"testkey": "testvalue"}); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if got, want := b.String(), "Created Template with labels: "; !strings.Contains(got, want) {
		fmt.Println("This is Got ", got)
		fmt.Println("This is Want ", want)
		t.Errorf("createModelArmorTemplateWithLabels: expected %q to contain %q", got, want)
	} else {
		template, err := getModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the labels
		if len(template.Labels) != 1 {
			t.Errorf("expected 1 label, got %d", len(template.Labels))
		}
		if template.Labels["testkey"] != "testvalue" {
			t.Errorf("expected label testkey to be testvalue, got %s", template.Labels["testkey"])
		}
	}
}

func TestCreateModelArmorTemplateWithBasicSDP(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplateWithBasicSDP(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if got, want := b.String(), "Created Template with basic SDP: "; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithBasicSDP: expected %q to contain %q", got, want)
	}
}

func TestCreateModelArmorTemplateWithAdvancedSDP(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	fmt.Println("BeforeTestSDPTemplate")
	inspectTemplateName, deideintifyTemplateName := testSDPTemplate(t, tc.ProjectID, "us-central1")
	fmt.Println("AfterTestSDPTemplate")

	var b bytes.Buffer
	if _, err := createModelArmorTemplateWithAdvancedSDP(&b, tc.ProjectID, "us-central1", templateID, inspectTemplateName, deideintifyTemplateName); err != nil {
		fmt.Println("Error is here")
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if got, want := b.String(), "Created Template with advanced SDP: "; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithAdvancedSDP: expected %q to contain %q", got, want)
	}
}

func TestDeleteModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if err := deleteModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Successfully deleted Model Armor template:"; !strings.Contains(got, want) {
		t.Errorf("deleteModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

func TestGetModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := getModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved template: "; !strings.Contains(got, want) {
		t.Errorf("getModelArmorTemplates: expected %q to contain %q", got, want)
	}
}

func TestListModelArmorTemplates(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := listModelArmorTemplates(&b, tc.ProjectID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Template: "; !strings.Contains(got, want) {
		t.Errorf("listModelArmorTemplates: expected %q to contain %q", got, want)
	}
}

func TestListModelArmorTemplatesWithFilter(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := listModelArmorTemplatesWithFilter(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Templates Found: "; !strings.Contains(got, want) {
		t.Errorf("listModelArmorTemplatesWithFilter: expected %q to contain %q", got, want)
	}
}

func TestUpdateTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := updateModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated Filter Config: "; !strings.Contains(got, want) {
		t.Errorf("updateModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

func TestUpdateTemplateMetadata(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := updateModelArmorTemplateMetadata(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated Model Armor Template Metadata: "; !strings.Contains(got, want) {
		t.Errorf("updateModelArmorTemplateMetadata: expected %q to contain %q", got, want)
	}
}

func TestUpdateTemplateLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := updateModelArmorTemplateLabels(&b, tc.ProjectID, "us-central1", templateID, map[string]string{"testkey": "testvalue"}); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated Model Armor Template Labels: "; !strings.Contains(got, want) {
		t.Errorf("updateModelArmorTemplateLabels: expected %q to contain %q", got, want)
	} else {
		template, err := getModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the labels
		if len(template.Labels) != 1 {
			t.Errorf("expected 1 label, got %d", len(template.Labels))
		}
		if template.Labels["testkey"] != "testvalue" {
			t.Errorf("expected label testkey to be testvalue, got %s", template.Labels["testkey"])
		}
	}
}

func TestUpdateTemplateWithMaskConfiguration(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := updateModelArmorTemplateWithMaskConfiguration(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated Model Armor Template: "; !strings.Contains(got, want) {
		t.Errorf("updateModelArmorTemplateWithMaskConfiguration: expected %q to contain %q", got, want)
	}
}





