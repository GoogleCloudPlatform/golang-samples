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

	return response, err
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

// TestSanitizeUserPrompt verifies that the sanitizeUserPrompt function
// It ensures the output contains the expected sanitization result.
func TestSanitizeUserPrompt(t *testing.T) {

	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	userPrompt := "how can I teach my child to brush their teeth?"
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID)
	var b bytes.Buffer
	if _, err := testModelArmorTemplate(t, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if err := sanitizeUserPrompt(&b, tc.ProjectID, locationID, templateID, userPrompt); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Sanitization Result: "; !strings.Contains(got, want) {
		t.Errorf("sanitizeUserPrompt: expected %q to contain %q", got, want)
	}
}
