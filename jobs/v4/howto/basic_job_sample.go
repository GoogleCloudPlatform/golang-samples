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

	talent "cloud.google.com/go/talent/apiv4beta1"
	"google.golang.org/api/iterator"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// [START job_search_create_job]

// createJob create a job as given.
func createJob(w io.Writer, projectId string, jobToCreate *talentpb.Job) (*talentpb.Job, error) {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewJobClient: %v", err)
		return nil, err
	}

  // Construct a createJob request.
	req := &talentpb.CreateJobRequest{
		Parent: "projects/" + projectId,
		Job: jobToCreate,
	}

	resp, err := c.CreateJob(ctx, req)
	if err != nil {
		fmt.Printf("Failed to create job: %v", err)
		return nil, err
	}

	fmt.Printf("Created job: %q\n", resp.GetName())

	return resp, nil
}

// [END job_search_create_job]

// [START job_search_get_job]

// getJob gets an existing job by name.
func getJob(w io.Writer, jobName string) (*talentpb.Job, error) {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewJobClient: %v", err)
		return nil, err
	}

	// Construct a getJob request.
	req := &talentpb.GetJobRequest{
		// The resource name of the job to retrieve.
    // The format is "projects/{project_id}/jobs/{job_id}".
		Name: jobName,
	}

	resp, err := c.GetJob(ctx, req)
	if err != nil {
		fmt.Printf("Failed to get job %s: %v", jobName, err)
		return nil, err
	}

	fmt.Fprintf(w, "Job: %q\n", resp.GetName())
	fmt.Printf("Job title: %v\n", resp.GetTitle())

	return resp, err
}

// [END job_search_get_job]

// [START job_search_delete_job]

// deleteJob deletes an existing job by name.
func deleteJob(w io.Writer, jobName string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewJobClient: %v", err)
		return err
	}

	// Construct a deleteJob request.
	req := &talentpb.DeleteJobRequest{
		// The resource name of the job to be deleted.
		// The format is "projects/{project_id}/jobs/{job_id}".
		Name: jobName,
	}

	if err := c.DeleteJob(ctx, req); err != nil {
		fmt.Printf("Delete(%s): %v", jobName, err)
		return err
	}

	fmt.Printf("Deleted job: %q\n", jobName)

	return err
}

// [END job_search_delete_job]

// [START job_search_list_jobs]

// listJobs lists jobs with a filter, for example
// `companyName="projects/my-project/companies/123"`.
func listJobs(w io.Writer, projectId, filter string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewJobClient: %v", err)
		return err
	}

	// Construct a listJobs request.
	req := &talentpb.ListJobsRequest{
		Parent: "projects/" + projectId,
		Filter: filter,
	}

	it := c.ListJobs(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			fmt.Printf("it.Next: %v", err)
			return err
		}
		fmt.Fprintf(w, "Listing job: %v\n", resp.GetName())
		fmt.Printf("Job title: %v\n", resp.GetTitle())
	}
}

// [END job_search_list_jobs]
