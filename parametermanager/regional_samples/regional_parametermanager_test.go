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

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
}

func testLocation(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_REGIONAL_SAMPLES_LOCATION")
	if v == "" {
		tb.Skip("testIamUser: missing GOLANG_REGIONAL_SAMPLES_LOCATION")
	}

	return v
}

func testClient(tb testing.TB) (*parametermanager.Client, context.Context) {
	tb.Helper()

	locationId := testLocation(tb)
	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationId)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		tb.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testParameter(tb testing.TB, projectID string, format parametermanagerpb.ParameterFormat) (*parametermanagerpb.Parameter, string) {
	tb.Helper()

	parameterID := testName(tb)
	locationId := testLocation(tb)

	client, ctx := testClient(tb)
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationId)
	parameter, err := client.CreateParameter(ctx, &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: format,
		},
	})
	if err != nil {
		tb.Fatalf("testParameter: failed to create parameter: %v", err)
	}

	return parameter, parameterID
}

func testParameterVersion(tb testing.TB, projectID, parameterID, payload string) (*parametermanagerpb.ParameterVersion, string) {
	tb.Helper()
	parameterVersionID := testName(tb)
	locationId := testLocation(tb)

	client, ctx := testClient(tb)
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
		tb.Fatalf("testParameterVersion: failed to create parameter version: %v", err)
	}

	return parameterVersion, parameterVersionID
}

func testCleanupParameter(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testClient(tb)

	if err := client.DeleteParameter(ctx, &parametermanagerpb.DeleteParameterRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			tb.Fatalf("testCleanupParameter: failed to delete parameter: %v", err)
		}
	}
}

func testCleanupParameterVersion(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testClient(tb)

	if err := client.DeleteParameterVersion(ctx, &parametermanagerpb.DeleteParameterVersionRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			tb.Fatalf("testCleanupParameterVersion: failed to delete parameter version: %v", err)
		}
	}
}

func testClientForSecret(tb testing.TB) (*secretmanager.Client, context.Context) {
	tb.Helper()

	locationId := testLocation(tb)

	ctx := context.Background()
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		tb.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testSecret(tb testing.TB, projectID string) *secretmanagerpb.Secret {
	tb.Helper()

	secretID := testName(tb)
	locationId := testLocation(tb)

	client, ctx := testClientForSecret(tb)
	secret, err := client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, locationId),
		SecretId: secretID,
		Secret:   &secretmanagerpb.Secret{},
	})
	if err != nil {
		tb.Fatalf("testSecret: failed to create secret: %v", err)
	}

	return secret
}

func testSecretVersion(tb testing.TB, parent string, payload []byte) *secretmanagerpb.SecretVersion {
	tb.Helper()

	client, ctx := testClientForSecret(tb)

	version, err := client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	})
	if err != nil {
		tb.Fatalf("testSecretVersion: failed to create secret version: %v", err)
	}
	return version
}

func testIamGrantAccess(tb testing.TB, name, member string) error {
	tb.Helper()

	client, ctx := testClientForSecret(tb)

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

func testCleanupSecret(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testClientForSecret(tb)

	if err := client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			tb.Fatalf("testCleanupSecret: failed to delete secret: %v", err)
		}
	}
}

func TestCreateRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)
	locationId := testLocation(t)

	var b bytes.Buffer
	if err := createRegionalParam(&b, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationId, parameterID))

	if got, want := b.String(), "Created regional parameter:"; !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

func TestCreateStructuredRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)
	locationId := testLocation(t)

	var b bytes.Buffer
	if err := createStructuredRegionalParam(&b, tc.ProjectID, locationId, parameterID, parametermanagerpb.ParameterFormat_JSON); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/%s/parameters/%s", tc.ProjectID, locationId, parameterID))

	if got, want := b.String(), "Created regional parameter with format:"; !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

func TestCreateStructuredRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameterVersionID := testName(t)
	locationId := testLocation(t)

	payload := `{"username": "test-user", "host": "localhost"}`
	var b bytes.Buffer
	if err := createStructuredRegionalParamVersion(&b, tc.ProjectID, locationId, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := b.String(), "Created regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestCreateRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	locationId := testLocation(t)

	payload := "test123"
	var b bytes.Buffer
	if err := createRegionalParamVersion(&b, tc.ProjectID, locationId, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := b.String(), "Created regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestCreateRegionalParamVersionWithSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	locationId := testLocation(t)
	secretID := testName(t)
	payload := fmt.Sprintf("projects/%s/locations/%s/secrets/%s/versions/latest", tc.ProjectID, locationId, secretID)
	var b bytes.Buffer
	if err := createRegionalParamVersionWithSecret(&b, tc.ProjectID, locationId, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := b.String(), "Created regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestGetRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	defer testCleanupParameter(t, parameter.Name)

	locationId := testLocation(t)
	var b bytes.Buffer
	if err := getRegionalParam(&b, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found regional parameter:"; !strings.Contains(got, want) {
		t.Errorf("GetParameter: expected %q to contain %q", got, want)
	}
}

func TestGetRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)
	locationId := testLocation(t)

	var b bytes.Buffer
	if err := getRegionalParamVersion(&b, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("GetParameterVersion: expected %q to contain %q", got, want)
	}
}

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

	var b bytes.Buffer
	time.Sleep(2 * time.Minute)
	if err := renderRegionalParamVersion(&b, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Rendered regional Parameter Version:"; !strings.Contains(got, want) {
		t.Errorf("RenderParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestListRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter1, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameter2, _ := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	locationId := testLocation(t)

	defer testCleanupParameter(t, parameter1.Name)
	defer testCleanupParameter(t, parameter2.Name)

	var b bytes.Buffer
	if err := listRegionalParam(&b, tc.ProjectID, locationId); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), fmt.Sprintf("Found regional parameter: %s with format %s \n", parameter1.Name, parameter1.Format); !strings.Contains(got, want) {
		t.Errorf("ListParameter: expected %q to contain %q", got, want)
	}

	if got, want := b.String(), fmt.Sprintf("Found regional parameter: %s with format %s \n", parameter2.Name, parameter2.Format); !strings.Contains(got, want) {
		t.Errorf("ListParameter: expected %q to contain %q", got, want)
	}
}

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

	var b bytes.Buffer
	if err := listRegionalParamVersion(&b, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), fmt.Sprintf("Found regional parameter version: %s with disabled: %v\n", parameterVersion1.Name, parameterVersion1.Disabled); !strings.Contains(got, want) {
		t.Errorf("ListParameterVersion: expected %q to contain %q", got, want)
	}

	if got, want := b.String(), fmt.Sprintf("Found regional parameter version: %s with disabled: %v\n", parameterVersion2.Name, parameterVersion2.Disabled); !strings.Contains(got, want) {
		t.Errorf("ListParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestDisableEnableRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	locationId := testLocation(t)

	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var b1 bytes.Buffer
	if err := disableRegionalParamVersion(&b1, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b1.String(), "Disabled regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("DisableParameterVersion: expected %q to contain %q", got, want)
	}

	var b2 bytes.Buffer
	if err := enableRegionalParamVersion(&b2, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b2.String(), "Enabled regional parameter version:"; !strings.Contains(got, want) {
		t.Errorf("DisableParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestDeleteRegionalParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	locationId := testLocation(t)
	defer testCleanupParameter(t, parameter.Name)

	var b bytes.Buffer
	if err := deleteRegionalParam(&b, tc.ProjectID, locationId, parameterID); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	_, err := client.GetParameter(ctx, &parametermanagerpb.GetParameterRequest{
		Name: parameter.Name,
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("DeleteParameter: expected %v to be not found", err)
	}
}

func TestDeleteRegionalParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)
	locationId := testLocation(t)

	var b bytes.Buffer
	if err := deleteRegionalParamVersion(&b, tc.ProjectID, locationId, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	_, err := client.GetParameterVersion(ctx, &parametermanagerpb.GetParameterVersionRequest{
		Name: parameterVersion.Name,
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("DeleteParameterVersion: expected %v to be not found", err)
	}
}
