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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func testLocation(t *testing.T) string {
	t.Helper()

	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatal(err.Error())
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

func TestUpdateTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var b bytes.Buffer
	if _, err := createModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, "us-central1", templateID))

	if _, err := updateModelArmorTemplate(&b, tc.ProjectID, "us-central1", templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated Filter Config: "; !strings.Contains(got, want) {
		t.Errorf("updateModelArmorTemplate: expected %q to contain %q", got, want)
	}
}
