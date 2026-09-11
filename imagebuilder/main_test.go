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

package main

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/googleapi"
)

func TestCreateBuildRequest(t *testing.T) {
	resetFlags(t)

	// Set flags
	*projectID = "test-project"
	*serviceAccount = "test-sa"
	*gcsWorkdir = "gs://test-bucket/workdir"
	*imageBuilderYAML = "gs://test-bucket/config.yaml"

	req := createBuildRequest()
	if req == nil {
		t.Fatal("createBuildRequest() returned nil")
	}

	want := &cloudbuild.Build{
		Steps: []*cloudbuild.BuildStep{
			{
				Name:   "us-docker.pkg.dev/image-builder-prod/release/builder",
				Script: "#!/usr/bin/env bash\n/build",
				Id:     "imagebuilder-customize",
				Results: []*cloudbuild.StepResult{
					{
						Name:            "base_image",
						AttestationType: "https://cloudbuild.googleapis.com/attestations/build_content_restrictions",
					},
					{
						Name: "image_builder_telemetry_metrics",
					},
				},
			},
			{
				Name:   "us-docker.pkg.dev/image-builder-prod/release/validator",
				Script: "#!/usr/bin/env bash\n/validate",
				Id:     "imagebuilder-validate",
				Results: []*cloudbuild.StepResult{
					{
						Name: "image_builder_telemetry_metrics",
					},
				},
			},
			{
				Name:   "us-docker.pkg.dev/image-builder-prod/release/builder",
				Script: "#!/usr/bin/env bash\n/publish",
				Id:     "imagebuilder-publish",
				Results: []*cloudbuild.StepResult{
					{
						Name: "image_builder_telemetry_metrics",
					},
				},
			},
		},
		ServiceAccount: "projects/test-project/serviceAccounts/test-sa",
		Substitutions: map[string]string{
			"_GCS_WORKDIR":                    "gs://test-bucket/workdir",
			"_SERVICE_ACCOUNT":                "projects/test-project/serviceAccounts/test-sa",
			"_IMAGE_OUTPUT_PATH":              "gce-image-builder/binaryOut/image.tar.gz",
			"_IMAGE_BUILDER_CONFIG_PATH":      "gs://test-bucket/config.yaml",
			"_ARTIFACT_REGISTRY_RESOURCE_URI": "projects/test-project/locations/us-central1/repositories/vm-images/packages/custom-os-images/versions/v${BUILD_ID}",
		},
		Options: &cloudbuild.BuildOptions{
			AutomapSubstitutions:  true,
			DynamicSubstitutions:  true,
			RequestedVerifyOption: "VERIFIED",
			SubstitutionOption:    "ALLOW_LOOSE",
			Logging:               "CLOUD_LOGGING_ONLY",
		},
		Artifacts: &cloudbuild.Artifacts{
			GenericArtifacts: []*cloudbuild.GenericArtifact{
				{
					Folder:       "/workspace/gce-image-builder/binaryOut",
					RegistryPath: "${_ARTIFACT_REGISTRY_RESOURCE_URI}",
				},
			},
		},
		Timeout: "3600s",
	}

	if diff := cmp.Diff(want, req); diff != "" {
		t.Errorf("createBuildRequest() mismatch (-want +got):\n%s", diff)
	}
}

func TestTriggerBuildInvalidProject(t *testing.T) {
	testutil.SystemTest(t)
	resetFlags(t)

	// Use an invalid project ID to verify the API call executes and fails as expected
	*projectID = "invalid-project-id-12345"
	*serviceAccount = "test-sa@invalid-project-id-12345.iam.gserviceaccount.com"
	*gcsWorkdir = "gs://test-bucket/workdir"
	*imageBuilderYAML = "gs://test-bucket/config.yaml"

	err := triggerBuild(io.Discard)
	if err == nil {
		t.Fatal("triggerBuild expected error for invalid project, got nil")
	}

	// The error should come from the API (not client initialization)
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) {
		t.Errorf("expected *googleapi.Error, got: %T (%v)", err, err)
	} else if gerr.Code != 400 && gerr.Code != 403 && gerr.Code != 404 {
		t.Errorf("expected HTTP 400, 403, or 404, got status %d: %v", gerr.Code, gerr)
	}
}

func TestMainFlagValidation(t *testing.T) {
	m := testutil.BuildMain(t)
	t.Cleanup(m.Cleanup)
	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	// Run with missing flags
	_, stdErr, err := m.Run(nil, 15*time.Second)
	if err == nil {
		t.Error("expected exit status 1 due to missing flags, got exit code 0")
	}
	if !strings.Contains(string(stdErr), "Error: --project_id, --gcs_workdir, and --service_account are required flags") {
		t.Errorf("expected error message about required flags, got stderr: %q", string(stdErr))
	}
}

func resetFlags(t *testing.T) {
	origProject := *projectID
	origSA := *serviceAccount
	origWorkdir := *gcsWorkdir
	origConfig := *imageBuilderYAML
	t.Cleanup(func() {
		*projectID = origProject
		*serviceAccount = origSA
		*gcsWorkdir = origWorkdir
		*imageBuilderYAML = origConfig
	})
}
