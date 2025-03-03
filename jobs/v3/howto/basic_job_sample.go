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

package howto

import (
	"context"
	"fmt"
	"io"
	"log"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// createJob create a job as given.
func createJob(w io.Writer, projectID string, jobToCreate *talent.Job) (*talent.Job, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	parent := "projects/" + projectID
	req := &talent.CreateJobRequest{
		Job: jobToCreate,
	}
	job, err := service.Projects.Jobs.Create(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create job %q, Err: %w", jobToCreate.RequisitionId, err)
	}
	return job, err
}

// getJob gets a job by name.
func getJob(w io.Writer, jobName string) (*talent.Job, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	job, err := service.Projects.Jobs.Get(jobName).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to get job %s: %w", jobName, err)
	}

	fmt.Fprintf(w, "Job: %q", job.Name)

	return job, err
}

// updateJob update a job with all fields except name.
func updateJob(w io.Writer, jobName string, jobToUpdate *talent.Job) (*talent.Job, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	req := &talent.UpdateJobRequest{
		Job: jobToUpdate,
	}
	job, err := service.Projects.Jobs.Patch(jobName, req).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to update job %s: %w", jobName, err)
	}

	return job, err
}

// updateJobWithMask updates a job by name with specific fields.
// mask is a comma separated list top-level fields of talent.Job.
func updateJobWithMask(w io.Writer, jobName string, mask string, jobToUpdate *talent.Job) (*talent.Job, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	req := &talent.UpdateJobRequest{
		Job:        jobToUpdate,
		UpdateMask: mask,
	}
	job, err := service.Projects.Jobs.Patch(jobName, req).Do()
	if err != nil {
		log.Fatalf("Failed to update job %s with field mask %s, Err: %v", jobName, mask, err)
	}

	return job, err
}

// deleteJob deletes an existing job by name.
func deleteJob(w io.Writer, jobName string) error {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return fmt.Errorf("talent.New: %w", err)
	}

	if _, err := service.Projects.Jobs.Delete(jobName).Do(); err != nil {
		return fmt.Errorf("Delete(%s): %w", jobName, err)
	}

	return err
}

// listJobs lists jobs with a filter, for example
// `companyName="projects/my-project/companies/123"`.
func listJobs(w io.Writer, projectID, filter string) (*talent.ListJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	parent := "projects/" + projectID
	resp, err := service.Projects.Jobs.List(parent).Filter(filter).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to list jobs with filter: %q: %w", filter, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.Jobs {
		fmt.Fprintf(w, "\t%q\n", j.Name)
	}

	return resp, err
}
