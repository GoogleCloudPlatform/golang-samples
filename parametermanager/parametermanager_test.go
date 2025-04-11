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

package parametermanager

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
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

// testParameter creates a parameter in the specified GCP project with the given format.
// It returns the created parameter and its ID or fails the test if parameter creation fails.
func testParameter(t *testing.T, projectID string, format parametermanagerpb.ParameterFormat) (*parametermanagerpb.Parameter, string) {
	t.Helper()
	parameterID := testName(t)

	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/global", projectID)
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

	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()
	parent := fmt.Sprintf("projects/%s/locations/global/parameters/%s", projectID, parameterID)

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

	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
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
	client, err := parametermanager.NewClient(ctx)
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

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	secret, err := client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", projectID),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
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
	client, err := secretmanager.NewClient(ctx)
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
	client, err := secretmanager.NewClient(ctx)
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
	client, err := secretmanager.NewClient(ctx)
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

// TestCreateStructuredParamVersion tests the createStructuredParamVersion function by creating a structured parameter version,
// then verifies if the parameter version was successfully created by checking the output.
func TestCreateStructuredParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameterVersionID := testName(t)
	payload := `{"username": "test-user", "host": "localhost"}`
	var buf bytes.Buffer
	if err := createStructuredParamVersion(&buf, tc.ProjectID, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := buf.String(), "Created parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestCreateParamVersion tests the createParamVersion function by creating a parameter version,
// then verifies if the parameter version was successfully created by checking the output.
func TestCreateParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	payload := "test123"
	var buf bytes.Buffer
	if err := createParamVersion(&buf, tc.ProjectID, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := buf.String(), "Created parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestCreateParamVersionWithSecret tests the createParamVersionWithSecret function by creating a parameter version with a secret reference,
// then verifies if the parameter version was successfully created by checking the output.
func TestCreateParamVersionWithSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	secretID := testName(t)
	payload := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", tc.ProjectID, secretID)
	var buf bytes.Buffer
	if err := createParamVersionWithSecret(&buf, tc.ProjectID, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := buf.String(), "Created parameter version with secret reference:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestGetParam tests the getParam function by creating a parameter,
// then attempts to retrieve the created parameter. It verifies if the parameter
// was successfully retrieved by checking the output.
func TestGetParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	defer testCleanupParameter(t, parameter.Name)

	var buf bytes.Buffer
	if err := getParam(&buf, tc.ProjectID, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found parameter %s with format JSON", parameter.Name); !strings.Contains(got, want) {
		t.Errorf("GetParameter: expected %q to contain %q", got, want)
	}
}

// TestListParam tests the listParam function by creating multiple parameters,
// then attempts to list the created parameters. It verifies if the parameters
// were successfully listed by checking the output.
func TestListParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter1, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameter2, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)

	defer testCleanupParameter(t, parameter1.Name)
	defer testCleanupParameter(t, parameter2.Name)

	var buf bytes.Buffer
	if err := listParams(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found parameter %s with format %s \n", parameter1.Name, parameter1.Format); !strings.Contains(got, want) {
		t.Errorf("ListParameter: expected %q to contain %q", got, want)
	}

	if got, want := buf.String(), fmt.Sprintf("Found parameter %s with format %s \n", parameter2.Name, parameter2.Format); !strings.Contains(got, want) {
		t.Errorf("ListParameter: expected %q to contain %q", got, want)
	}
}

// TestListParamVersions tests the listParamVersion function by creating a parameter and its versions,
// then attempts to list the created parameter versions. It verifies if the parameter versions
// were successfully listed by checking the output.
func TestListParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion1, _ := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	parameterVersion2, _ := testParameterVersion(t, tc.ProjectID, parameterID, payload)

	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion1.Name)
	defer testCleanupParameterVersion(t, parameterVersion2.Name)

	var buf bytes.Buffer
	if err := listParamVersions(&buf, tc.ProjectID, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found parameter version %s with disabled state in %v\n", parameterVersion1.Name, parameterVersion1.Disabled); !strings.Contains(got, want) {
		t.Errorf("ListParameterVersion: expected %q to contain %q", got, want)
	}

	if got, want := buf.String(), fmt.Sprintf("Found parameter version %s with disabled state in %v\n", parameterVersion2.Name, parameterVersion2.Disabled); !strings.Contains(got, want) {
		t.Errorf("ListParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestGetParamVersion tests the getParamVersion function by creating a parameter and its version,
// then attempts to retrieve the created parameter version. It verifies if the parameter version
// was successfully retrieved by checking the output.
func TestGetParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var buf bytes.Buffer
	if err := getParamVersion(&buf, tc.ProjectID, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found parameter version %s with disabled state in %v", parameterVersion.Name, parameterVersion.Disabled); !strings.Contains(got, want) {
		t.Errorf("GetParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestRenderParamVersion tests the renderParamVersion function by creating a parameter,
// its version, and a secret. It then attempts to render the created parameter version
// and verifies if the parameter version was successfully rendered by checking the output.
func TestRenderParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	secret := testSecret(t, tc.ProjectID)
	testSecretVersion(t, secret.Name, []byte("very secret data"))
	payload := fmt.Sprintf(`{"username": "test-user","password": "__REF__(//secretmanager.googleapis.com/%s/versions/latest)"}`, secret.Name)
	if err := testIamGrantAccess(t, secret.Name, parameter.PolicyMember.IamPolicyUidPrincipal); err != nil {
		t.Fatal(err)
	}
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)

	defer testCleanupSecret(t, secret.Name)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var buf bytes.Buffer
	if err := renderParamVersion(&buf, tc.ProjectID, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Rendered parameter version:"; !strings.Contains(got, want) {
		t.Errorf("RenderParameterVersion: expected %q to contain %q", got, want)
	}

	if got, want := buf.String(), `Rendered payload: {"username": "test-user","password": "very secret data"}`; !strings.Contains(got, want) {
		t.Errorf("RenderParameterVersion: expected %q to contain %q", got, want)
	}
}
