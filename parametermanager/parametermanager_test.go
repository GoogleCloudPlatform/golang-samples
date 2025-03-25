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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
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

func testParameter(t *testing.T, projectID string, format parametermanagerpb.ParameterFormat) (*parametermanagerpb.Parameter, string) {
	parameterID := testName(t)

	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
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

func testParameterVersion(t *testing.T, projectID, parameterID, payload string) (*parametermanagerpb.ParameterVersion, string) {
	parameterVersionID := testName(t)

	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
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

func testCleanupParameter(t *testing.T, name string) {
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
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

func testCleanupParameterVersion(t *testing.T, name string) {
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
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

func TestDisableParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var b bytes.Buffer
	if err := disableParamVersion(&b, tc.ProjectID, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Disabled parameter version"; !strings.Contains(got, want) {
		t.Errorf("DisableParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestEnableParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var b bytes.Buffer
	if err := enableParamVersion(&b, tc.ProjectID, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Enabled parameter version"; !strings.Contains(got, want) {
		t.Errorf("EnableParameterVersion: expected %q to contain %q", got, want)
	}
}

func TestDeleteParam(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	defer testCleanupParameter(t, parameter.Name)

	var b bytes.Buffer
	if err := deleteParam(&b, tc.ProjectID, parameterID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Deleted parameter"; !strings.Contains(got, want) {
		t.Errorf("DeleteParameter: expected %q to contain %q", got, want)
	}
}

func TestDeleteParamVersion(t *testing.T) {
	tc := testutil.SystemTest(t)

	parameter, parameterID := testParameter(t, tc.ProjectID, parametermanagerpb.ParameterFormat_JSON)
	payload := `{"username": "test-user", "host": "localhost"}`
	parameterVersion, parameterVersionID := testParameterVersion(t, tc.ProjectID, parameterID, payload)
	defer testCleanupParameter(t, parameter.Name)
	defer testCleanupParameterVersion(t, parameterVersion.Name)

	var b bytes.Buffer
	if err := deleteParamVersion(&b, tc.ProjectID, parameterID, parameterVersionID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Deleted parameter version"; !strings.Contains(got, want) {
		t.Errorf("DeleteParameterVersion: expected %q to contain %q", got, want)
	}
}
