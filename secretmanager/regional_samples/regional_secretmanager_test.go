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

package regional_secretmanager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"testing"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
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

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
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
		},
	})
	if err != nil {
		tb.Fatalf("testSecret: failed to create secret: %v", err)
	}

	return secret, secretID
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

func TestCreateRegionalSecretWithAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createRegionalSecretWithAnnotations"
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithAnnotations(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID))

	if got, want := b.String(), "Created secret with annotations:"; !strings.Contains(got, want) {
		t.Errorf("createSecret: expected %q to contain %q", got, want)
	}
}

func TestEditRegionalSecretAnnotation(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := editRegionalSecretAnnotation(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)

	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Annotations, map[string]string{"annotationkey": "updatedannotationvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("editSecretAnnotations: expected %q to be %q", got, want)
	}
}

func TestViewRegionalSecretAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := viewRegionalSecretAnnotations(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret"; !strings.Contains(got, want) {
		t.Errorf("viewRegionalSecretAnnotations: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Annotations, map[string]string{"annotationkey": "annotationvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("viewRegionalSecretAnnotations: expected %q to be %q", got, want)
	}
}
