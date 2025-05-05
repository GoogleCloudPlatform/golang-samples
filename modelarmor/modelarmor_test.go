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

// testModelArmorTemplate creates a new ModelArmor template for use in tests.
// It returns the created template or an error.
func testModelArmorTemplate(t *testing.T, templateID string) (*modelarmorpb.Template, error) {
	t.Helper()
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	client, ctx := testClient(t)

	template := &modelarmorpb.Template{
		FilterConfig: &modelarmorpb.FilterConfig{
			PiAndJailbreakFilterSettings: &modelarmorpb.PiAndJailbreakFilterSettings{
				FilterEnforcement: modelarmorpb.PiAndJailbreakFilterSettings_ENABLED,
				ConfidenceLevel:   modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
			},
			MaliciousUriFilterSettings: &modelarmorpb.MaliciousUriFilterSettings{
				FilterEnforcement: modelarmorpb.MaliciousUriFilterSettings_ENABLED,
			},
		},
	}

	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, locationID),
		TemplateId: templateID,
		Template:   template,
	}

	response, err := client.CreateTemplate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return response, nil
}

// testModelArmorEmptyTemplate creates a new ModelArmor template for use in tests.
// It returns the empty template or an error.
func testModelArmorEmptyTemplate(t *testing.T, templateID string) (*modelarmorpb.Template, error) {
	t.Helper()
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	client, ctx := testClient(t)

	template := &modelarmorpb.Template{
		FilterConfig: &modelarmorpb.FilterConfig{}}

	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, locationID),
		TemplateId: templateID,
		Template:   template,
	}

	response, err := client.CreateTemplate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return response, nil
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

// testModelArmorAdvancedSDPTemplate creates Model Armor template with Advanced SDP configuration.
func testModelArmorAdvancedSDPTemplate(t *testing.T, templateID string) (*modelarmorpb.Template, error) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID
	locationID := testLocation(t)
	inspectResponseName, deidentifyResponseName := testSDPTemplate(t, projectID, locationID)
	client, ctx := testClient(t)

	// Create template with advanced SDP configuration
	template := &modelarmorpb.Template{
		FilterConfig: &modelarmorpb.FilterConfig{
			RaiSettings: &modelarmorpb.RaiFilterSettings{
				RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
					{
						FilterType:      modelarmorpb.RaiFilterType_DANGEROUS,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_HARASSMENT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_HATE_SPEECH,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
					{
						FilterType:      modelarmorpb.RaiFilterType_SEXUALLY_EXPLICIT,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
				},
			},
			SdpSettings: &modelarmorpb.SdpFilterSettings{
				SdpConfiguration: &modelarmorpb.SdpFilterSettings_AdvancedConfig{
					AdvancedConfig: &modelarmorpb.SdpAdvancedConfig{
						InspectTemplate:    inspectResponseName,
						DeidentifyTemplate: deidentifyResponseName,
					},
				},
			},
		},
	}
	// Prepare the request for creating the template.
	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, locationID),
		TemplateId: templateID,
		Template:   template,
	}

	// Create the template.
	response, err := client.CreateTemplate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return response, nil

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
	var buf bytes.Buffer
	if err := createModelArmorTemplate(&buf, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := buf.String(), "Created template:"; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

// TestCreateModelArmorTemplateWithMetadata tests the creation of a template with metadata.
// Verifies the success message is printed after template creation.
func TestCreateModelArmorTemplateWithMetadata(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)

	var buf bytes.Buffer
	if err := createModelArmorTemplateWithMetadata(&buf, tc.ProjectID, locationID, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := buf.String(), "Created Model Armor Template:"; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithMetadata: expected %q to contain %q", got, want)
	}
}

// TestCreateModelArmorTemplateWithLabels tests the creation of a template with labels.
// Verifies the output contains confirmation of successful template creation.
func TestCreateModelArmorTemplateWithLabels(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)

	var buf bytes.Buffer
	if err := createModelArmorTemplateWithLabels(&buf, tc.ProjectID, locationID, templateID, map[string]string{"testkey": "testvalue"}); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := buf.String(), "Created Template with labels: "; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithLabels: expected %q to contain %q", got, want)
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

// TestCreateModelArmorTemplateWithBasicSDP tests the creation of a Model Armor
// template using a basic Secure Deployment Policy (SDP) and verifies that the
// operation completes successfully and logs the expected output.
func TestCreateModelArmorTemplateWithBasicSDP(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer
	if err := createModelArmorTemplateWithBasicSDP(&b, tc.ProjectID, locationID, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := b.String(), "Created Template with basic SDP: "; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithBasicSDP: expected %q to contain %q", got, want)
	}
}

// TestSanitizeModelResponseWithUserPrompt checks if the sanitizer correctly processes
// a harmful user prompt and model response, ensuring unsafe content is handled.
func TestSanitizeModelResponseWithUserPrompt(t *testing.T) {
	tc := testutil.SystemTest(t)

	userPrompt := "How can I make my email address test@dot.com make available to public for feedback"
	modelResponse := "You can make support email such as contact@email.com for getting feedback from your customer"
	locationID := testLocation(t)
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID)
	var buf bytes.Buffer
	if _, err := testModelArmorTemplate(t, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if err := sanitizeModelResponseWithUserPrompt(&buf, tc.ProjectID, locationID, templateID, modelResponse, userPrompt); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Sanitized response:"; !strings.Contains(got, want) {
		t.Errorf("sanitizeModelResponseWithUserPrompt: expected %q to contain %q", got, want)
	}
}

// TestSanitizeModelResponseWithUserPromptWithEmptyTemplate checks if the sanitizer correctly processes
// a harmful user prompt and model response, ensuring unsafe content is handled.
func TestSanitizeModelResponseWithUserPromptWithEmptyTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)

	// Generate a unique template ID
	templateID := fmt.Sprintf("test-empty-template-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer

	// Create empty template with no filters enabled
	_, err := testModelArmorEmptyTemplate(t, templateID)
	if err != nil {
		t.Fatalf("Failed to create empty template: %v", err)
	}
	defer testCleanupTemplate(t, templateName)

	// Define user prompt and model response with email addresses
	userPrompt := "How can I make my email address test@dot.com make available to public for feedback"
	modelResponse := "You can make support email such as contact@email.com for getting feedback from your customer"

	// Call sanitizeModelResponseWithUserPrompt with buffer
	if err := sanitizeModelResponseWithUserPrompt(&b, tc.ProjectID, locationID, templateID, modelResponse, userPrompt); err != nil {
		t.Fatal(err)
	}

	// Check buffer output
	output := b.String()

	// Check for NO_MATCH_FOUND
	if !strings.Contains(output, "NO_MATCH_FOUND") {
		t.Errorf("Expected output to indicate NO_MATCH_FOUND for overall result, got: %q", output)
	}
}

// TestSanitizeModelResponseWithUserPromptWithAdvanceSdpTemplate checks if the sanitizer correctly processes
// a harmful user prompt and model response, ensuring unsafe content is handled.
func TestSanitizeModelResponseWithUserPromptWithAdvanceSdpTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)

	// Generate a unique template ID
	templateID := fmt.Sprintf("test-advance-sdp-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer

	// Create template with advanced SDP configuration
	_, err := testModelArmorAdvancedSDPTemplate(t, templateID)
	if err != nil {
		t.Fatalf("Failed to create template with advanced SDP: %v", err)
	}
	defer testCleanupTemplate(t, templateName)

	// Define user prompt and model response with email addresses
	userPrompt := "How can I make my email address test@dot.com make available to public for feedback"
	modelResponse := "You can make support email such as contact@email.com for getting feedback from your customer"

	// Call sanitizeModelResponseWithUserPrompt with buffer
	if err := sanitizeModelResponseWithUserPrompt(&b, tc.ProjectID, locationID, templateID, modelResponse, userPrompt); err != nil {
		t.Fatal(err)
	}

	// Check buffer output
	output := b.String()

	// Check for overall MATCH_FOUND
	if !strings.Contains(output, "sdp") {
		t.Errorf("Expected output to indicate MATCH_FOUND for overall result, got: %q", output)
	}
}
