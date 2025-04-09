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

package modelarmor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"github.com/joho/godotenv"
	// modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
)

func testLocation(t *testing.T) string {
	t.Helper()

	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatalf(err.Error())
	}

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

func testClient(t *testing.T) (*modelarmor.Client, context.Context) {
	t.Helper()

	ctx := context.Background()

	locationId := testLocation(t)

	//Endpoint to send the request to regional server
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationId)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client, ctx
}

func testCleanupTemplate(t *testing.T, templateName string) {
	t.Helper()

	client, ctx := testClient(t)
	if err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
		}
	}

}

func TestGetProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)

	var b bytes.Buffer
	if _, err := getProjectFloorSettings(&b, tc.ProjectID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved floor setting:"; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestGetOrganizationFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatalf(err.Error())
	}

	organizationID := os.Getenv("GOLANG_SAMPLES_ORGANIZATION_ID")
	var b bytes.Buffer
	if _, err := getOrganizationFloorSettings(&b, organizationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved org floor setting:"; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestGetFolderFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatalf(err.Error())
	}

	folderID := os.Getenv("GOLANG_SAMPLES_FOLDER_ID")
	var b bytes.Buffer
	if _, err := getFolderFloorSettings(&b, folderID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestUpdateFolderFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatalf(err.Error())
	}
	folderID := os.Getenv("GOLANG_SAMPLES_FOLDER_ID")
	var b bytes.Buffer
	if _, err := updateFolderFloorSettings(&b, folderID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateFolderFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestUpdateOrganizationFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatalf(err.Error())
	}

	organizationID := os.Getenv("GOLANG_SAMPLES_ORGANIZATION_ID")
	var b bytes.Buffer
	if _, err := updateOrganizationFloorSettings(&b, organizationID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated org floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateOrganizationFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestUpdateProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)
	var b bytes.Buffer
	if _, err := updateProjectFloorSettings(&b, tc.ProjectID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated project floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateProjectFloorSettings: expected %q to contain %q", got, want)
	}
}
