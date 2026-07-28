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

// Package main demonstrates how to submit a VM Image Builder build using public Go API
// Client SDKs that runs on a defined schedule.
//
// Build the binary using "go build" and invoke it like:
//
// For a scheduled job:
//
//	./schedule \
//		--project_id <YOUR_PROJECT_ID> \
//		--gcs_workdir gs://<GCS_BUCKET>/<FOLDER>/ \
//		--service_account <YOUR_SERVICE_ACCOUNT_EMAIL> \
//		--config_path gs://<SOURCE_GCS_BUCKET>/imagebuilder.yaml \
//		--cron_schedule '0 9 * * *'
//
// For submitting a one-off job:
//
//	./schedule \
//		--project_id <YOUR_PROJECT_ID> \
//		--gcs_workdir gs://<GCS_BUCKET>/<FOLDER>/ \
//		--service_account <YOUR_SERVICE_ACCOUNT_EMAIL> \
//		--config_path gs://<SOURCE_GCS_BUCKET>/imagebuilder.yaml
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/cloudscheduler/v1"
	"google.golang.org/api/googleapi"
)

var (
	projectID        = flag.String("project_id", "", "The target GCP Project ID")
	region           = flag.String("region", "us-central1", "GCP region of resources")
	gcsWorkdir       = flag.String("gcs_workdir", "", "The GCS workspace directory URI, e.g. gs://my-bucket/workdir/")
	serviceAccount   = flag.String("service_account", "", "Service account email to run the OS image customization pipeline")
	imageBuilderYAML = flag.String("config_path", "imagebuilder.yaml", "GCS Path to the imagebuilder.yaml configuration file")
	arRepositoryID   = flag.String("ar_repo_id", "vm-images", "Name of the target generic Artifact Registry repository")
	arPackageName    = flag.String("ar_package_name", "custom-os-images", "Package identifier for VM images")
	cronSchedule     = flag.String("cron_schedule", "", "Optional cron schedule to run this build periodically (e.g. '0 9 * * *')")
)

func main() {
	flag.Parse()

	if *projectID == "" || *serviceAccount == "" || *gcsWorkdir == "" {
		fmt.Fprintln(os.Stderr, "Error: --project_id, --gcs_workdir, and --service_account are required flags")
		flag.Usage()
		os.Exit(1)
	}

	if *cronSchedule != "" {
		fmt.Printf("Scheduling build with cron schedule: %q\n", *cronSchedule)
		if err := scheduleBuild(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to schedule Cloud Build: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Trigger Cloud Build using target config
		fmt.Println("Submitting build request to Cloud Build API...")
		if err := triggerBuild(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to trigger Cloud Build: %v\n", err)
			os.Exit(1)
		}
	}
}

// createBuildRequest defines and structures the Cloud Build steps (customize, validate,
// and publish) along with required substitutions and options to build the VM image.
func createBuildRequest() *cloudbuild.Build {
	sa := fmt.Sprintf("projects/%s/serviceAccounts/%s", *projectID, *serviceAccount)
	subs := map[string]string{
		"_GCS_WORKDIR":               *gcsWorkdir,
		"_SERVICE_ACCOUNT":           sa,
		"_IMAGE_OUTPUT_PATH":         "gce-image-builder/binaryOut/image.tar.gz",
		"_IMAGE_BUILDER_CONFIG_PATH": *imageBuilderYAML,
		"_ARTIFACT_REGISTRY_RESOURCE_URI": fmt.Sprintf(
			"projects/%s/locations/%s/repositories/%s/packages/%s/versions/v${BUILD_ID}",
			*projectID, *region, *arRepositoryID, *arPackageName,
		),
	}

	return &cloudbuild.Build{
		Steps: []*cloudbuild.BuildStep{
			{
				Name: "us-docker.pkg.dev/image-builder-prod/release/builder",
				// Map substitutions to environment variables directly to make them accessible within container without
				// explicitly passing them in the step. See
				// https://docs.cloud.google.com/build/docs/configuring-builds/substitute-variable-values#mapping_substitutions_to_environment_variables
				Script: `#!/usr/bin/env bash
/build`,
				Id: "imagebuilder-customize",
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
				Name: "us-docker.pkg.dev/image-builder-prod/release/validator",
				Script: `#!/usr/bin/env bash
/validate`,
				Id: "imagebuilder-validate",
				Results: []*cloudbuild.StepResult{
					{
						Name: "image_builder_telemetry_metrics",
					},
				},
			},
			{
				Name: "us-docker.pkg.dev/image-builder-prod/release/builder",
				Script: `#!/usr/bin/env bash
/publish`,
				Id: "imagebuilder-publish",
				Results: []*cloudbuild.StepResult{
					{
						Name: "image_builder_telemetry_metrics",
					},
				},
			},
		},
		ServiceAccount: sa,
		Substitutions:  subs,
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
}

// triggerBuild constructs a Cloud Build API request and submits it.
// It maps the execution steps (builder, validator, publisher) as well as the
// necessary environment variables and substitutions dynamically.
func triggerBuild(w io.Writer) error {
	ctx := context.Background()
	cbService, err := cloudbuild.NewService(ctx)
	if err != nil {
		return fmt.Errorf("cloudbuild.NewService: %w", err)
	}

	req := createBuildRequest()

	parent := fmt.Sprintf("projects/%s/locations/%s", *projectID, *region)
	op, err := cbService.Projects.Locations.Builds.Create(parent, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Builds.Create: %w", err)
	}

	// Retrieve build trigger details
	var buildMetadata cloudbuild.BuildOperationMetadata
	if err := json.Unmarshal(op.Metadata, &buildMetadata); err != nil {
		return fmt.Errorf("failed to retrieve build metadata: %w", err)
	}
	buildID := buildMetadata.Build.Id
	logURL := buildMetadata.Build.LogUrl
	fmt.Fprintf(w, "Build successfully triggered! ID: %q\n", buildID)
	fmt.Fprintf(w, "Monitor real-time logs at: %q\n", logURL)

	return nil
}

// scheduleBuild creates a manual Cloud Build trigger and a Cloud Scheduler job to invoke it.
func scheduleBuild(w io.Writer) error {
	ctx := context.Background()
	cbService, err := cloudbuild.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Cloud Build service: %w", err)
	}

	csService, err := cloudscheduler.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Cloud Scheduler service: %w", err)
	}

	triggerName := "scheduled-image-builder"
	triggerParent := fmt.Sprintf("projects/%s/locations/%s", *projectID, *region)

	// 1. Create the manual trigger
	trigger := &cloudbuild.BuildTrigger{
		Name:        triggerName,
		Description: "Manual trigger for scheduled VM image builds",
		EventType:   "MANUAL",
		Build:       createBuildRequest(),
	}

	fmt.Fprintln(w, "Creating manual Cloud Build trigger...")
	createdTrigger, err := cbService.Projects.Locations.Triggers.Create(triggerParent, trigger).Context(ctx).Do()
	if err != nil {
		var gerr *googleapi.Error
		if errors.As(err, &gerr) && gerr.Code == 409 {
			fmt.Fprintln(w, "Note: manual Cloud Build trigger already exists.")
			createdTrigger = trigger
		} else {
			return fmt.Errorf("failed to create Cloud Build trigger: %w", err)
		}
	} else {
		fmt.Fprintf(w, "Trigger created! ID: %q\n", createdTrigger.Id)
	}

	// 2. Create the Cloud Scheduler job
	jobName := "scheduled-image-builder-job"
	jobParent := fmt.Sprintf("projects/%s/locations/%s", *projectID, *region)

	// Default scheduling SA.
	// https://docs.cloud.google.com/build/docs/schedule-builds
	sa := fmt.Sprintf("cloud-build-trigger-scheduler@%s.iam.gserviceaccount.com", *projectID)

	job := &cloudscheduler.Job{
		Name:        fmt.Sprintf("%s/jobs/%s", jobParent, jobName),
		Description: "Invokes the manual Cloud Build trigger on schedule",
		Schedule:    *cronSchedule,
		HttpTarget: &cloudscheduler.HttpTarget{
			Uri: fmt.Sprintf(
				"https://cloudbuild.googleapis.com/v1/projects/%s/locations/%s/triggers/%s:run",
				*projectID, *region, createdTrigger.Name,
			),
			HttpMethod: "POST",
			OauthToken: &cloudscheduler.OAuthToken{
				ServiceAccountEmail: sa,
			},
		},
	}

	fmt.Fprintf(w, "Creating Cloud Scheduler job %q with SA %q...\n", jobName, sa)
	createdJob, err := csService.Projects.Locations.Jobs.Create(jobParent, job).Context(ctx).Do()
	if err != nil {
		var gerr *googleapi.Error
		if errors.As(err, &gerr) && gerr.Code == 409 {
			fmt.Fprintln(w, "Note: Cloud Scheduler job already exists.")
		} else {
			return fmt.Errorf("failed to create Cloud Scheduler job: %w", err)
		}
	} else {
		fmt.Fprintln(w, "Successfully scheduled build! Job details:")
		fmt.Fprintf(w, "Name: %q\n", createdJob.Name)
		fmt.Fprintf(w, "Schedule: %q\n", createdJob.Schedule)
		fmt.Fprintf(w, "State: %q\n", createdJob.State)
	}

	return nil
}
