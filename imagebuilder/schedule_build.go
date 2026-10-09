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

// [START compute_imagebuilder_schedule_build]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/cloudscheduler/v1"
)

// scheduleBuild creates a manual Cloud Build trigger that runs an Image
// Builder pipeline, and a Cloud Scheduler job that runs the trigger every day.
//
// Before you run this sample, upload your imagebuilder.yaml configuration file
// to Cloud Storage, create a generic Artifact Registry repository, and create a
// service account to run the build.
//
// The job runs the trigger as the default scheduling service account,
// cloud-build-trigger-scheduler@PROJECT_ID.iam.gserviceaccount.com. Grant that
// service account the Cloud Build Editor role, and the Service Account User
// role on the service account that runs the build. For details, see
// https://cloud.google.com/build/docs/schedule-builds.
func scheduleBuild(w io.Writer, projectID, serviceAccount, gcsWorkdir, configPath string) error {
	// projectID := "my-project-id"
	// serviceAccount := "image-builder@my-project-id.iam.gserviceaccount.com"
	// gcsWorkdir := "gs://my-bucket/workdir/"
	// configPath := "gs://my-bucket/imagebuilder.yaml"
	// Region of the trigger, job, and repository.
	region := "us-central1"
	// Generic Artifact Registry repository.
	repository := "vm-images"
	// Package that receives the image.
	packageName := "custom-os-images"
	// Name for the new build trigger and Cloud Scheduler job.
	triggerName := "scheduled-image-builder"
	jobName := "scheduled-image-builder-job"
	// Every day at 09:00 UTC, in unix-cron format.
	schedule := "0 9 * * *"

	ctx := context.Background()
	cbService, err := cloudbuild.NewService(ctx)
	if err != nil {
		return fmt.Errorf("cloudbuild.NewService: %w", err)
	}
	csService, err := cloudscheduler.NewService(ctx)
	if err != nil {
		return fmt.Errorf("cloudscheduler.NewService: %w", err)
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	sa := fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, serviceAccount)

	// Create a manual trigger that runs the Image Builder pipeline.
	trigger := &cloudbuild.BuildTrigger{
		Name:        triggerName,
		Description: "Manual trigger for scheduled VM image builds",
		EventType:   "MANUAL",
		Build: &cloudbuild.Build{
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
		},
	}
	createdTrigger, err := cbService.Projects.Locations.Triggers.Create(parent, trigger).
		Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Triggers.Create: %w", err)
	}
	fmt.Fprintf(w, "Created trigger %s (ID: %s)\n", createdTrigger.Name, createdTrigger.Id)

	// Create a Cloud Scheduler job that runs the trigger on the schedule.
	schedulerSA := fmt.Sprintf(
		"cloud-build-trigger-scheduler@%s.iam.gserviceaccount.com", projectID,
	)
	job := &cloudscheduler.Job{
		Name:        fmt.Sprintf("%s/jobs/%s", parent, jobName),
		Description: "Runs the Image Builder trigger on a schedule",
		Schedule:    schedule,
		HttpTarget: &cloudscheduler.HttpTarget{
			Uri: fmt.Sprintf(
				"https://cloudbuild.googleapis.com/v1/projects/%s/locations/%s/triggers/%s:run",
				projectID, region, createdTrigger.Id,
			),
			HttpMethod: "POST",
			OauthToken: &cloudscheduler.OAuthToken{
				ServiceAccountEmail: schedulerSA,
			},
		},
	}
	createdJob, err := csService.Projects.Locations.Jobs.Create(parent, job).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Jobs.Create: %w", err)
	}
	fmt.Fprintf(w, "Created job %s with schedule %q\n", createdJob.Name, createdJob.Schedule)
	return nil
}

// [END compute_imagebuilder_schedule_build]
