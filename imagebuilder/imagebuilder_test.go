// Copyright 2026 Google LLC
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

package imagebuilder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/googleapi"
)

// testRegion is the region that the samples use.
const testRegion = "us-central1"

// TestSubmitBuild submits a build to the test project, and then cancels it,
// because a full Image Builder pipeline takes up to an hour and creates Compute
// Engine resources.
func TestSubmitBuild(t *testing.T) {
	tc := testutil.SystemTest(t)
	serviceAccount := os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
	if serviceAccount == "" {
		t.Skip("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL not set")
	}
	// Cloud Build doesn't check these paths when it accepts a build, so they
	// don't need to exist.
	gcsWorkdir := "gs://" + tc.ProjectID + "-imagebuilder-test/workdir/"
	configPath := "gs://" + tc.ProjectID + "-imagebuilder-test/imagebuilder.yaml"

	var buf bytes.Buffer
	if err := submitBuild(&buf, tc.ProjectID, serviceAccount, gcsWorkdir, configPath); err != nil {
		t.Fatalf("submitBuild: %v", err)
	}
	got := buf.String()
	m := regexp.MustCompile(`Build submitted\. ID: (\S+)`).FindStringSubmatch(got)
	if m == nil {
		t.Fatalf("submitBuild output = %q, want a build ID", got)
	}

	// Cancel the build so that the pipeline doesn't run.
	ctx := context.Background()
	cbService, err := cloudbuild.NewService(ctx)
	if err != nil {
		t.Fatalf("cloudbuild.NewService: %v", err)
	}
	name := fmt.Sprintf("projects/%s/locations/%s/builds/%s", tc.ProjectID, testRegion, m[1])
	req := &cloudbuild.CancelBuildRequest{}
	if _, err := cbService.Projects.Locations.Builds.Cancel(name, req).Context(ctx).Do(); err != nil {
		t.Errorf("Builds.Cancel(%q): %v", name, err)
	}
}

// TestScheduleBuild uses a project that doesn't exist, because the sample
// creates a trigger and a Cloud Scheduler job with fixed names, which would
// conflict between concurrent test runs. The test verifies that the sample
// reaches the API and that the API rejects the request, rather than the sample
// failing during client setup.
func TestScheduleBuild(t *testing.T) {
	testutil.SystemTest(t)

	const invalidProjectID = "invalid-project-id-12345"
	err := scheduleBuild(io.Discard, invalidProjectID,
		"test-sa@"+invalidProjectID+".iam.gserviceaccount.com",
		"gs://test-bucket/workdir/", "gs://test-bucket/imagebuilder.yaml")
	if err == nil {
		t.Fatal("got nil error, want an API error for the invalid project")
	}
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) {
		t.Fatalf("got error %v (%T), want a *googleapi.Error", err, err)
	}
	if gerr.Code != 400 && gerr.Code != 403 && gerr.Code != 404 {
		t.Errorf("got HTTP status %d (%v), want 400, 403, or 404", gerr.Code, gerr)
	}
}
