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
)

// testLocation returns the location for testing from the environment variable.
// Skips the test if the environment variable is not set.
func testLocation(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

// testClient creates and returns a new ModelArmor client and context.
// It uses a region-specific endpoint based on the environment variable.
func testClient(t *testing.T) (*modelarmor.Client, context.Context) {
	t.Helper()

	ctx := context.Background()

	locationId := testLocation(t)

	// Create options for Model Armor client.
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationId))
	//Endpoint to send the request to regional server
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
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
		return nil, fmt.Errorf("failed to create template: %v", err)
	}

	return response, err
}

// testCleanupTemplate deletes a ModelArmor template created during a test.
// Ignores the error if the template is already deleted.
func testCleanupTemplate(t *testing.T, templateName string) {
	t.Helper()

	client, ctx := testClient(t)
	if err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
		}
	}

}

// TestUpdateTemplateWithMaskConfiguration verifies that a Model Armor template
// can be updated with a mask configuration. It creates a test template, performs
// the update, and checks the output for confirmation.
func TestUpdateTemplateWithMaskConfiguration(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var buf bytes.Buffer
	if _, err := testModelArmorTemplate(t, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if err := updateModelArmorTemplateWithMaskConfiguration(&buf, tc.ProjectID, locationID, templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Updated Model Armor Template: "; !strings.Contains(got, want) {
		t.Errorf("updateModelArmorTemplateWithMaskConfiguration: expected %q to contain %q", got, want)
	}
}
