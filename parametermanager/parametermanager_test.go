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

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
}

func testClient(tb testing.TB) (*parametermanager.Client, context.Context) {
	tb.Helper()

	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		tb.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testParameter(tb testing.TB, projectID string, format parametermanagerpb.ParameterFormat) (*parametermanagerpb.Parameter, string) {
	tb.Helper()

	parameterID := testName(tb)

	client, ctx := testClient(tb)
	parent := fmt.Sprintf("projects/%s/locations/global", projectID)
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

	client, ctx := testClient(tb)
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

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		tb.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testSecret(tb testing.TB, projectID string) *secretmanagerpb.Secret {
	tb.Helper()

	secretID := testName(tb)

	client, ctx := testClientForSecret(tb)
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

func TestCreateParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)

	var b bytes.Buffer
	if err := createParam(&b, tc.ProjectID, parameterID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/global/parameters/%s", tc.ProjectID, parameterID))

	if got, want := b.String(), "Created parameter:"; !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

func TestCreateStructuredParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameterID := testName(t)

	var b bytes.Buffer
	if err := createStructuredParam(&b, tc.ProjectID, parameterID, parametermanagerpb.ParameterFormat_JSON); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, fmt.Sprintf("projects/%s/locations/global/parameters/%s", tc.ProjectID, parameterID))

	if got, want := b.String(), fmt.Sprintf("Created parameter %s with format JSON", fmt.Sprintf("projects/%s/locations/global/parameters/%s", tc.ProjectID, parameterID)); !strings.Contains(got, want) {
		t.Errorf("createParameter: expected %q to contain %q", got, want)
	}
}

func TestCreateStructuredParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	parameterVersionID := testName(t)
	payload := `{"username": "test-user", "host": "localhost"}`
	var b bytes.Buffer
	if err := createStructuredParamVersion(&b, tc.ProjectID, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := b.String(), "Created parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestCreateParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	payload := "test123"
	var b bytes.Buffer
	if err := createParamVersion(&b, tc.ProjectID, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := b.String(), "Created parameter version:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestCreateParamVersionWithSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_UNFORMATTED)
	parameterVersionID := testName(t)
	secretID := testName(t)
	payload := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", tc.ProjectID, secretID)
	var b bytes.Buffer
	if err := createParamVersionWithSecret(&b, tc.ProjectID, parameterID, parameterVersionID, payload); err != nil {
		t.Fatal(err)
	}
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, fmt.Sprintf("%s/versions/%s", parameter.Name, parameterVersionID))

	if got, want := b.String(), "Created parameter version with secret reference:"; !strings.Contains(got, want) {
		t.Errorf("createParameterVersion: expected %q to contain %q", got, want)
	}
}
