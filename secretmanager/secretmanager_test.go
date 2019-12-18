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
	"reflect"
	"strings"
	"testing"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	secretspb "google.golang.org/genproto/googleapis/cloud/secrets/v1beta1"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func testClient(tb testing.TB) (*secretmanager.Client, context.Context) {
	tb.Helper()

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		tb.Fatalf("failed to create client: %v", err)
	}
	return client, ctx
}

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("failed to generate uuid: %v", err)
	}
	return u.String()
}

func testSecret(tb testing.TB, projectID string) (*secretspb.Secret, func()) {
	tb.Helper()

	client, ctx := testClient(tb)
	name := testName(tb)

	secret, err := client.CreateSecret(ctx, &secretspb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", projectID),
		SecretId: name,
		Secret: &secretspb.Secret{
			Replication: &secretspb.Replication{
				Replication: &secretspb.Replication_Automatic_{
					Automatic: &secretspb.Replication_Automatic{},
				},
			},
		},
	})
	if err != nil {
		tb.Fatalf("failed to create secret: %v", err)
	}

	return secret, func() { testCleanupSecret(tb, secret.Name) }
}

func testSecretVersion(tb testing.TB, parent string, payload []byte) *secretspb.SecretVersion {
	tb.Helper()

	client, ctx := testClient(tb)

	version, err := client.AddSecretVersion(ctx, &secretspb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretspb.SecretPayload{
			Data: payload,
		},
	})
	if err != nil {
		tb.Fatalf("failed to create secret version: %v", err)
	}
	return version
}

func testCleanupSecret(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testClient(tb)

	_ = client.DeleteSecret(ctx, &secretspb.DeleteSecretRequest{
		Name: name,
	})
}

func TestAccessSecretVersion(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	payload := []byte("my-secret")
	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()
	version := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := accessSecretVersion(&b, version.Name); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), string(payload); !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestAddSecretVersion(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()

	var b bytes.Buffer
	if err := addSecretVersion(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), "Added secret version:"; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestCreateSecret(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)
	name := testName(t)

	var b bytes.Buffer
	if err := createSecret(&b, parent, name); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("%s/secrets/%s", parent, name))

	if act, exp := b.String(), "Created secret:"; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestDeleteSecret(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()

	if err := deleteSecret(secret.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	_, err := client.GetSecret(ctx, &secretspb.GetSecretRequest{
		Name: secret.Name,
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("expected %v to be not found", err)
	}
}

func TestDestroySecretVersion(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	payload := []byte("my-secret")
	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()
	version := testSecretVersion(t, secret.Name, payload)

	if err := destroySecretVersion(version.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	v, err := client.GetSecretVersion(ctx, &secretspb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if act, exp := v.State, secretspb.SecretVersion_DESTROYED; act != exp {
		t.Errorf("expected %v to be %v", act, exp)
	}
}

func TestDisableEnableSecretVersion(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	payload := []byte("my-secret")
	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()
	version := testSecretVersion(t, secret.Name, payload)

	if err := disableSecretVersion(version.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	v, err := client.GetSecretVersion(ctx, &secretspb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if act, exp := v.State, secretspb.SecretVersion_DISABLED; act != exp {
		t.Errorf("expected %v to be %v", act, exp)
	}

	if err := enableSecretVersion(version.Name); err != nil {
		t.Fatal(err)
	}

	v, err = client.GetSecretVersion(ctx, &secretspb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if act, exp := v.State, secretspb.SecretVersion_ENABLED; act != exp {
		t.Errorf("expected %v to be %v", act, exp)
	}
}

func TestGetSecretVersion(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	payload := []byte("my-secret")
	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()
	version := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := getSecretVersion(&b, version.Name); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), "Found secret version"; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestGetSecret(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()

	var b bytes.Buffer
	if err := getSecret(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), "Found secret"; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestListSecretVersions(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	payload := []byte("my-secret")
	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()

	version1 := testSecretVersion(t, secret.Name, payload)
	version2 := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := listSecretVersions(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), fmt.Sprintf("%s with state ENABLED", version1.Name); !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}

	if act, exp := b.String(), fmt.Sprintf("%s with state ENABLED", version2.Name); !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestListSecrets(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	secret1, cleanup1 := testSecret(t, tc.ProjectID)
	defer cleanup1()

	secret2, cleanup2 := testSecret(t, tc.ProjectID)
	defer cleanup2()

	var b bytes.Buffer
	if err := listSecrets(&b, fmt.Sprintf("projects/%s", tc.ProjectID)); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), secret1.Name; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}

	if act, exp := b.String(), secret2.Name; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}
}

func TestUpdateSecret(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	secret, cleanup := testSecret(t, tc.ProjectID)
	defer cleanup()

	var b bytes.Buffer
	if err := updateSecret(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if act, exp := b.String(), "Updated secret"; !strings.Contains(act, exp) {
		t.Errorf("expected %q to contain %q", act, exp)
	}

	client, ctx := testClient(t)
	s, err := client.GetSecret(ctx, &secretspb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if act, exp := s.Labels, map[string]string{"secretmanager": "rocks"}; !reflect.DeepEqual(act, exp) {
		t.Errorf("expected %q to be %q", act, exp)
	}
}
