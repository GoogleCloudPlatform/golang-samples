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
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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

// testNewClient creates a Parameter Manager client for the test.
func testNewClient(t *testing.T) *parametermanager.Client {
	t.Helper()

	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", testLocation(t))
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("testNewClient: failed to create client: %v", err)
	}
	return client
}

// testLocationPath returns the parent resource path of the test location.
func testLocationPath(t *testing.T, projectID string) string {
	t.Helper()

	return fmt.Sprintf("projects/%s/locations/%s", projectID, testLocation(t))
}

// testTemplate creates a template with the given format in the specified GCP project.
// It returns the created template and its ID or fails the test if template creation fails.
func testTemplate(t *testing.T, projectID string, format parametermanagerpb.TemplateFormat) (*parametermanagerpb.Template, string) {
	t.Helper()

	templateID := testName(t)
	client := testNewClient(t)
	defer client.Close()

	template, err := client.CreateTemplate(context.Background(), &parametermanagerpb.CreateTemplateRequest{
		Parent:     testLocationPath(t, projectID),
		TemplateId: templateID,
		Template: &parametermanagerpb.Template{
			Format: format,
		},
	})
	if err != nil {
		t.Fatalf("testTemplate: failed to create template: %v", err)
	}

	return template, templateID
}

// testTemplateVersion creates a version of a template with the given payload.
// It returns the created template version and its ID or fails the test if creation fails.
func testTemplateVersion(t *testing.T, templateName, payload string) (*parametermanagerpb.TemplateVersion, string) {
	t.Helper()

	versionID := testName(t)
	client := testNewClient(t)
	defer client.Close()

	version, err := client.CreateTemplateVersion(context.Background(), &parametermanagerpb.CreateTemplateVersionRequest{
		Parent:            templateName,
		TemplateVersionId: versionID,
		TemplateVersion: &parametermanagerpb.TemplateVersion{
			Payload: &parametermanagerpb.TemplateVersionPayload{
				Data: []byte(payload),
			},
		},
	})
	if err != nil {
		t.Fatalf("testTemplateVersion: failed to create template version: %v", err)
	}

	return version, versionID
}

// testCleanupTemplate deletes the specified template in the GCP project.
// It fails the test if the template deletion fails.
func testCleanupTemplate(t *testing.T, name string) {
	t.Helper()

	client := testNewClient(t)
	defer client.Close()

	err := client.DeleteTemplate(context.Background(), &parametermanagerpb.DeleteTemplateRequest{
		Name: name,
	})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
	}
}

// testCleanupTemplateVersion deletes the specified template version in the GCP project.
// It fails the test if the template version deletion fails.
func testCleanupTemplateVersion(t *testing.T, name string) {
	t.Helper()

	client := testNewClient(t)
	defer client.Close()

	err := client.DeleteTemplateVersion(context.Background(), &parametermanagerpb.DeleteTemplateVersionRequest{
		Name: name,
	})
	if err == nil {
		return
	}
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Fatalf("testCleanupTemplateVersion: failed to delete template version: %v", err)
	}
}

// testTag returns the pre-provisioned tag key and tag value used by the tag tests.
// The tests are skipped if GOLANG_SAMPLES_TAG_KEY or GOLANG_SAMPLES_TAG_VALUE is not set.
func testTag(t *testing.T) (string, string) {
	t.Helper()

	key := os.Getenv("GOLANG_SAMPLES_TAG_KEY")
	value := os.Getenv("GOLANG_SAMPLES_TAG_VALUE")
	if key == "" || value == "" {
		t.Skip("testTag: missing GOLANG_SAMPLES_TAG_KEY or GOLANG_SAMPLES_TAG_VALUE")
	}

	return key, value
}

// testTagBindings returns the tag values bound to the named parameter, read through Resource Manager.
func testTagBindings(t *testing.T, parameterName string) []string {
	t.Helper()

	ctx := context.Background()
	endpoint := fmt.Sprintf("%s-cloudresourcemanager.googleapis.com:443", testLocation(t))
	client, err := resourcemanager.NewTagBindingsClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("testTagBindings: failed to create client: %v", err)
	}
	defer client.Close()

	var values []string
	it := client.ListTagBindings(ctx, &resourcemanagerpb.ListTagBindingsRequest{
		Parent: fmt.Sprintf("//parametermanager.googleapis.com/%s", parameterName),
	})
	for {
		binding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("testTagBindings: failed to list tag bindings: %v", err)
		}
		values = append(values, binding.TagValue)
	}
	return values
}

// TestCreateRegionalParamTemplate tests the createRegionalParamTemplate function by creating a template,
// then verifies if the template was successfully created by checking the output.
func TestCreateRegionalParamTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := testName(t)
	locationId := testLocation(t)
	var buf bytes.Buffer
	if err := createRegionalParamTemplate(&buf, tc.ProjectID, locationId, templateID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("%s/templates/%s", testLocationPath(t, tc.ProjectID), templateID))

	if got, want := buf.String(), "Created regional parameter template:"; !strings.Contains(got, want) {
		t.Errorf("createRegionalParamTemplate: expected %q to contain %q", got, want)
	}
	if got, want := buf.String(), "TEMPLATE_FORMAT_JSON"; !strings.Contains(got, want) {
		t.Errorf("createRegionalParamTemplate: expected %q to contain %q", got, want)
	}
}

// TestCreateRegionalParamTemplateVersion tests the createRegionalParamTemplateVersion function by creating a
// template version with placeholders, then verifies the output and the stored payload.
func TestCreateRegionalParamTemplateVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	versionID := testName(t)
	payload := `{"username": "{{.username}}"}`
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, fmt.Sprintf("%s/versions/%s", template.Name, versionID))

	var buf bytes.Buffer
	if err := createRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, versionID, payload); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Created regional parameter template version:"; !strings.Contains(got, want) {
		t.Errorf("createRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	version, err := client.GetTemplateVersion(context.Background(), &parametermanagerpb.GetTemplateVersionRequest{
		Name: fmt.Sprintf("%s/versions/%s", template.Name, versionID),
	})
	if err != nil {
		t.Fatalf("failed to get template version: %v", err)
	}
	if got := string(version.Payload.Data); got != payload {
		t.Errorf("createRegionalParamTemplateVersion: got payload %q, want %q", got, payload)
	}
}

// TestListRegionalParamTemplates tests the listRegionalParamTemplates function by creating templates,
// then verifies that they are listed.
func TestListRegionalParamTemplates(t *testing.T) {
	tc := testutil.SystemTest(t)

	template1, templateID1 := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	template2, templateID2 := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_YAML)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template1.Name)
	defer testCleanupTemplate(t, template2.Name)

	var buf bytes.Buffer
	if err := listRegionalParamTemplates(&buf, tc.ProjectID, locationId); err != nil {
		t.Fatal(err)
	}

	for _, id := range []string{templateID1, templateID2} {
		if got, want := buf.String(), fmt.Sprintf("Found regional parameter template: %s/templates/%s", testLocationPath(t, tc.ProjectID), id); !strings.Contains(got, want) {
			t.Errorf("listRegionalParamTemplates: expected %q to contain %q", got, want)
		}
	}
}

// TestGetRegionalParamTemplate tests the getRegionalParamTemplate function by creating a template,
// then verifies that it is retrieved.
func TestGetRegionalParamTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_YAML)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)

	var buf bytes.Buffer
	if err := getRegionalParamTemplate(&buf, tc.ProjectID, locationId, templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter template %s with format TEMPLATE_FORMAT_YAML", template.Name); !strings.Contains(got, want) {
		t.Errorf("getRegionalParamTemplate: expected %q to contain %q", got, want)
	}
}

// TestListRegionalParamTemplateVersions tests the listRegionalParamTemplateVersions function by creating
// template versions, then verifies that they are listed.
func TestListRegionalParamTemplateVersions(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	version1, _ := testTemplateVersion(t, template.Name, `{"a": "{{.a}}"}`)
	version2, _ := testTemplateVersion(t, template.Name, `{"b": "{{.b}}"}`)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, version1.Name)
	defer testCleanupTemplateVersion(t, version2.Name)

	var buf bytes.Buffer
	if err := listRegionalParamTemplateVersions(&buf, tc.ProjectID, locationId, templateID); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{version1.Name, version2.Name} {
		if got, want := buf.String(), fmt.Sprintf("Found regional parameter template version: %s", name); !strings.Contains(got, want) {
			t.Errorf("listRegionalParamTemplateVersions: expected %q to contain %q", got, want)
		}
	}
}

// TestGetRegionalParamTemplateVersion tests the getRegionalParamTemplateVersion function by creating a
// template version, then verifies that it and its payload are retrieved.
func TestGetRegionalParamTemplateVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	payload := `{"username": "{{.username}}"}`
	version, versionID := testTemplateVersion(t, template.Name, payload)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, version.Name)

	var buf bytes.Buffer
	if err := getRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, versionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Found regional parameter template version %s with disabled state in false", version.Name); !strings.Contains(got, want) {
		t.Errorf("getRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}
	if got, want := buf.String(), fmt.Sprintf("Payload: %s", payload); !strings.Contains(got, want) {
		t.Errorf("getRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}
}

// TestUpdateRegionalParamTemplateLabels tests the updateRegionalParamTemplateLabels function by creating a
// template, then verifies that the label was applied.
func TestUpdateRegionalParamTemplateLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)

	var buf bytes.Buffer
	if err := updateRegionalParamTemplateLabels(&buf, tc.ProjectID, locationId, templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "environment:test"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalParamTemplateLabels: expected %q to contain %q", got, want)
	}
}

// TestDeleteRegionalParamTemplate tests the deleteRegionalParamTemplate function by creating a template,
// deleting it, then verifies that it no longer exists.
func TestDeleteRegionalParamTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)

	var buf bytes.Buffer
	if err := deleteRegionalParamTemplate(&buf, tc.ProjectID, locationId, templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Deleted regional parameter template:"; !strings.Contains(got, want) {
		t.Errorf("deleteRegionalParamTemplate: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	_, err := client.GetTemplate(context.Background(), &parametermanagerpb.GetTemplateRequest{Name: template.Name})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("deleteRegionalParamTemplate: expected NotFound after deletion, got %v", err)
	}
}

// TestDisableRegionalParamTemplateVersion tests the disableRegionalParamTemplateVersion function by creating a
// template version, disabling it, then verifies the disabled state.
func TestDisableRegionalParamTemplateVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	version, versionID := testTemplateVersion(t, template.Name, `{"a": "{{.a}}"}`)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, version.Name)

	var buf bytes.Buffer
	if err := disableRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, versionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Disabled regional parameter template version:"; !strings.Contains(got, want) {
		t.Errorf("disableRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	got, err := client.GetTemplateVersion(context.Background(), &parametermanagerpb.GetTemplateVersionRequest{Name: version.Name})
	if err != nil {
		t.Fatalf("failed to get template version: %v", err)
	}
	if !got.Disabled {
		t.Errorf("disableRegionalParamTemplateVersion: expected template version to be disabled")
	}
}

// TestEnableRegionalParamTemplateVersion tests the enableRegionalParamTemplateVersion function by creating a
// template version, disabling it, enabling it, then verifies the enabled state.
func TestEnableRegionalParamTemplateVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	version, versionID := testTemplateVersion(t, template.Name, `{"a": "{{.a}}"}`)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, version.Name)

	client := testNewClient(t)
	defer client.Close()
	if _, err := client.UpdateTemplateVersion(context.Background(), &parametermanagerpb.UpdateTemplateVersionRequest{
		TemplateVersion: &parametermanagerpb.TemplateVersion{Name: version.Name, Disabled: true},
		UpdateMask:      &fieldmaskpb.FieldMask{Paths: []string{"disabled"}},
	}); err != nil {
		t.Fatalf("failed to disable template version: %v", err)
	}

	var buf bytes.Buffer
	if err := enableRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, versionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Enabled regional parameter template version:"; !strings.Contains(got, want) {
		t.Errorf("enableRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}

	got, err := client.GetTemplateVersion(context.Background(), &parametermanagerpb.GetTemplateVersionRequest{Name: version.Name})
	if err != nil {
		t.Fatalf("failed to get template version: %v", err)
	}
	if got.Disabled {
		t.Errorf("enableRegionalParamTemplateVersion: expected template version to be enabled")
	}
}

// TestDeleteRegionalParamTemplateVersion tests the deleteRegionalParamTemplateVersion function by creating a
// template version, deleting it, then verifies that it no longer exists.
func TestDeleteRegionalParamTemplateVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	version, versionID := testTemplateVersion(t, template.Name, `{"a": "{{.a}}"}`)
	locationId := testLocation(t)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, version.Name)

	var buf bytes.Buffer
	if err := deleteRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, versionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Deleted regional parameter template version:"; !strings.Contains(got, want) {
		t.Errorf("deleteRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	_, err := client.GetTemplateVersion(context.Background(), &parametermanagerpb.GetTemplateVersionRequest{Name: version.Name})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("deleteRegionalParamTemplateVersion: expected NotFound after deletion, got %v", err)
	}
}

// TestRenderRegionalParamTemplateVersion tests the renderRegionalParamTemplateVersion function. The template
// references a Secret Manager secret, which the parameter's identity is granted access to, and the
// test verifies that the secret value appears in the rendered payload.
func TestRenderRegionalParamTemplateVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	secret := testSecret(t, tc.ProjectID)
	testSecretVersion(t, secret.Name, []byte("very secret data"))
	if err := testIamGrantAccess(t, secret.Name, parameter.PolicyMember.IamPolicyUidPrincipal); err != nil {
		t.Fatal(err)
	}
	templatePayload := `{"username": "{{.username}}", "password": "{{.password}}"}`
	templateVersion, templateVersionID := testTemplateVersion(t, template.Name, templatePayload)
	parameterPayload := fmt.Sprintf(`{"username": "test-user", "password": "__REF__(//secretmanager.googleapis.com/%s/versions/latest)"}`, secret.Name)
	parameterVersion, _ := testParameterVersion(t, tc.ProjectID, parameterID, parameterPayload)
	locationId := testLocation(t)
	defer testCleanupSecret(t, secret.Name)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, templateVersion.Name)

	var buf bytes.Buffer
	time.Sleep(2 * time.Minute)
	if err := renderRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, templateVersionID, parameterVersion.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Rendered regional parameter template version:"; !strings.Contains(got, want) {
		t.Errorf("renderRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}
	if got, want := buf.String(), fmt.Sprintf("Template payload: %s", templatePayload); !strings.Contains(got, want) {
		t.Errorf("renderRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}
	expected := `{"username": "test-user", "password": "very secret data"}`
	if got, want := buf.String(), fmt.Sprintf("Rendered payload: %s", expected); !strings.Contains(got, want) {
		t.Errorf("renderRegionalParamTemplateVersion: expected %q to contain %q", got, want)
	}
}

// TestRenderRegionalParamTemplateVersionMissingSecret verifies that rendering fails when the parameter
// version references a secret that does not exist.
func TestRenderRegionalParamTemplateVersionMissingSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	template, templateID := testTemplate(t, tc.ProjectID, parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON)
	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	templateVersion, templateVersionID := testTemplateVersion(t, template.Name, `{"password": "{{.password}}"}`)
	parameterPayload := fmt.Sprintf(`{"password": "__REF__(//secretmanager.googleapis.com/projects/%s/locations/%s/secrets/%s/versions/latest)"}`, tc.ProjectID, testLocation(t), testName(t))
	parameterVersion, _ := testParameterVersion(t, tc.ProjectID, parameterID, parameterPayload)
	locationId := testLocation(t)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)
	defer testCleanupTemplate(t, template.Name)
	defer testCleanupTemplateVersion(t, templateVersion.Name)

	var buf bytes.Buffer
	err := renderRegionalParamTemplateVersion(&buf, tc.ProjectID, locationId, templateID, templateVersionID, parameterVersion.Name)
	if err == nil {
		t.Fatalf("renderRegionalParamTemplateVersion: expected an error for a missing secret, got output %q", buf.String())
	}
	if terr, ok := grpcstatus.FromError(errors.Unwrap(err)); !ok || terr.Code() != grpccodes.FailedPrecondition {
		t.Errorf("renderRegionalParamTemplateVersion: expected FailedPrecondition for a missing secret, got %v", err)
	}
}

// TestCreateRegionalParamWithTags tests the createRegionalParamWithTags function by creating a parameter with
// a tag, then verifies the tag binding through Resource Manager (tags are never returned by
// Parameter Manager).
func TestCreateRegionalParamWithTags(t *testing.T) {
	tc := testutil.SystemTest(t)
	tagKey, tagValue := testTag(t)

	parameterID := testName(t)
	locationId := testLocation(t)
	var buf bytes.Buffer
	if err := createRegionalParamWithTags(&buf, tc.ProjectID, locationId, parameterID, tagKey, tagValue); err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("%s/parameters/%s", testLocationPath(t, tc.ProjectID), parameterID)
	defer testCleanupParameter(t, name)

	if got, want := buf.String(), "Created regional parameter"; !strings.Contains(got, want) {
		t.Errorf("createRegionalParamWithTags: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	parameter, err := client.GetParameter(context.Background(), &parametermanagerpb.GetParameterRequest{Name: name})
	if err != nil {
		t.Fatalf("failed to get parameter: %v", err)
	}
	if got := testTagBindings(t, parameter.Name); len(got) != 1 || got[0] != tagValue {
		t.Errorf("createRegionalParamWithTags: got tag bindings %v, want [%s]", got, tagValue)
	}
}

// TestBindRegionalTagsToParam tests the bindRegionalTagsToParam function by creating a parameter, binding an
// existing tag value, then verifies the tag binding through Resource Manager.
func TestBindRegionalTagsToParam(t *testing.T) {
	tc := testutil.SystemTest(t)
	_, tagValue := testTag(t)

	parameterID := testName(t)
	locationId := testLocation(t)
	var buf bytes.Buffer
	if err := bindRegionalTagsToParam(&buf, tc.ProjectID, locationId, parameterID, tagValue); err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("%s/parameters/%s", testLocationPath(t, tc.ProjectID), parameterID)
	defer testCleanupParameter(t, name)

	if got, want := buf.String(), fmt.Sprintf("Bound tag value %s to regional parameter", tagValue); !strings.Contains(got, want) {
		t.Errorf("bindRegionalTagsToParam: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	parameter, err := client.GetParameter(context.Background(), &parametermanagerpb.GetParameterRequest{Name: name})
	if err != nil {
		t.Fatalf("failed to get parameter: %v", err)
	}
	if got := testTagBindings(t, parameter.Name); len(got) != 1 || got[0] != tagValue {
		t.Errorf("bindRegionalTagsToParam: got tag bindings %v, want [%s]", got, tagValue)
	}
}

// TestGetRegionalParamTags tests the getRegionalParamTags function by creating a parameter with a tag,
// then verifies the tag binding is listed.
func TestGetRegionalParamTags(t *testing.T) {
	tc := testutil.SystemTest(t)
	tagKey, tagValue := testTag(t)

	parameterID := testName(t)
	locationId := testLocation(t)
	if err := createRegionalParamWithTags(io.Discard, tc.ProjectID, locationId, parameterID, tagKey, tagValue); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("%s/parameters/%s", testLocationPath(t, tc.ProjectID), parameterID))

	var buf bytes.Buffer
	if err := getRegionalParamTags(&buf, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Found tag binding on regional parameter"; !strings.Contains(got, want) {
		t.Errorf("getRegionalParamTags: expected %q to contain %q", got, want)
	}
	if got, want := buf.String(), tagValue; !strings.Contains(got, want) {
		t.Errorf("getRegionalParamTags: expected %q to contain %q", got, want)
	}
}

// TestCreateRegionalParamVersionWithChecksum tests the createRegionalParamVersionWithChecksum function by
// creating a version with a client-computed CRC32C, then verifies the checksum source and value.
func TestCreateRegionalParamVersionWithChecksum(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	versionID := testName(t)
	payload := "checksum payload"
	locationId := testLocation(t)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, versionID))

	var buf bytes.Buffer
	if err := createRegionalParamVersionWithChecksum(&buf, tc.ProjectID, locationId, parameterID, versionID, payload); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Created regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createRegionalParamVersionWithChecksum: expected %q to contain %q", got, want)
	}
	if got, want := buf.String(), "checksum source USER_SPECIFIED"; !strings.Contains(got, want) {
		t.Errorf("createRegionalParamVersionWithChecksum: expected %q to contain %q", got, want)
	}

	client := testNewClient(t)
	defer client.Close()
	version, err := client.GetParameterVersion(context.Background(), &parametermanagerpb.GetParameterVersionRequest{
		Name: fmt.Sprintf("%s/versions/%s", parameter.Name, versionID),
		View: parametermanagerpb.View_FULL,
	})
	if err != nil {
		t.Fatalf("failed to get parameter version: %v", err)
	}
	want := int64(crc32.Checksum([]byte(payload), crc32.MakeTable(crc32.Castagnoli)))
	if version.Payload.DataCrc32C == nil || *version.Payload.DataCrc32C != want {
		t.Errorf("createRegionalParamVersionWithChecksum: got data_crc32c %v, want %d", version.Payload.DataCrc32C, want)
	}
}

// TestGetRegionalParamVersionVerifyChecksum tests the getRegionalParamVersionVerifyChecksum function by
// creating a version without a checksum, then verifies the server-generated checksum.
func TestGetRegionalParamVersionVerifyChecksum(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	version, versionID := testParameterVersion(t, tc.ProjectID, parameterID, "checksum payload")
	locationId := testLocation(t)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, version.Name)

	var buf bytes.Buffer
	if err := getRegionalParamVersionVerifyChecksum(&buf, tc.ProjectID, locationId, parameterID, versionID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Verified checksum of regional parameter version"; !strings.Contains(got, want) {
		t.Errorf("getRegionalParamVersionVerifyChecksum: expected %q to contain %q", got, want)
	}
	if got, want := buf.String(), "SERVER_GENERATED"; !strings.Contains(got, want) {
		t.Errorf("getRegionalParamVersionVerifyChecksum: expected %q to contain %q", got, want)
	}
}

// TestCreateRegionalParamVersionChecksumMismatch verifies that Parameter Manager rejects a version whose
// client-supplied CRC32C does not match the payload.
func TestCreateRegionalParamVersionChecksumMismatch(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	defer testCleanupParameter(t, parameter.Name)

	payload := []byte("checksum payload")
	wrong := int64(crc32.Checksum(payload, crc32.MakeTable(crc32.Castagnoli))) + 1

	client := testNewClient(t)
	defer client.Close()
	_, err := client.CreateParameterVersion(context.Background(), &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parameter.Name,
		ParameterVersionId: testName(t),
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data:       payload,
				DataCrc32C: &wrong,
			},
		},
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.InvalidArgument {
		t.Fatalf("expected InvalidArgument for a mismatched checksum, got %v", err)
	}
	if !strings.Contains(err.Error(), "CHECKSUM_MISMATCH") {
		t.Errorf("expected error to contain CHECKSUM_MISMATCH, got %v", err)
	}
}
