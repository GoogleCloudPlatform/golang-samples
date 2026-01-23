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
	"path"
	"reflect"
	"strings"
	"time"

	"testing"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	regional_secretmanager "github.com/GoogleCloudPlatform/golang-samples/secretmanager/regional_samples"
	"github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	grpcstatus "google.golang.org/grpc/status"
)

func testLocation(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_REGIONAL_SAMPLES_LOCATION")
	if v == "" {
		tb.Skip("testIamUser: missing GOLANG_REGIONAL_SAMPLES_LOCATION")
	}

	return v
}

func testClient(tb testing.TB) (*secretmanager.Client, context.Context) {
	tb.Helper()

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		tb.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testRegionalClient(tb testing.TB) (*secretmanager.Client, context.Context) {
	tb.Helper()

	ctx := context.Background()

	locationId := testLocation(tb)

	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		tb.Fatalf("testRegionalClient: failed to create regional client: %v", err)
	}
	return client, ctx
}

func testResourceManagerTagsKeyClient(tb testing.TB) (*resourcemanager.TagKeysClient, context.Context) {
	tb.Helper()
	ctx := context.Background()

	client, err := resourcemanager.NewTagKeysClient(ctx)
	if err != nil {
		tb.Fatalf("testResourceManagerTagsKeyClient: failed to create client: %v", err)
	}
	return client, ctx

}

func testResourceManagerTagsValueClient(tb testing.TB) (*resourcemanager.TagValuesClient, context.Context) {
	tb.Helper()
	ctx := context.Background()

	client, err := resourcemanager.NewTagValuesClient(ctx)
	if err != nil {
		tb.Fatalf("testResourceManagerTagsValueClient: failed to create client: %v", err)
	}
	return client, ctx

}

func testKmsKey(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_SAMPLES_KMS_KEY")
	if v == "" {
		tb.Skip("testKmsKey: missing GOLANG_SAMPLES_KMS_KEY")
	}

	return v
}

func testResourceManagerTagBindingsClient(tb testing.TB) (*resourcemanager.TagBindingsClient, context.Context) {
	tb.Helper()
	ctx := context.Background()

	client, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		tb.Fatalf("testResourceManagerTagBindingsClient: failed to create client: %v", err)
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

func testTopic(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_SAMPLES_TOPIC_NAME")
	if v == "" {
		tb.Skip("testTopic: missing GOLANG_SAMPLES_TOPIC_NAME")
	}

	return v
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
			Labels: map[string]string{
				"labelkey": "labelvalue",
			},
			Annotations: map[string]string{
				"annotationkey": "annotationvalue",
			},
		},
	})
	if err != nil {
		tb.Fatalf("testSecret: failed to create secret: %v", err)
	}

	return secret
}

func testRegionalSecret(tb testing.TB, projectID string) (*secretmanagerpb.Secret, string) {
	tb.Helper()

	secretID := testName(tb)

	locationID := testLocation(tb)
	client, ctx := testRegionalClient(tb)
	secret, err := client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, locationID),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Annotations: map[string]string{
				"annotationkey": "annotationvalue",
			},
			Labels: map[string]string{
				"labelkey": "labelvalue",
			},
		},
	})
	if err != nil {
		tb.Fatalf("testSecret: failed to create secret: %v", err)
	}

	return secret, secretID
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

func testRegionalSecretVersion(tb testing.TB, parent string, payload []byte) *secretmanagerpb.SecretVersion {
	tb.Helper()

	client, ctx := testRegionalClient(tb)

	version, err := client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	})
	if err != nil {
		tb.Fatalf("testSecretVersion: failed to create regional secret version: %v", err)
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

func testCleanupRegionalSecret(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testRegionalClient(tb)

	if err := client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			tb.Fatalf("testCleanupSecret: failed to delete regional secret: %v", err)
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

func TestAccessRegionalSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	testRegionalSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := regional_secretmanager.AccessRegionalSecretVersion(&b, tc.ProjectID, locationID, secretID, "1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), string(payload); !strings.Contains(got, want) {
		t.Errorf("accessRegionalSecretVersion: expected %q to contain %q", got, want)
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

func TestAddRegionalSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.AddRegionalSecretVersion(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Added regional secret version:"; !strings.Contains(got, want) {
		t.Errorf("addSecretVersion: expected %q to contain %q", got, want)
	}
}

func TestConsumeEventNotification(t *testing.T) {
	v, err := ConsumeEventNotification(context.Background(), PubSubMessage{
		Attributes: PubSubAttributes{
			SecretId:  "projects/p/secrets/s",
			EventType: "SECRET_UPDATE",
		},
		Data: []byte("hello!"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := v, `Received SECRET_UPDATE for projects/p/secrets/s. New metadata: "hello!".`; !strings.Contains(got, want) {
		t.Errorf("consumeEventNotification: expected %q to contain %q", got, want)
	}
}

func TestCreateSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createSecret"

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)

	var b bytes.Buffer
	if err := createSecret(&b, parent, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	if got, want := b.String(), "Created secret:"; !strings.Contains(got, want) {
		t.Errorf("createSecret: expected %q to contain %q", got, want)
	}
}

func TestCreateSecretWithTTL(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createSecretTTL"

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)

	duration := time.Second * 70

	var b bytes.Buffer
	if err := createSecretWithTTL(&b, parent, secretID, duration); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	if got, want := b.String(), "Created secret with ttl:"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithTTL: expected %q to contain %q", got, want)
	}
}

func TestCreateSecretWithLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createSecretWithLabels"

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)

	var b bytes.Buffer
	if err := createSecretWithLabels(&b, parent, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	if got, want := b.String(), "Created secret with labels:"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithLabels: expected %q to contain %q", got, want)
	}
}

func TestCreateSecretWithAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createSecretWithAnnotations"

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)

	var b bytes.Buffer
	if err := createSecretWithAnnotations(&b, parent, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	if got, want := b.String(), "Created secret with annotations:"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithAnnotations: expected %q to contain %q", got, want)
	}
}

func TestCreateRegionalSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createRegionalSecret"
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.CreateRegionalSecret(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID))

	if got, want := b.String(), "Created regional secret:"; !strings.Contains(got, want) {
		t.Errorf("createSecret: expected %q to contain %q", got, want)
	}
}

func TestCreateUserManagedReplicationSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createUmmrSecret"
	locations := []string{"us-east1", "us-east4", "us-west1"}

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)

	var b bytes.Buffer
	if err := createUserManagedReplicationSecret(&b, parent, secretID, locations); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	if got, want := b.String(), "Created secret with user managed replication:"; !strings.Contains(got, want) {
		t.Errorf("createUserManagedReplicationSecret: expected %q to contain %q", got, want)
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

func TestDeleteSecretLabel(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := deleteSecretLabel(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	client, ctx := testClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{}; reflect.DeepEqual(got, want) {
		t.Errorf("deleteSecretLabel: expected %q to be %q", got, want)
	}
}

func TestDeleteRegionalSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretId := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	if err := regional_secretmanager.DeleteRegionalSecret(tc.ProjectID, locationID, secretId); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	_, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("deleteRegionalSecret: expected %v to be not found", err)
	}

}

func TestDeleteSecretWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	if err := deleteSecretWithEtag(secret.Name, secret.Etag); err != nil {
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

func TestDeleteRegionalSecretWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretId := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	if err := regional_secretmanager.DeleteRegionalSecretWithEtag(tc.ProjectID, locationID, secretId, secret.Etag); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	_, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
		t.Errorf("deleteRegionalSecret: expected %v to be not found", err)
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

func TestDestroyRegionalSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	version := testRegionalSecretVersion(t, secret.Name, payload)

	if err := regional_secretmanager.DestroyRegionalSecretVersion(tc.ProjectID, locationID, secretID, "1"); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	v, err := client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_DESTROYED; got != want {
		t.Errorf("testRegionalSecretVersion: expected %v to be %v", got, want)
	}
}

func TestDestroySecretVersionWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)
	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version := testSecretVersion(t, secret.Name, payload)

	if err := destroySecretVersionWithEtag(version.Name, version.Etag); err != nil {
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

func TestDestroyRegionalSecretVersionWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)
	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	version := testRegionalSecretVersion(t, secret.Name, payload)

	if err := regional_secretmanager.DestroyRegionalSecretVersionWithEtag(tc.ProjectID, locationID, secretID, "1", version.Etag); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	v, err := client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_DESTROYED; got != want {
		t.Errorf("testRegionalSecretVersion: expected %v to be %v", got, want)
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

func TestDisableEnableRegionalSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	version := testRegionalSecretVersion(t, secret.Name, payload)

	if err := regional_secretmanager.DisableRegionalSecretVersion(tc.ProjectID, locationID, secretID, "1"); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	v, err := client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_DISABLED; got != want {
		t.Errorf("testRegionalSecretVersion: expected %v to be %v", got, want)
	}

	if err := regional_secretmanager.EnableRegionalSecretVersion(tc.ProjectID, locationID, secretID, "1"); err != nil {
		t.Fatal(err)
	}

	v, err = client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_ENABLED; got != want {
		t.Errorf("testRegionalSecretVersion: expected %v to be %v", got, want)
	}
}

func TestDisableEnableSecretVersionWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version := testSecretVersion(t, secret.Name, payload)

	if err := disableSecretVersionWithEtag(version.Name, version.Etag); err != nil {
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

	if err := enableSecretVersionWithEtag(version.Name, v.Etag); err != nil {
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

func TestDisableEnableRegionalSecretVersionWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationId := testLocation(t)

	version := testRegionalSecretVersion(t, secret.Name, payload)

	if err := regional_secretmanager.DisableRegionalSecretVersionWithEtag(tc.ProjectID, locationId, secretID, "1", version.Etag); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	v, err := client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_DISABLED; got != want {
		t.Errorf("testRegionalSecretVersion: expected %v to be %v", got, want)
	}

	if err := regional_secretmanager.EnableRegionalSecretVersionWithEtag(tc.ProjectID, locationId, secretID, "1", v.Etag); err != nil {
		t.Fatal(err)
	}

	v, err = client.GetSecretVersion(ctx, &secretmanagerpb.GetSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v.State, secretmanagerpb.SecretVersion_ENABLED; got != want {
		t.Errorf("testRegionalSecretVersion: expected %v to be %v", got, want)
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

func TestGetRegionalSecretVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	testRegionalSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := regional_secretmanager.GetRegionalSecretVersion(&b, tc.ProjectID, locationID, secretID, "1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found regional secret version"; !strings.Contains(got, want) {
		t.Errorf("testRegionalSecretVersion: expected %q to contain %q", got, want)
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

func TestGetRegionalSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretdID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.GetRegionalSecret(&b, tc.ProjectID, locationID, secretdID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found regional secret"; !strings.Contains(got, want) {
		t.Errorf("getRegionalSecret: expected %q to contain %q", got, want)
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

func TestIamGrantAccessWithRegionalSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	iamUser := testIamUser(t)

	var b bytes.Buffer
	if err := regional_secretmanager.IamGrantAccessWithRegionalSecret(&b, tc.ProjectID, locationID, secretID, iamUser); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated IAM policy"; !strings.Contains(got, want) {
		t.Errorf("getRegionalSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
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

func TestIamRevokeAccessWithRegionalSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	iamUser := testIamUser(t)

	var b bytes.Buffer
	if err := regional_secretmanager.IamRevokeAccessWithRegionalSecret(&b, tc.ProjectID, locationID, secretID, iamUser); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated IAM policy"; !strings.Contains(got, want) {
		t.Errorf("getRegionalSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
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

func TestListRegionalSecretVersions(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	version1 := testRegionalSecretVersion(t, secret.Name, payload)
	version2 := testRegionalSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := regional_secretmanager.ListRegionalSecretVersions(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), fmt.Sprintf("%s with state ENABLED", version1.Name); !strings.Contains(got, want) {
		t.Errorf("listSecretVersions: expected %q to contain %q", got, want)
	}

	if got, want := b.String(), fmt.Sprintf("%s with state ENABLED", version2.Name); !strings.Contains(got, want) {
		t.Errorf("listSecretVersions: expected %q to contain %q", got, want)
	}
}

func TestListSecretVersionsWithFilter(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	version1 := testSecretVersion(t, secret.Name, payload)
	version2 := testSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := listSecretVersionsWithFilter(&b, secret.Name, fmt.Sprintf("name:%s", version1.Name)); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), fmt.Sprintf("%s with state ENABLED", version1.Name); !strings.Contains(got, want) {
		t.Errorf("listSecretVersions: expected %q to contain %q", got, want)
	}

	if got, lacked := b.String(), fmt.Sprintf("%s with state ENABLED", version2.Name); strings.Contains(got, lacked) {
		t.Errorf("listSecretVersions: expected %q to not contain %q", got, lacked)
	}
}

func TestListRegionalSecretVersionsWithFilter(t *testing.T) {
	tc := testutil.SystemTest(t)

	payload := []byte("my-secret")
	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	version1 := testRegionalSecretVersion(t, secret.Name, payload)
	version2 := testRegionalSecretVersion(t, secret.Name, payload)

	var b bytes.Buffer
	if err := regional_secretmanager.ListRegionalSecretVersionsWithFilter(&b, tc.ProjectID, locationID, secretID, fmt.Sprintf("name:%s", version1.Name)); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), fmt.Sprintf("%s with state ENABLED", version1.Name); !strings.Contains(got, want) {
		t.Errorf("listSecretVersions: expected %q to contain %q", got, want)
	}

	if got, lacked := b.String(), fmt.Sprintf("%s with state ENABLED", version2.Name); strings.Contains(got, lacked) {
		t.Errorf("listSecretVersions: expected %q to not contain %q", got, lacked)
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

func TestViewSecretLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := viewSecretLabels(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret"; !strings.Contains(got, want) {
		t.Errorf("viewSecretLabels: expected %q to contain %q", got, want)
	}

	client, ctx := testClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{"labelkey": "labelvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("viewSecretLabels: expected %q to be %q", got, want)
	}
}

func TestViewSecretAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := viewSecretAnnotations(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret"; !strings.Contains(got, want) {
		t.Errorf("viewSecretAnnotations: expected %q to contain %q", got, want)
	}

	client, ctx := testClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Annotations, map[string]string{"annotationkey": "annotationvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("viewSecretAnnotations: expected %q to be %q", got, want)
	}
}

func TestListRegionalSecrets(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret1, _ := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret1.Name)

	secret2, _ := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret2.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.ListRegionalSecrets(&b, tc.ProjectID, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), secret1.Name; !strings.Contains(got, want) {
		t.Errorf("listRegionalSecrets: expected %q to contain %q", got, want)
	}

	if got, want := b.String(), secret2.Name; !strings.Contains(got, want) {
		t.Errorf("listRegionalSecrets: expected %q to contain %q", got, want)
	}
}

func TestListSecretsWithFilter(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret1 := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret1.Name)

	secret2 := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret2.Name)

	var b bytes.Buffer
	if err := listSecretsWithFilter(&b, fmt.Sprintf("projects/%s", tc.ProjectID), fmt.Sprintf("name:%s", secret1.Name)); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), secret1.Name; !strings.Contains(got, want) {
		t.Errorf("listSecrets: expected %q to contain %q", got, want)
	}

	if got, lacked := b.String(), secret2.Name; strings.Contains(got, lacked) {
		t.Errorf("listSecrets: expected %q to not contain %q", got, lacked)
	}
}

func TestListRegionalSecretsWithFilter(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret1, _ := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret1.Name)

	secret2, _ := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret2.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.ListRegionalSecretsWithFilter(&b, tc.ProjectID, locationID, fmt.Sprintf("name:%s", secret1.Name)); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), secret1.Name; !strings.Contains(got, want) {
		t.Errorf("listRegionalSecrets: expected %q to contain %q", got, want)
	}

	if got, lacked := b.String(), secret2.Name; strings.Contains(got, lacked) {
		t.Errorf("listRegionalSecrets: expected %q to not contain %q", got, lacked)
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

func TestCreateUpdateSecretLabel(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := createUpdateSecretLabel(&b, secret.Name); err != nil {
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

	if got, want := s.Labels, map[string]string{"labelkey": "updatedlabelvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("createUpdateSecretLabel: expected %q to be %q", got, want)
	}
}

func TestEditSecretAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := editSecretAnnotation(&b, secret.Name); err != nil {
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

	if got, want := s.Annotations, map[string]string{"annotationkey": "updatedannotationvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("editSecretAnnotation: expected %q to be %q", got, want)
	}
}

func TestRegionalUpdateSecret(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.UpdateRegionalSecret(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated regional secret"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)

	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{"secretmanager": "rocks"}; !reflect.DeepEqual(got, want) {
		t.Errorf("updateRegionalSecret: expected %q to be %q", got, want)
	}
}

func TestUpdateSecretWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	var b bytes.Buffer
	if err := updateSecretWithEtag(&b, secret.Name, secret.Etag); err != nil {
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

func TestUpdateRegionalSecretWithEtag(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := regional_secretmanager.UpdateRegionalSecretWithEtag(&b, tc.ProjectID, locationID, secretID, secret.Etag); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated regional secret"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{"secretmanager": "rocks"}; !reflect.DeepEqual(got, want) {
		t.Errorf("updateRegionalSecret: expected %q to be %q", got, want)
	}
}

func TestUpdateSecretWithAlias(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)

	testSecretVersion(t, secret.Name, []byte("my-secret"))

	var b bytes.Buffer
	if err := updateSecretWithAlias(&b, secret.Name); err != nil {
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

	if got, want := s.VersionAliases, map[string]int64{"test": 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("updateSecret: expected %q to be %q", got, want)
	}
}

func TestUpdateRegionalSecretWithAlias(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	testRegionalSecretVersion(t, secret.Name, []byte("my-secret"))

	var b bytes.Buffer
	if err := regional_secretmanager.UpdateRegionalSecretWithAlias(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated regional secret"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.VersionAliases, map[string]int64{"test": 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("updateRegionalSecret: expected %q to be %q", got, want)
	}
}

func testCreateTagKey(tb testing.TB, projectID string) *resourcemanagerpb.TagKey {
	tb.Helper()

	client, ctx := testResourceManagerTagsKeyClient(tb)
	parent := fmt.Sprintf("projects/%s", projectID)
	tagKeyName := testName(tb)
	tagKeyDescription := "creating tag key for secretmanager tags sample"

	tagKeyOperation, err := client.CreateTagKey(ctx, &resourcemanagerpb.CreateTagKeyRequest{
		TagKey: &resourcemanagerpb.TagKey{
			Parent:      parent,
			ShortName:   tagKeyName,
			Description: tagKeyDescription,
		},
	})
	if err != nil {
		tb.Fatalf("testCreateTagKey: failed to create tagKey: %v", err)
	}

	createdTagKey, err := tagKeyOperation.Wait(ctx)
	if err != nil {
		tb.Fatalf("testCreateTagKey: failed to create TagKey after waiting for operation: %v", err)
	}

	return createdTagKey
}

func testCreateTagValue(tb testing.TB, tagKeyId string) *resourcemanagerpb.TagValue {
	tb.Helper()

	client, ctx := testResourceManagerTagsValueClient(tb)
	tagValueName := testName(tb)
	tagKeyDescription := "creating TagValue for secretmanager tags sample"

	tagKeyOperation, err := client.CreateTagValue(ctx, &resourcemanagerpb.CreateTagValueRequest{
		TagValue: &resourcemanagerpb.TagValue{
			Parent:      tagKeyId,
			ShortName:   tagValueName,
			Description: tagKeyDescription,
		},
	})
	if err != nil {
		tb.Fatalf("testCreateTagValue: failed to create tagValue: %v", err)
	}

	createdTagValue, err := tagKeyOperation.Wait(ctx)
	if err != nil {
		tb.Fatalf("testCreateTagValue: failed to create TagValue after waiting for operation: %v", err)
	}

	return createdTagValue
}

func testCleanupTagKey(tb testing.TB, tagKeyName string) {
	tb.Helper()

	client, ctx := testResourceManagerTagsKeyClient(tb)

	tagKeyOperation, err := client.DeleteTagKey(ctx, &resourcemanagerpb.DeleteTagKeyRequest{
		Name: tagKeyName,
	})
	if err != nil {
		tb.Fatalf("testCleanupTagKey: failed to delete tagKey: %v", err)
		return
	}

	_, err = tagKeyOperation.Wait(ctx)
	if err != nil {
		tb.Fatalf("testCleanupTagKey: failed to delete TagKey after waiting for operation: %v", err)
	}
}

// Polling to clean up the tag value because, after deleting a secret, it takes some time for the tag value to become unbound.
func testCleanupTagValue(tb testing.TB, tagValueName string) {
	tb.Helper()

	client, ctx := testResourceManagerTagsValueClient(tb)

	maxPollingDuration := 10 * time.Minute
	initialDelay := 2 * time.Second
	maxBackoffDelay := 30 * time.Second

	startTime := time.Now()
	attempt := 0

	for time.Since(startTime) < maxPollingDuration {
		attempt++

		tagValueOperation, err := client.DeleteTagValue(ctx, &resourcemanagerpb.DeleteTagValueRequest{
			Name: tagValueName,
		})

		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.NotFound {
				tb.Logf("Tag value %s already deleted (or never existed) after %v.", tagValueName, time.Since(startTime))
				return
			}

			if ok && s.Code() == codes.FailedPrecondition && strings.Contains(s.Message(), "attached to resources") {
				delay := initialDelay * time.Duration(1<<uint(attempt-1))
				if delay > maxBackoffDelay {
					delay = maxBackoffDelay
				}
				time.Sleep(delay)
				continue
			}

			tb.Errorf("testCleanupTagValue: failed to initiate delete for tag value %s due to unrecoverable error: %v", tagValueName, err)
			return
		}

		_, err = tagValueOperation.Wait(ctx)
		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.NotFound {
				tb.Logf("Tag value %s was deleted during operation wait (or already gone).", tagValueName)
				return
			}
			if ok && s.Code() == codes.FailedPrecondition && strings.Contains(s.Message(), "attached to resources") {
				delay := initialDelay * time.Duration(1<<uint(attempt-1))
				if delay > maxBackoffDelay {
					delay = maxBackoffDelay
				}
				time.Sleep(delay)
				continue
			}

			tb.Errorf("testCleanupTagValue: failed to delete tag value %s after waiting for operation due to unrecoverable error: %v", tagValueName, err)
			return
		}

		tb.Logf("Successfully deleted tag value %s after %v (attempt %d).", tagValueName, time.Since(startTime), attempt)
		return
	}
	tb.Errorf("testCleanupTagValue: failed to delete tag value %s after %v (max duration reached). It might still be attached.", tagValueName, maxPollingDuration)
}

func TestCreateSecretWithTags(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createSecretWithTags"

	parent := fmt.Sprintf("projects/%s", tc.ProjectID)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.Name)
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.Name)

	t.Logf("Secret ID used: %s", secretID)
	t.Logf("Tag Key used: %s", tagKey.GetName())
	t.Logf("Tag Value used: %s", tagValue.Name)

	var b bytes.Buffer
	if err := createSecretWithTags(&b, parent, secretID, tagKey.GetName(), tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	defer testCleanupSecret(t, fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID))

	if got, want := b.String(), "Created secret with tags:"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithTags: expected %q to contain %q", got, want)
	}

}

func TestCreateSecretWithCMEK(t *testing.T) {
	tc := testutil.SystemTest(t)
	kmsKeyName := testKmsKey(t)

	secretId := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretId)
	defer testCleanupSecret(t, secretName)

	var b bytes.Buffer
	if err := createSecretWithCMEK(&b, tc.ProjectID, secretId, kmsKeyName); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	if !strings.Contains(got, secretId) {
		t.Errorf("createSecretWithCMEK: output %q did not contain secretId %q", got, secretId)
	}
	if !strings.Contains(got, kmsKeyName) {
		t.Errorf("createSecretWithCMEK: output %q did not contain kmsKeyName %q", got, kmsKeyName)
	}

	// Verify CMEK key with GetSecret.
	client, ctx := testClient(t)
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetReplication().GetAutomatic().GetCustomerManagedEncryption().GetKmsKeyName() != kmsKeyName {
		t.Errorf("CMEK key name mismatch: got %q, want %q", secret.GetReplication().GetAutomatic().GetCustomerManagedEncryption().GetKmsKeyName(), kmsKeyName)
	}
}

func TestCreateSecretWithDelayedDestroy(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	ttl := 24 * time.Hour

	var b bytes.Buffer
	if err := createSecretWithDelayedDestroy(&b, tc.ProjectID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithDelayedDestroy: expected %q to contain %q", got, want)
	}

	// Verify Delayed Destroy with GetSecret.
	client, ctx := testClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetVersionDestroyTtl() == nil {
		t.Fatal("GetSecret: VersionDestroyTtl is nil, expected non-nil")
	}
	if secret.GetVersionDestroyTtl().AsDuration() != ttl {
		t.Errorf("VersionDestroyTtl mismatch: got %v, want %v", secret.GetVersionDestroyTtl().AsDuration(), ttl)
	}
}

func TestUpdateSecretWithDelayedDestroy(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	client, ctx := testClient(t)

	var b bytes.Buffer
	if err := createSecretWithDelayedDestroy(&b, tc.ProjectID, secretID); err != nil {
		t.Fatal(err)
	}

	newDelayedDestroy := 48 * time.Hour
	if err := updateSecretWithDelayedDestroy(&b, secretName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateSecretWithDelayedDestroy: expected %q to contain %q", got, want)
	}

	// Verify TTL with GetSecret.
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetVersionDestroyTtl() == nil {
		t.Fatal("GetSecret: VersionDestroyTtl is nil, expected non-nil")
	}
	if secret.GetVersionDestroyTtl().AsDuration() != newDelayedDestroy {
		t.Errorf("VersionDestroyTtl mismatch: got %v, want %v", secret.GetVersionDestroyTtl().AsDuration(), newDelayedDestroy)
	}
}

func TestDeleteSecretVersionDestroyTTL(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	client, ctx := testClient(t)

	var b bytes.Buffer
	if err := createSecretWithDelayedDestroy(&b, tc.ProjectID, secretID); err != nil {
		t.Fatal(err)
	}

	if err := deleteSecretVersionDestroyTTL(&b, secretName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "removed version_destroy_ttl"; !strings.Contains(got, want) {
		t.Errorf("deleteSecretVersionDestroyTTL: expected %q to contain %q", got, want)
	}

	// Verify TTL with GetSecret.
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetVersionDestroyTtl() != nil {
		t.Fatal("GetSecret: VersionDestroyTtl is not nil, expected nil")
	}
}

func TestCreateSecretWithRotation(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	topicName := testTopic(t)
	client, ctx := testClient(t)

	rotationPeriod := 24 * time.Hour

	var b bytes.Buffer
	if err := createSecretWithRotation(&b, tc.ProjectID, secretID, topicName); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithRotation: expected %q to contain %q", got, want)
	}

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})

	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetRotation() == nil {
		t.Fatal("GetSecret: Rotation is nil, expected non-nil")
	}
	if secret.GetRotation().GetRotationPeriod().AsDuration() != rotationPeriod {
		t.Errorf("RotationPeriod mismatch: got %v, want %v", secret.GetRotation().GetRotationPeriod().AsDuration(), rotationPeriod)
	}
	if secret.GetRotation().GetNextRotationTime() == nil {
		t.Fatal("GetSecret: NextRotationTime is nil, expected non-nil")
	}
}

func TestUpdateSecretRotationPeriod(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	topicName := testTopic(t)
	client, ctx := testClient(t)

	var b bytes.Buffer
	if err := createSecretWithRotation(&b, tc.ProjectID, secretID, topicName); err != nil {
		t.Fatal(err)
	}

	// Update rotation period.
	updatedRotationPeriod := 48 * time.Hour
	b.Reset()
	if err := updateSecretRotationPeriod(&b, secretName); err != nil {
		t.Fatal(err)
	}

	got := b.String()
	if !strings.Contains(got, secretID) {
		t.Errorf("updateSecretRotationPeriod: output %q did not contain secretID %q", got, secretID)
	}

	// Verify rotation period with GetSecret.
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetRotation() == nil {
		t.Fatal("GetSecret: Rotation is nil, expected non-nil")
	}
	if secret.GetRotation().GetRotationPeriod().AsDuration() != updatedRotationPeriod {
		t.Errorf("RotationPeriod mismatch: got %v, want %v", secret.GetRotation().GetRotationPeriod().AsDuration(), updatedRotationPeriod)
	}
}

func TestDeleteSecretRotation(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	topicName := testTopic(t)
	client, ctx := testClient(t)

	var b bytes.Buffer
	if err := createSecretWithRotation(&b, tc.ProjectID, secretID, topicName); err != nil {
		t.Fatal(err)
	}

	// Remove rotation.
	if err := deleteSecretRotation(&b, secretName); err != nil {
		t.Fatal(err)
	}

	got := b.String()
	if !strings.Contains(got, secretID) {
		t.Errorf("deleteSecretRotation: output %q did not contain secretID %q", got, secretID)
	}

	// Verify rotation is removed with GetSecret.
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetRotation() != nil {
		t.Errorf("Rotation mismatch: got %v, want nil", secret.GetRotation())
	}
}

func TestCreateSecretWithTopic(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	topicName := testTopic(t)
	client, ctx := testClient(t)

	var b bytes.Buffer
	if err := createSecretWithTopic(&b, tc.ProjectID, secretID, topicName); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithTopic: expected %q to contain %q", got, want)
	}

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})

	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if len(secret.GetTopics()) != 1 || secret.GetTopics()[0].GetName() != topicName {
		t.Errorf("Topics mismatch: got %v, want %s", secret.GetTopics(), topicName)
	}
}

func TestCreateSecretWithExpireTime(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	// Expire time in 1 hour.
	expire := time.Now().Add(time.Hour)

	var b bytes.Buffer
	if err := createSecretWithExpireTime(&b, tc.ProjectID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createSecretWithExpireTime: expected %q to contain %q", got, want)
	}

	// Verify expire time with GetSecret.
	client, ctx := testClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetExpireTime() == nil {
		t.Fatal("GetSecret: ExpireTime is nil, expected non-nil")
	}

	if diff := secret.GetExpireTime().AsTime().Unix() - expire.Unix(); diff > 1 || diff < -1 {
		t.Errorf("ExpireTime mismatch: got %v, want %v", secret.GetExpireTime().AsTime(), expire)
	}
}

func TestUpdateSecretExpiration(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	// Create with expire time in 1 hour.
	var b bytes.Buffer
	if err := createSecretWithExpireTime(&b, tc.ProjectID, secretID); err != nil {
		t.Fatal(err)
	}

	// Update expire time to 2 hours.
	newExpire := time.Now().Add(2 * time.Hour)
	if err := updateSecretExpiration(&b, secretName); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateSecretExpiration: expected %q to contain %q", got, want)
	}

	// Verify expire time with GetSecret.
	client, ctx := testClient(t)
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetExpireTime() == nil {
		t.Fatal("GetSecret: ExpireTime is nil, expected non-nil")
	}

	if diff := secret.GetExpireTime().AsTime().Unix() - newExpire.Unix(); diff > 1 || diff < -1 {
		t.Errorf("ExpireTime mismatch: got %v, want %v", secret.GetExpireTime().AsTime(), newExpire)
	}
}

func TestRemoveExpiration(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	var b bytes.Buffer
	if err := createSecretWithExpireTime(&b, tc.ProjectID, secretID); err != nil {
		t.Fatal(err)
	}

	// Remove expire time.
	b.Reset()
	if err := deleteExpiration(&b, secretName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Removed expiration"; !strings.Contains(got, want) {
		t.Errorf("deleteExpiration: expected %q to contain %q", got, want)
	}

	// Verify expire time is removed with GetSecret.
	client, ctx := testClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetExpireTime() != nil {
		t.Errorf("GetSecret: ExpireTime is %v, expected nil", secret.GetExpireTime())
	}
}

func TestBindTagsToSecret(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.GetName())
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.GetName())

	var b bytes.Buffer
	if err := bindTagsToSecret(&b, tc.ProjectID, secretID, tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	if got, want := b.String(), "Tag binding created"; !strings.Contains(got, want) {
		t.Errorf("bindTagsToSecret: expected %q to contain %q", got, want)
	}

	// Verify binding exists with API
	ctx := context.Background()
	tagBindingsClient, ctx := testResourceManagerTagBindingsClient(t)
	defer tagBindingsClient.Close()

	parent := "//secretmanager.googleapis.com/" + secretName
	it := tagBindingsClient.ListTagBindings(ctx, &resourcemanagerpb.ListTagBindingsRequest{
		Parent: parent,
	})

	found := false
	for {
		binding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Failed to list tag bindings for verification: %v", err)
		}
		if binding.TagValue == tagValue.GetName() && path.Base(binding.GetParent()) == path.Base(secretName) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Tag binding for %s with value %s not found after creation", secretName, tagValue.GetName())
	}
}

func TestListTagBindings(t *testing.T) {
	tc := testutil.SystemTest(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.GetName())
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.GetName())

	secretID := testName(t)

	var b bytes.Buffer
	if err := bindTagsToSecret(&b, tc.ProjectID, secretID, tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	b.Reset()
	if err := listTagBindings(&b, secretName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), tagValue.GetName(); !strings.Contains(got, want) {
		t.Errorf("listTagBindings: expected %q to contain %q", got, want)
	}
}

func TestDetachTag(t *testing.T) {
	tc := testutil.SystemTest(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.GetName())
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.GetName())

	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/secrets/%s", tc.ProjectID, secretID)
	defer testCleanupSecret(t, secretName)

	var b bytes.Buffer
	if err := bindTagsToSecret(&b, tc.ProjectID, secretID, tagValue.GetName()); err != nil {
		t.Fatal(err)
	}

	b.Reset()
	if err := detachTag(&b, secretName, tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Detached tag value"; !strings.Contains(got, want) {
		t.Errorf("detachTag: expected %q to contain %q", got, want)
	}

	b.Reset()
	if err := listTagBindings(&b, secretName); err != nil {
		t.Fatal(err)
	}
	if got, dontwant := b.String(), tagValue.GetName(); strings.Contains(got, dontwant) {
		t.Errorf("listTagBindings after detach: expected %q not to contain %q", got, dontwant)
	}
}

func TestDeleteSecretAnnotation(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret := testSecret(t, tc.ProjectID)
	defer testCleanupSecret(t, secret.Name)
	annotationKey := "annotationkey"

	var b bytes.Buffer

	if err := deleteSecretAnnotation(&b, secret.Name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Deleted annotation"; !strings.Contains(got, want) {
		t.Errorf("deleteSecretAnnotation: expected %q to contain %q", got, want)
	}

	client, ctx := testClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := s.Annotations[annotationKey]; ok {
		t.Errorf("deleteSecretAnnotation: key %q still present after deletion", annotationKey)
	}
}
