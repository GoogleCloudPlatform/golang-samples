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

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func testName(t *testing.T) string {
	u, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
}

func testLocation(t *testing.T) string {
	v := os.Getenv("GOLANG_REGIONAL_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testLocation: missing GOLANG_REGIONAL_SAMPLES_LOCATION")
	}

	return v
}

func testParameterWithKmsKey(t *testing.T, projectID, kms_key string) (*parametermanagerpb.Parameter, string) {
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

func testCleanupParameter(t *testing.T, name string) {
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

func testCleanupKeyVersions(t *testing.T, name string) {
	ctx := context.Background()

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	if _, err := client.DestroyCryptoKeyVersion(ctx, &kmspb.DestroyCryptoKeyVersionRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupKeyVersion: failed to delete key version: %v", err)
		}
	}
}

func testCreateKeyRing(t *testing.T, projectID, keyRingId string) {
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

func testCreateKeyHSM(t *testing.T, projectID, keyRing, id string) {
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

func TestCreateRegionalParamWithKmsKey(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)
	locationID := testLocation(t)

	keyId := testName(t)
	testCreateKeyRing(t, tc.ProjectID, "go-test-key-ring")
	testCreateKeyHSM(t, tc.ProjectID, "go-test-key-ring", keyId)
	kms_key := fmt.Sprintf("projects/%s/locations/%s/keyRings/go-test-key-ring/cryptoKeys/%s", tc.ProjectID, locationID, keyId)

	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationID, parameterID))
	defer testCleanupKeyVersions(t, fmt.Sprintf("%s/cryptoKeyVersions/1", kms_key))

	var b bytes.Buffer
	if err := createRegionalParamWithKmsKey(&b, tc.ProjectID, locationID, parameterID, kms_key); err != nil {
		t.Fatalf("Failed to create regional parameter: %v", err)
	}
	if got, want := b.String(), fmt.Sprintf("Created regional parameter %s with kms_key %s\n", fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationID, parameterID), kms_key); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

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

	var b bytes.Buffer
	if err := updateRegionalParamKmsKey(&b, tc.ProjectID, locationID, parameterID, kms_key); err != nil {
		t.Fatalf("Failed to update regional parameter: %v", err)
	}
	if got, want := b.String(), fmt.Sprintf("Updated regional parameter %s with kms_key %s\n", fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationID, parameterID), kms_key); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

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

	var b bytes.Buffer
	if err := removeRegionalParamKmsKey(&b, tc.ProjectID, locationID, parameterID); err != nil {
		t.Fatalf("Failed to create regional parameter: %v", err)
	}
	if got, want := b.String(), fmt.Sprintf("Removed kms_key for regional parameter %s\n", fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationID, parameterID)); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}
