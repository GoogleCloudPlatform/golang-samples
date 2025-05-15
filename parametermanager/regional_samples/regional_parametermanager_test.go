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

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
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

// testParameterWithKmsKey creates a parameter with a KMS key in the specified GCP project.
// It returns the created parameter and its ID or fails the test if parameter creation fails.
func testParameterWithKmsKey(t *testing.T, projectID, kms_key string) (*parametermanagerpb.Parameter, string) {
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
			Format: parametermanagerpb.ParameterFormat_UNFORMATTED,
			KmsKey: &kms_key,
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

	ctx := context.Background()
	locationId := testLocation(t)

	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationId)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	err = client.DeleteParameter(ctx, &parametermanagerpb.DeleteParameterRequest{
		Name: name,
	})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupParameter: failed to delete parameter: %v", err)
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

	err = client.DeleteParameterVersion(ctx, &parametermanagerpb.DeleteParameterVersionRequest{
		Name: name,
	})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupParameterVersion: failed to delete parameter version: %v", err)
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

	err = client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupSecret: failed to delete secret: %v", err)
	}
}

// testCleanupKeyVersions deletes the specified key version in the GCP project.
// It fails the test if the key version deletion fails.
func testCleanupKeyVersions(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	_, err = client.DestroyCryptoKeyVersion(ctx, &kmspb.DestroyCryptoKeyVersionRequest{
		Name: name,
	})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupKeyVersion: failed to delete key version: %v", err)
	}
}

// testCreateKeyRing creates a key ring in the specified GCP project.
// It fails the test if the key ring creation fails.
func testCreateKeyRing(t *testing.T, projectID, keyRingId string) {
	t.Helper()
	ctx := context.Background()
	locationID := testLocation(t)

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Check if key ring already exists
	req := &kmspb.GetKeyRingRequest{
		Name: parent + "/keyRings/" + keyRingId,
	}
	_, err = client.GetKeyRing(ctx, req)
	if err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("failed to get key ring: %v", err)
		}
		// Key ring not found, create it
		req := &kmspb.CreateKeyRingRequest{
			Parent:    parent,
			KeyRingId: keyRingId,
		}
		_, err = client.CreateKeyRing(ctx, req)
		if err != nil {
			t.Fatalf("failed to create key ring: %v", err)
		}
	}
}

// testCreateKeyHSM creates a HSM key in the specified key ring in the GCP project.
// It fails the test if the key creation fails.
func testCreateKeyHSM(t *testing.T, projectID, keyRing, id string) {
	t.Helper()
	ctx := context.Background()
	locationID := testLocation(t)
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", projectID, locationID, keyRing)

	// Check if key already exists
	req := &kmspb.GetCryptoKeyRequest{
		Name: parent + "/cryptoKeys/" + id,
	}
	_, err = client.GetCryptoKey(ctx, req)
	if err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("failed to get crypto key: %v", err)
		}
		// Key not found, create it
		req := &kmspb.CreateCryptoKeyRequest{
			Parent:      parent,
			CryptoKeyId: id,
			CryptoKey: &kmspb.CryptoKey{
				Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
				VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
					ProtectionLevel: kmspb.ProtectionLevel_HSM,
					Algorithm:       kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
				},
			},
		}
		_, err = client.CreateCryptoKey(ctx, req)
		if err != nil {
			t.Fatalf("failed to create crypto key: %v", err)
		}
	}
}

// TestCreateRegionalParam tests the createRegionalParam function by creating a regional parameter,
// then verifies if the parameter was successfully created by checking the output.
func TestCreateRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)
	locationId := testLocation(t)

	var buf bytes.Buffer
	if err := createRegionalParam(&buf, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationId, parameterID))

	if got, want := buf.String(), "Created regional parameter:"; !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

// TestCreateStructuredRegionalParam tests the createStructuredRegionalParam function by creating a structured regional parameter,
// then verifies if the parameter was successfully created by checking the output.
func TestCreateStructuredRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)
	locationId := testLocation(t)

	var buf bytes.Buffer
	if err := createStructuredRegionalParam(&buf, tc.ProjectID, locationId, parameterID, parametermanagerpb.ParameterFormat_JSON); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationId, parameterID))

	if got, want := buf.String(), fmt.Sprintf("Created regional parameter %s with format JSON", fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationId, parameterID)); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

// TestCreateStructuredRegionalParamVersion tests the createStructuredRegionalParamVersion function by creating a structured regional parameter version,
// then verifies if the parameter version was successfully created by checking the output.
func TestCreateStructuredRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameterVersionID := testName(t)
	locationId := testLocation(t)

	payload := `{"username": "test-user", "host": "localhost"}`
	var buf bytes.Buffer
	if err := createStructuredRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := buf.String(), "Created regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestCreateRegionalParamVersion tests the createRegionalParamVersion function by creating a regional parameter version,
// then verifies if the parameter version was successfully created by checking the output.
func TestCreateRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	locationId := testLocation(t)

	payload := "test123"
	var buf bytes.Buffer
	if err := createRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := buf.String(), "Created regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestCreateRegionalParamVersionWithSecret tests the createRegionalParamVersionWithSecret function by creating a regional parameter version with a secret reference,
// then verifies if the parameter version was successfully created by checking the output.
func TestCreateRegionalParamVersionWithSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	locationId := testLocation(t)
	secretID := testName(t)
	payload := fmt.Sprintf("projects/%s/locations/%s/secrets/%s/versions/latest", tc.ProjectID, locationId, secretID)
	var buf bytes.Buffer
	if err := createRegionalParamVersionWithSecret(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := buf.String(), "Created regional parameter version with secret reference:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestDisableRegionalParamVersion tests the disableRegionalParamVersion function by creating a parameter and its version,
// then attempts to disable the created parameter version. It verifies if the parameter version
// was successfully disabled by checking the output.
func TestDisableRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	locationId := testLocation(t)

	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var buf bytes.Buffer
	if err := disableRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Disabled regional parameter version"; !strings.Contains(got, want) {
		t.Errorf("DisableParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestEnableRegionalParamVersion tests the enableRegionalParamVersion function by creating a parameter and its version,
// then attempts to enable the created parameter version. It verifies if the parameter version
// was successfully enabled by checking the output.
func TestEnableRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	locationId := testLocation(t)

	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var buf bytes.Buffer
	if err := enableRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Enabled regional parameter version"; !strings.Contains(got, want) {
		t.Errorf("EnableParameterVersion: expected %q to contain %q", got, want)
	}
}

// TestDeleteRegionalParam tests the deleteRegionalParam function by creating a parameter,
// then attempts to delete the created parameter. It verifies if the parameter
// was successfully deleted by checking the output.
func TestDeleteRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	locationId := testLocation(t)
	defer testCleanupParameter(t, parameter.Name)

	var buf bytes.Buffer
	if err := deleteRegionalParam(&buf, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Deleted regional parameter"; !strings.Contains(got, want) {
		t.Errorf("DeleteParameter: expected %q to contain %q", got, want)
	}
}

// TestDeleteRegionalParamVersion tests the deleteRegionalParamVersion function by creating a parameter and its version,
// then attempts to delete the created parameter version. It verifies if the parameter version
// was successfully deleted by checking the output.
func TestDeleteRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)
	locationId := testLocation(t)

	var buf bytes.Buffer
	if err := deleteRegionalParamVersion(&buf, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Deleted regional parameter version"; !strings.Contains(got, want) {
		t.Errorf("DeleteParameterVersion: expected %q to contain %q", got, want)
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

// TestCreateRegionalParamWithKmsKey tests the createRegionalParamWithKmsKey function by creating a regional parameter with a KMS key,
// and verifies if the parameter was successfully created by checking the output.
func TestCreateRegionalParamWithKmsKey(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)
	locationID := testLocation(t)
	parameterName := fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationID, parameterID)

	keyId := testName(t)
	testCreateKeyRing(t, tc.ProjectID, "go-test-key-ring")
	testCreateKeyHSM(t, tc.ProjectID, "go-test-key-ring", keyId)
	kms_key := fmt.Sprintf("projects/%s/locations/%s/keyRings/go-test-key-ring/cryptoKeys/%s", tc.ProjectID, locationID, keyId)

	defer testCleanupParameter(t, parameterName)
	defer testCleanupKeyVersions(t, fmt.Sprintf("%s/cryptoKeyVersions/1", kms_key))

	var buf bytes.Buffer
	if err := createRegionalParamWithKmsKey(&buf, tc.ProjectID, locationID, parameterID, kms_key); err != nil {
		t.Fatalf("Failed to create regional parameter: %v", err)
	}
	if got, want := buf.String(), fmt.Sprintf("Created regional parameter %s with kms_key %s", parameterName, kms_key); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

// TestUpdateRegionalParamKmsKey tests the updateRegionalParamKmsKey function by creating a regional parameter with a KMS key,
// updating the KMS key, and verifying if the parameter was successfully updated by checking the output.
func TestUpdateRegionalParamKmsKey(t *testing.T) {
	tc := testutil.SystemTest(t)

	locationID := testLocation(t)

	keyId := testName(t)
	testCreateKeyRing(t, tc.ProjectID, "go-test-key-ring")
	testCreateKeyHSM(t, tc.ProjectID, "go-test-key-ring", keyId)
	kms_key := fmt.Sprintf("projects/%s/locations/%s/keyRings/go-test-key-ring/cryptoKeys/%s", tc.ProjectID, locationID, keyId)

	parameter, parameterID := testParameterWithKmsKey(t, tc.ProjectID, kms_key)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupKeyVersions(t, fmt.Sprintf("%s/cryptoKeyVersions/1", kms_key))

	var buf bytes.Buffer
	if err := updateRegionalParamKmsKey(&buf, tc.ProjectID, locationID, parameterID, kms_key); err != nil {
		t.Fatalf("Failed to update regional parameter: %v", err)
	}
	if got, want := buf.String(), fmt.Sprintf("Updated regional parameter %s with kms_key %s", parameter.Name, kms_key); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

// TestRemoveRegionalParamKmsKey tests the removeRegionalParamKmsKey function by creating a regional parameter with a KMS key,
// removing the KMS key, and verifying if the KMS key was successfully removed by checking the output.
func TestRemoveRegionalParamKmsKey(t *testing.T) {
	tc := testutil.SystemTest(t)

	locationID := testLocation(t)

	keyId := testName(t)
	testCreateKeyRing(t, tc.ProjectID, "go-test-key-ring")
	testCreateKeyHSM(t, tc.ProjectID, "go-test-key-ring", keyId)
	kms_key := fmt.Sprintf("projects/%s/locations/%s/keyRings/go-test-key-ring/cryptoKeys/%s", tc.ProjectID, locationID, keyId)

	parameter, parameterID := testParameterWithKmsKey(t, tc.ProjectID, kms_key)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupKeyVersions(t, fmt.Sprintf("%s/cryptoKeyVersions/1", kms_key))

	var buf bytes.Buffer
	if err := removeRegionalParamKmsKey(&buf, tc.ProjectID, locationID, parameterID); err != nil {
		t.Fatalf("Failed to create regional parameter: %v", err)
	}
	if got, want := buf.String(), fmt.Sprintf("Removed kms_key for regional parameter %s", parameter.Name); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}
