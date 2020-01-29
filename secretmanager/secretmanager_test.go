// Copyright 2019 Google LLC
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

package secretmanager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func testClient(tb testing.TB) (*secretmanager.Client, context.Context) {
	tb.Helper()

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		tb.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
}

func testSecret(tb testing.TB, projectID string) *secretmanagerpb.Secret {
	tb.Helper()

	secretID := testName(tb)

	client, ctx := testClient(tb)
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

	client, ctx := testClient(tb)

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

func testCleanupSecret(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testClient(tb)

	if err := client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			tb.Fatalf("testCleanupSecret: failed to delete secret: %v", err)
		}
	}
}

func testIamUser(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
	if v == "" {
		tb.Skip("testIamUser: missing GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
	}

	return fmt.Sprintf("serviceAccount:%s", v)
}

func TestAccessSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := accessSecretVersion(&b, version.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), string(payload); !strings.Contains(got, want) {
		t.Errorf("accessSecretVersion: expected %q to contain %q", got, want)
	}
}

func TestAddSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := addSecretVersion(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Added secret version:"; !strings.Contains(got, want) {
		t.Errorf("addSecretVersion: expected %q to contain %q", got, want)
	}
}

func TestCreateSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createSecret"

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	var b bytes.Buffer
	if err := createSecret(&b, parent, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created secret:"; !strings.Contains(got, want) {
		t.Errorf("createSecret: expected %q to contain %q", got, want)
	}
}

func TestDeleteSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	if err := deleteSecret(secret.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	_, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("deleteSecret: expected %v to be not found", err)
	}
}

func TestDestroySecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version := testSecretVersion(t, secret.Name, payload)

	if err := destroySecretVersion(version.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	v, err := client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_DESTROYED; got != want {
		t.Errorf("testSecretVersion: expected %v to be %v", got, want)
	}
}

func TestDisableEnableSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version := testSecretVersion(t, secret.Name, payload)

	if err := disableSecretVersion(version.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	v, err := client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_DISABLED; got != want {
		t.Errorf("testSecretVersion: expected %v to be %v", got, want)
	}

	if err := enableSecretVersion(version.Name); err != nil {
		t.Fatal(err)
	}

	v, err = client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_ENABLED; got != want {
		t.Errorf("testSecretVersion: expected %v to be %v", got, want)
	}
}

func TestGetSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := getSecretVersion(&b, version.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret version"; !strings.Contains(got, want) {
		t.Errorf("testSecretVersion: expected %q to contain %q", got, want)
	}
}

func TestGetSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := getSecret(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret"; !strings.Contains(got, want) {
		t.Errorf("getSecret: expected %q to contain %q", got, want)
	}
}

func TestIamGrantAccess(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	iamUser := testIamUser(t)

	var b bytes.Buffer
	if err := iamGrantAccess(&b, secret.Name, iamUser); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated IAM policy"; !strings.Contains(got, want) {
		t.Errorf("getSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testClient(t)
	policy, err := client.IAM(secret.Name).Policy(ctx)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	members := policy.Members("roles/secretmanager.secretAccessor")
	for _, m := range members {
		if m == iamUser {
			found = true
		}
	}

	if !found {
		t.Errorf("expected %q to include %q", members, iamUser)
	}
}

func TestIamRevokeAccess(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	iamUser := testIamUser(t)

	var b bytes.Buffer
	if err := iamRevokeAccess(&b, secret.Name, iamUser); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated IAM policy"; !strings.Contains(got, want) {
		t.Errorf("getSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testClient(t)
	policy, err := client.IAM(secret.Name).Policy(ctx)
	if err != nil {
		t.Fatal(err)
	}

	members := policy.Members("roles/secretmanager.secretAccessor")
	for _, m := range members {
		if m == iamUser {
			t.Errorf("expected %q to not include %q", members, iamUser)
		}
	}
}

func TestListSecretVersions(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version1 := testSecretVersion(t, secret.Name, payload)
	version2 := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := listSecretVersions(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), fmt.Sprintf("%s with state ENABLED", version1.Name); !strings.Contains(got, want) {
		t.Errorf("listSecretVersions: expected %q to contain %q", got, want)
	}

	if got, want := b.String(), fmt.Sprintf("%s with state ENABLED", version2.Name); !strings.Contains(got, want) {
		t.Errorf("listSecretVersions: expected %q to contain %q", got, want)
	}
}

func TestListSecrets(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret1 := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret1.Name)

	secret2 := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret2.Name)

	var b bytes.Buffer
	if err := listSecrets(&b, fmt.Sprintf("projects/%s", tc.ProjectID)); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), secret1.Name; !strings.Contains(got, want) {
		t.Errorf("listSecrets: expected %q to contain %q", got, want)
	}

	if got, want := b.String(), secret2.Name; !strings.Contains(got, want) {
		t.Errorf("listSecrets: expected %q to contain %q", got, want)
	}
}

func TestUpdateSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := updateSecret(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{"secretmanager": "rocks"}; !reflect.DeepEqual(got, want) {
		t.Errorf("updateSecret: expected %q to be %q", got, want)
	}
}
