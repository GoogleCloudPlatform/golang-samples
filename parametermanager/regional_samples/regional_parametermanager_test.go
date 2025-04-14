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

package regional_parametermanager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

// testName generates a unique name for testing purposes by creating a new UUID.
// It returns the UUID as a string or fails the test if UUID generation fails.
func testName(t *testing.T) string {
	t.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
}

// testLocation retrieves the location for testing purposes from the environment variable
// GOLANG_REGIONAL_SAMPLES_LOCATION. If the environment variable is not set,
// the test is skipped.
func testLocation(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_REGIONAL_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_REGIONAL_SAMPLES_LOCATION")
	}

	return v
}

// testParameter creates a parameter in the specified GCP project with the given format.
// It returns the created parameter and its ID or fails the test if parameter creation fails.
func testParameter(t *testing.T, projectID string, format parametermanagerpb.ParameterFormat) (*parametermanagerpb.Parameter, string) {
	t.Helper()

	parameterID := testName(t)
	locationId := testLocation(t)

	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationId)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationId)
	parameter, err := client.CreateParameter(ctx, &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: format,
		},
	})
	if err != nil {
		t.Fatalf("testParameter: failed to create parameter: %v", err)
	}

	return parameter, parameterID
}

// testParameterVersion creates a version of a parameter with the given payload in the specified GCP project.
// It returns the created parameter version and its ID or fails the test if parameter version creation fails.
func testParameterVersion(t *testing.T, projectID, parameterID, payload string) (*parametermanagerpb.ParameterVersion, string) {
	t.Helper()

	parameterVersionID := testName(t)
	locationId := testLocation(t)

	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationId)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/parameters/%s", projectID, locationId, parameterID)

	parameterVersion, err := client.CreateParameterVersion(ctx, &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parent,
		ParameterVersionId: parameterVersionID,
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data: []byte(payload),
			},
		},
	})
	if err != nil {
		t.Fatalf("testParameterVersion: failed to create parameter version: %v", err)
	}

	return parameterVersion, parameterVersionID
}

// testCleanupParameter deletes the specified parameter in the GCP project.
// It fails the test if the parameter deletion fails.
func testCleanupParameter(t *testing.T, name string) {
	t.Helper()

	locationId := testLocation(t)
	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationId)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	if err := client.DeleteParameter(ctx, &parametermanagerpb.DeleteParameterRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupParameter: failed to delete parameter: %v", err)
		}
	}
}

// testCleanupParameterVersion deletes the specified parameter version in the GCP project.
// It fails the test if the parameter version deletion fails.
func testCleanupParameterVersion(t *testing.T, name string) {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)

	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationId)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	if err := client.DeleteParameterVersion(ctx, &parametermanagerpb.DeleteParameterVersionRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupParameterVersion: failed to delete parameter version: %v", err)
		}
	}
}

// testSecret creates a secret in the specified GCP project.
// It returns the created secret or fails the test if secret creation fails.
func testSecret(t *testing.T, projectID string) *secretmanagerpb.Secret {
	t.Helper()

	secretID := testName(t)
	locationId := testLocation(t)

	ctx := context.Background()
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	secret, err := client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, locationId),
		SecretId: secretID,
		Secret:   &secretmanagerpb.Secret{},
	})
	if err != nil {
		t.Fatalf("testSecret: failed to create secret: %v", err)
	}

	return secret
}

// testSecretVersion creates a version of a secret with the given payload in the specified GCP project.
// It returns the created secret version or fails the test if secret version creation fails.
func testSecretVersion(t *testing.T, parent string, payload []byte) *secretmanagerpb.SecretVersion {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)

	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	version, err := client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	})
	if err != nil {
		t.Fatalf("testSecretVersion: failed to create secret version: %v", err)
	}
	return version
}

// testIamGrantAccess grants the specified member access permissions to the secret.
func testIamGrantAccess(t *testing.T, name, member string) error {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)

	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	handle := client.IAM(name)
	policy, err := handle.Policy(ctx)
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}

	// Grant the member access permissions.
	policy.Add(member, "roles/secretmanager.secretAccessor")
	if err = handle.SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	return nil
}

// testCleanupSecret deletes the specified secret in the GCP project.
// It fails the test if the secret deletion fails.
func testCleanupSecret(t *testing.T, name string) {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)

	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	if err := client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupSecret: failed to delete secret: %v", err)
		}
	}
}

// TestGetRegionalParam tests the getRegionalParam function by creating a parameter,
// then attempts to retrieve the created parameter. It verifies if the parameter
// was successfully retrieved by checking the output.
func TestGetRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	defer testCleanupParameter(t, parameter.Name)

	locationId := testLocation(t)
	var buf bytes.Buffer
	if err := getRegionalParam(&buf, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter %s with format JSON", parameter.Name); !strings.Contains(got, want) {
		t.Errorf("GetParameter: expected %q to contain %q", got, want)
	}
}

// TestGetRegionalParamVersion tests the getRegionalParamVersion function by creating a parameter and its version,
// then attempts to retrieve the created parameter version. It verifies if the parameter version
// was successfully retrieved by checking the output.
func TestGetRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)
	locationId := testLocation(t)

	var buf bytes.Buffer
	if err := getRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter version %s with disabled state in %v", parameterVersion.Name, parameterVersion.Disabled); !strings.Contains(got, want) {
		t.Errorf("GetParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestRenderRegionalParamVersion tests the renderRegionalParamVersion function by creating a parameter,
// its version, and a secret. It then attempts to render the created parameter version
// and verifies if the parameter version was successfully rendered by checking the output.
func TestRenderRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	secret := testSecret(t, tc.ProjectID)
	testSecretVersion(t, secret.Name, []byte("very secret data"))
	payload := fmt.Sprintf(`{"username": "test-user","password": "__REF__(//secretmanager.googleapis.com/%s/versions/latest)"}`, secret.Name)
	if err := testIamGrantAccess(t, secret.Name, parameter.PolicyMember.IamPolicyUidPrincipal); err != nil {
		t.Fatal(err)
	}
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	locationId := testLocation(t)

	defer testCleanupSecret(t, secret.Name)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var buf bytes.Buffer
	time.Sleep(2 * time.Minute)
	if err := renderRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Rendered regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("RenderParameterVersion: expected %q to contain %q", got, want)
	}
	expectedPayload := `{"username": "test-user","password": "very secret data"}`
	if got, want := buf.String(), fmt.Sprintf("Rendered payload: %s", expectedPayload); !strings.Contains(got, want) {
		t.Errorf("RenderParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestListRegionalParam tests the listRegionalParam function by creating multiple parameters,
// then attempts to list the created parameters. It verifies if the parameters
// were successfully listed by checking the output.
func TestListRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter1, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameter2, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	locationId := testLocation(t)

	defer testCleanupParameter(t, parameter1.Name)
	defer testCleanupParameter(t, parameter2.Name)

	var buf bytes.Buffer
	if err := listRegionalParam(&buf, tc.ProjectID, locationId); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter %s with format %s", parameter1.Name, parameter1.Format); !strings.Contains(got, want) {
		t.Errorf("ListParameter: expected %q to contain %q", got, want)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter %s with format %s", parameter2.Name, parameter2.Format); !strings.Contains(got, want) {
		t.Errorf("ListParameter: expected %q to contain %q", got, want)
	}
}

// TestListRegionalParamVersion tests the listRegionalParamVersion function by creating a parameter and its versions,
// then attempts to list the created parameter versions. It verifies if the parameter versions
// were successfully listed by checking the output.
func TestListRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion1, _ := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	parameterVersion2, _ := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	locationId := testLocation(t)

	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion1.Name)
	defer testCleanupParameterVersion(t, parameterVersion2.Name)

	var buf bytes.Buffer
	if err := listRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter version %s with disabled state in %v", parameterVersion1.Name, parameterVersion1.Disabled); !strings.Contains(got, want) {
		t.Errorf("ListParameterVersion: expected %q to contain %q", got, want)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter version %s with disabled state in %v", parameterVersion2.Name, parameterVersion2.Disabled); !strings.Contains(got, want) {
		t.Errorf("ListParameterVersion: expected %q to contain %q", got, want)
	}
}
