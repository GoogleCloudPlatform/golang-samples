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
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

// testLocation retrieves the GOLANG_SAMPLES_LOCATION environment variable
// used to determine the region for running the test.
// Skips the test if the environment variable is not set.
func testLocation(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testLocation: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

// testClient initializes and returns a new Model Armor API client and context
// targeting the endpoint based on the specified location.
func testClient(t *testing.T) (*modelarmor.Client, context.Context) {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)
	// Create option for Model Armor client.
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationId))
	// Create a new client using the regional endpoint
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
	}

	return client, ctx
}

// testCleanupTemplate deletes the specified Model Armor template if it exists,
// ignoring the error if the template is already deleted.
func testCleanupTemplate(t *testing.T, templateName string) {
	t.Helper()

	client, ctx := testClient(t)
	if err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName}); err != nil {
		// Ignore NotFound errors (template may already be deleted)
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
		}
	}
}

// testSDPTemplate creates DLP inspect and deidentify templates for use in tests.
func testSDPTemplate(t *testing.T, projectID string, locationID string) (string, string) {
	inspectTemplateID := fmt.Sprintf("model-armor-inspect-template-%s", uuid.New().String())
	deidentifyTemplateID := fmt.Sprintf("model-armor-deidentify-template-%s", uuid.New().String())
	apiEndpoint := fmt.Sprintf("dlp.%s.rep.googleapis.com:443", locationID)
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	infoTypes := []*dlppb.InfoType{
		{Name: "EMAIL_ADDRESS"},
		{Name: "PHONE_NUMBER"},
		{Name: "US_INDIVIDUAL_TAXPAYER_IDENTIFICATION_NUMBER"},
	}

	ctx := context.Background()
	dlpClient, err := dlp.NewClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		t.Fatalf("Getting error while creating the client: %v", err)
	}
	defer dlpClient.Close()

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

	// Cleanup the templates after test.
	defer func() {
		time.Sleep(5 * time.Second)
		err := dlpClient.DeleteInspectTemplate(ctx, &dlppb.DeleteInspectTemplateRequest{Name: inspectResponse.Name})
		if err != nil {
			t.Errorf("failed to delete inspect template: %v, err: %v", inspectResponse.Name, err)
		}
		err = dlpClient.DeleteDeidentifyTemplate(ctx, &dlppb.DeleteDeidentifyTemplateRequest{Name: deidentifyResponse.Name})
		if err != nil {
			t.Errorf("failed to delete deidentify template: %v, err: %v", deidentifyResponse.Name, err)
		}
	}()

	return inspectResponse.Name, deidentifyResponse.Name
}

// TestCreateModelArmorTemplateWithAdvancedSDP tests creating a
// Model Armor template with advanced SDP using DLP templates.
func TestCreateModelArmorTemplateWithAdvancedSDP(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	inspectTemplateName, deideintifyTemplateName := testSDPTemplate(t, tc.ProjectID, testLocation(t))
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, testLocation(t), templateID)
	var buf bytes.Buffer
	if err := createModelArmorTemplateWithAdvancedSDP(&buf, tc.ProjectID, testLocation(t), templateID, inspectTemplateName, deideintifyTemplateName); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := buf.String(), "Created Template with advanced SDP: "; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithAdvancedSDP: expected %q to contain %q", got, want)
	}
}

// TestCreateModelArmorTemplate verifies the creation of a Model Armor template.
// It ensures the output contains a confirmation message after creation.
func TestCreateModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, testLocation(t), templateID)
	var b bytes.Buffer
	if err := createModelArmorTemplate(&b, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := b.String(), "Created template:"; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

// TestDeleteModelArmorTemplate verifies the deletion of a Model Armor template.
// It ensures the output contains a confirmation message after deletion.
func TestDeleteModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var buf bytes.Buffer
	// Create template first to ensure it exists for deletion
	if err := createModelArmorTemplate(&buf, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}

	// Attempt to delete the template
	if err := deleteModelArmorTemplate(&buf, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Successfully deleted Model Armor template:"; !strings.Contains(got, want) {
		t.Errorf("deleteModelArmorTemplate: expected %q to contain %q", got, want)
	}
}
