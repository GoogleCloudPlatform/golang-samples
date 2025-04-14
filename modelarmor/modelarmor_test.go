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

// testLocation returns the region where the tests should run,
// based on the GOLANG_SAMPLES_LOCATION environment variable.
// Skips the test if the variable is not set.
func testLocation(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testLocation: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

// testClient creates a new Model Armor API client using the regional endpoint
// derived from the environment. It returns the client and the context.
func testClient(t *testing.T) (*modelarmor.Client, context.Context) {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)


	// Set the endpoint for the regional replica based on location
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationId))
	client, err := modelarmor.NewClient(ctx,opts)
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
		return nil, fmt.Errorf("failed to create template: %v", err)
	}

	return response, err
}

// testCleanupTemplate deletes a given template by name if it exists.
// Ignores NotFound errors, which means the resource has already been deleted.
func testCleanupTemplate(t *testing.T, templateName string) {
	t.Helper()

	client, ctx := testClient(t)
	if err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
		}
	}
}

// TestCreateModelArmorTemplate tests the creation of a basic Model Armor template.
// Verifies that the output includes a success message after creation.
func TestCreateModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, locationID, templateID)
	var b bytes.Buffer
	if _, err := testModelArmorTemplate(t, templateID); err != nil {
		t.Fatal(err)
	}

	defer testCleanupTemplate(t, templateName)

	if got, want := b.String(), "Created template:"; !strings.Contains(got, want) {
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

	var b bytes.Buffer
	if err := createModelArmorTemplateWithMetadata(&b, tc.ProjectID, locationID, templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := b.String(), "Created Model Armor Template:"; !strings.Contains(got, want) {
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

	var b bytes.Buffer
	if err := createModelArmorTemplateWithLabels(&b, tc.ProjectID, locationID, templateID, map[string]string{"testkey": "testvalue"}); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := b.String(), "Created Template with labels: "; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplateWithLabels: expected %q to contain %q", got, want)
	}
}
