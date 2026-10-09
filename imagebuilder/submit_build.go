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

// Package imagebuilder contains samples that run Image Builder pipelines on
// Cloud Build.
package imagebuilder

// [START compute_imagebuilder_submit_build]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"google.golang.org/api/cloudbuild/v1"
)

// submitBuild submits a Cloud Build build that uses Image Builder to
// customize, validate, and publish a VM image to Compute Engine and Artifact
// Registry.
//
// Before you run this sample, upload your imagebuilder.yaml configuration file
// to Cloud Storage, create a generic Artifact Registry repository, and create a
// service account to run the build.
func submitBuild(w io.Writer, projectID, serviceAccount, gcsWorkdir, configPath string) error {
	// projectID := "my-project-id"
	// serviceAccount := "image-builder@my-project-id.iam.gserviceaccount.com"
	// gcsWorkdir := "gs://my-bucket/workdir/"
	// configPath := "gs://my-bucket/imagebuilder.yaml"
	// Region of the build and the repository.
	region := "us-central1"
	// Generic Artifact Registry repository.
	repository := "vm-images"
	// Package that receives the image.
	packageName := "custom-os-images"

	ctx := context.Background()
	cbService, err := cloudbuild.NewService(ctx)
	if err != nil {
		return fmt.Errorf("cloudbuild.NewService: %w", err)
	}

	sa := fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, serviceAccount)
	req := &cloudbuild.Build{
		ServiceAccount: sa,
		// AutomapSubstitutions exposes these substitutions to every step as
		// environment variables. See
		// https://cloud.google.com/build/docs/configuring-builds/substitute-variable-values#mapping_substitutions_to_environment_variables
		Substitutions: map[string]string{
			"_GCS_WORKDIR":               gcsWorkdir,
			"_SERVICE_ACCOUNT":           sa,
			"_IMAGE_OUTPUT_PATH":         "image-builder/binaryOut",
			"_IMAGE_BUILDER_CONFIG_PATH": configPath,
			"_ARTIFACT_REGISTRY_RESOURCE_URI": fmt.Sprintf(
				"projects/%s/locations/%s/repositories/%s/packages/%s/versions/v${BUILD_ID}",
				projectID, region, repository, packageName,
			),
		},
		Steps: []*cloudbuild.BuildStep{
			{
				// Customize the image as defined in the configuration file.
				Id:     "imagebuilder-customize",
				Name:   "us-docker.pkg.dev/image-builder-official/release/builder:stable",
				Script: "#!/usr/bin/env bash\n/build",
				Results: []*cloudbuild.StepResult{
					{
						Name:            "base_image",
						AttestationType: "https://cloudbuild.googleapis.com/attestations/build_content_restrictions",
					},
					{Name: "image_builder_telemetry_metrics"},
				},
			},
			{
				// Validate the customized image.
				Id:      "imagebuilder-validate",
				Name:    "us-docker.pkg.dev/image-builder-official/release/validator:stable",
				Script:  "#!/usr/bin/env bash\n/validate",
				Results: []*cloudbuild.StepResult{{Name: "image_builder_telemetry_metrics"}},
			},
			{
				// Publish the image to Compute Engine and Artifact Registry.
				Id:      "imagebuilder-publish",
				Name:    "us-docker.pkg.dev/image-builder-official/release/builder:stable",
				Script:  "#!/usr/bin/env bash\n/publish",
				Results: []*cloudbuild.StepResult{{Name: "image_builder_telemetry_metrics"}},
			},
		},
		Options: &cloudbuild.BuildOptions{
			AutomapSubstitutions:  true,
			DynamicSubstitutions:  true,
			RequestedVerifyOption: "VERIFIED",
			SubstitutionOption:    "ALLOW_LOOSE",
			// Builds that run as a user-specified service account must set a
			// logging option.
			Logging: "CLOUD_LOGGING_ONLY",
		},
		Artifacts: &cloudbuild.Artifacts{
			GenericArtifacts: []*cloudbuild.GenericArtifact{
				{
					Folder:       "${_IMAGE_OUTPUT_PATH}",
					RegistryPath: "${_ARTIFACT_REGISTRY_RESOURCE_URI}",
				},
			},
		},
		Timeout: "3600s",
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	op, err := cbService.Projects.Locations.Builds.Create(parent, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Builds.Create: %w", err)
	}

	// The operation metadata contains the build that was created.
	var metadata cloudbuild.BuildOperationMetadata
	if err := json.Unmarshal(op.Metadata, &metadata); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	fmt.Fprintf(w, "Build submitted. ID: %s\n", metadata.Build.Id)
	fmt.Fprintf(w, "Logs: %s\n", metadata.Build.LogUrl)
	return nil
}

// [END compute_imagebuilder_submit_build]
