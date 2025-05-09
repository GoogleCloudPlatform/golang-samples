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
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
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

// testAllFilterTemplate creates a new ModelArmor template with all filters enabled.
// It returns the template ID and filter config for cleanup.
func testAllFilterTemplate(t *testing.T, templateID string) (*modelarmorpb.Template, *modelarmorpb.FilterConfig, error) {
	t.Helper()
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	client, ctx := testClient(t)

	filterConfig := &modelarmorpb.FilterConfig{
		RaiSettings: &modelarmorpb.RaiFilterSettings{
			RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
				{
					FilterType:      modelarmorpb.RaiFilterType_DANGEROUS,
					ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
				},
				{
					FilterType:      modelarmorpb.RaiFilterType_HARASSMENT,
					ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
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
		PiAndJailbreakFilterSettings: &modelarmorpb.PiAndJailbreakFilterSettings{
			FilterEnforcement: modelarmorpb.PiAndJailbreakFilterSettings_ENABLED,
			ConfidenceLevel:   modelarmorpb.DetectionConfidenceLevel_MEDIUM_AND_ABOVE,
		},
		MaliciousUriFilterSettings: &modelarmorpb.MaliciousUriFilterSettings{
			FilterEnforcement: modelarmorpb.MaliciousUriFilterSettings_ENABLED,
		},
	}

	template := &modelarmorpb.Template{
		FilterConfig: filterConfig,
	}

	req := &modelarmorpb.CreateTemplateRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, locationID),
		TemplateId: templateID,
		Template:   template,
	}

	// When creating the client or making the call
	retryOpts := []gax.CallOption{
		gax.WithRetry(func() gax.Retryer {
			return gax.OnCodes([]codes.Code{
				codes.Unavailable,
				codes.DeadlineExceeded,
			}, gax.Backoff{
				Initial:    1 * time.Second,
				Max:        30 * time.Second,
				Multiplier: 2,
			})
		}),
	}

	// Using retry mechanism similar to retry_ma_create_template
	var response *modelarmorpb.Template
	var err error

	response, err = client.CreateTemplate(ctx, req, retryOpts...)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create template: %w", err)
	}

	return response, filterConfig, nil
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

// testBasicSDPTemplate creates Model Armor template with basic SDP configuration.
func testBasicSDPTemplate(t *testing.T, templateID string) (*modelarmorpb.Template, error) {
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
			SdpSettings: &modelarmorpb.SdpFilterSettings{
				SdpConfiguration: &modelarmorpb.SdpFilterSettings_BasicConfig{
					BasicConfig: &modelarmorpb.SdpBasicConfig{
						FilterEnforcement: modelarmorpb.SdpBasicConfig_ENABLED,
					},
				},
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
	err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupTemplate: failed to delete template %v", err)
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

// TestSanitizeModelResponse verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response.
func TestSanitizeModelResponse(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	modelResponse := "Unsanitized model output"
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var buf bytes.Buffer
	if _, err := testModelArmorTemplate(t, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if err := sanitizeModelResponse(&buf, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Sanitization Result: "; !strings.Contains(got, want) {
		t.Errorf("sanitizeModelResponse: expected %q to contain %q", got, want)
	}
}

// TestSanitizeModelResponseWithAllRaiFilterTemplate verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response with all rai filters template.
func TestSanitizeModelResponseWithAllRaiFilterTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)

	// Generate a unique template ID
	templateID := fmt.Sprintf("test-all-rai-filter-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer

	// Create template with all filters
	_, _, err := testAllFilterTemplate(t, templateID)
	if err != nil {
		t.Fatalf("Failed to create template with all filters: %v", err)
	}
	defer testCleanupTemplate(t, templateName)

	// Define a benign model response
	modelResponse := "To make cheesecake without oven, you'll need to follow these steps...."

	// Call sanitizeModelResponse with buffer
	if err := sanitizeModelResponse(&b, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	// Check that the output contains the expected sanitization result
	output := b.String()

	// Check for overall NO_MATCH_FOUND
	if !strings.Contains(output, "NO_MATCH_FOUND") {
		t.Errorf("Expected output to indicate NO_MATCH_FOUND for overall result, got: %q", output)
	}
}

// TestSanitizeModelResponseWithMaliciousUrlTemplate verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response with all malicious URL filters template.
func TestSanitizeModelResponseWithMaliciousUrlTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)

	// Generate a unique template ID
	templateID := fmt.Sprintf("test-malicious-url-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer

	// Create template with all filters
	_, _, err := testAllFilterTemplate(t, templateID)
	if err != nil {
		t.Fatalf("Failed to create template with all filters: %v", err)
	}
	defer testCleanupTemplate(t, templateName)

	// Define a model response with a malicious URL
	modelResponse := "You can use this to make a cake: https://testsafebrowsing.appspot.com/s/malware.html"

	// Call sanitizeModelResponse with buffer
	if err := sanitizeModelResponse(&b, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	// Check that the output contains the expected sanitization result
	output := b.String()

	// Check for overall MATCH_FOUND
	if !strings.Contains(output, "MATCH_FOUND") {
		t.Errorf("Expected output to indicate MATCH_FOUND for overall result, got: %q", output)
	}

}

// TestSanitizeModelResponseWithCsamTemplate verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response with all csam filters template.
func TestSanitizeModelResponseWithCsamTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)

	// Generate a unique template ID
	templateID := fmt.Sprintf("test-csam-filter-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer

	// Create template with all filters
	_, _, err := testAllFilterTemplate(t, templateID)
	if err != nil {
		t.Fatalf("Failed to create template with all filters: %v", err)
	}
	defer testCleanupTemplate(t, templateName)

	// Define a benign model response related to children
	modelResponse := "Here is how to teach long division to a child"

	// Call sanitizeModelResponse with buffer
	if err := sanitizeModelResponse(&b, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	// Check that the output contains the expected sanitization result
	output := b.String()

	// Check for overall NO_MATCH_FOUND
	if !strings.Contains(output, "NO_MATCH_FOUND") {
		t.Errorf("Expected output to indicate NO_MATCH_FOUND for overall result, got: %q", output)
	}
}

// TestSanitizeModelResponseWithEmptyTemplate verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response with all empty filters template.
func TestSanitizeModelResponseWithEmptyTemplate(t *testing.T) {
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

	// Define model response with sensitive information
	// (which won't be flagged because the template is empty)
	modelResponse := "For following email 1l6Y2@example.com found following associated phone number: 954-321-7890 and this ITIN: 988-86-1234"

	// Call sanitizeModelResponse with buffer
	if err := sanitizeModelResponse(&b, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	// Check buffer output
	output := b.String()

	// Check for NO_MATCH_FOUND
	if !strings.Contains(output, "NO_MATCH_FOUND") {
		t.Errorf("Expected output to indicate NO_MATCH_FOUND for overall result, got: %q", output)
	}
}

// TestSanitizeModelResponseWithBasicSdpTemplate verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response with all basic SDP filters template.
func TestSanitizeModelResponseWithBasicSdpTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)

	// Generate a unique template ID
	templateID := fmt.Sprintf("test-basic-sdp-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer

	// Create template with basic SDP configuration (inspection only, no deidentification)
	_, err := testBasicSDPTemplate(t, templateID)
	if err != nil {
		t.Fatalf("Failed to create template with basic SDP: %v", err)
	}
	defer testCleanupTemplate(t, templateName)

	// Define model response with sensitive information
	modelResponse := "For following email 1l6Y2@example.com found following associated phone number: 954-321-7890 and this ITIN: 988-86-1234"

	// Call sanitizeModelResponse with buffer
	if err := sanitizeModelResponse(&b, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	// Check buffer output
	output := b.String()

	// Check that the findings include the expected info types
	infoTypes := []string{"US_INDIVIDUAL_TAXPAYER_IDENTIFICATION_NUMBER"}
	for _, infoType := range infoTypes {
		if !strings.Contains(output, infoType) {
			t.Errorf("Expected output to contain finding for %s, but it wasn't found: %q", infoType, output)
		}
	}
}

// TestSanitizeModelResponseWithAdvanceSdpTemplate verifies that the sanitizeModelResponse function
// returns a properly formatted sanitization result for a given model response with all advanced SDP filters template.
func TestSanitizeModelResponseWithAdvanceSdpTemplate(t *testing.T) {
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

	// Define model response with sensitive information
	modelResponse := "For following email 1l6Y2@example.com found following associated phone number: 954-321-7890 and this ITIN: 988-86-1234"

	// Call sanitizeModelResponse with buffer
	if err := sanitizeModelResponse(&b, tc.ProjectID, locationID, templateID, modelResponse); err != nil {
		t.Fatal(err)
	}

	// Check buffer output
	output := b.String()

	// Check that sensitive information is redacted in the output
	if strings.Contains(output, "[REDACTED]") {
		t.Errorf("Expected email to be redacted in the output, but it was found: %q", output)
	}
}
