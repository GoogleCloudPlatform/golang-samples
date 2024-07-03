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
	"time"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START job_custom_attribute_job]

// constructJobWithCustomAttributes constructs a job with custom attributes.
func constructJobWithCustomAttributes(companyName string, jobTitle string) *talent.Job {
	// requisitionID shoud be the unique ID in your system
	requisitionID := fmt.Sprintf("job-with-custom-attribute-%d", time.Now().UnixNano())

	job := &talent.Job{
		RequisitionId: requisitionID,
		Title:         jobTitle,
		CompanyName:   companyName,
		ApplicationInfo: &talent.ApplicationInfo{
			Uris: []string{"https://googlesample.com/career"},
		},
		Description: "Design, devolop, test, deploy, maintain and improve software.",
		CustomAttributes: map[string]talent.CustomAttribute{
			"someFieldString": {
				Filterable:   true,
				StringValues: []string{"someStrVal"},
			},
			"someFieldLong": {
				Filterable: true,
				LongValues: []int64{900},
			},
		},
	}
	return job
}

// [END job_custom_attribute_job]

// [START job_custom_attribute_filter_string_value]

// filterOnStringValueCustomAttribute searches for jobs on a string value custom
// atrribute.
func filterOnStringValueCustomAttribute(w io.Writer, projectID string) (*talent.SearchJobsResponse, error) {
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
	req := &talent.SearchJobsRequest{
		JobQuery: &talent.JobQuery{
			CustomAttributeFilter: "NOT EMPTY(someFieldString)",
		},
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talent.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		JobView: "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with string value custom attribute: %w", err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_custom_attribute_filter_string_value]

// [START job_custom_attribute_filter_long_value]

// filterOnLongValueCustomAttribute searches for jobs on a long value custom
// atrribute.
func filterOnLongValueCustomAttribute(w io.Writer, projectID string) (*talent.SearchJobsResponse, error) {
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
	req := &talent.SearchJobsRequest{
		JobQuery: &talent.JobQuery{
			CustomAttributeFilter: "someFieldLong < 1000",
		},
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talent.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		JobView: "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with long value custom attribute: %w", err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_custom_attribute_filter_long_value]

// [START job_custom_attribute_filter_multi_attributes]

// filterOnLongValueCustomAttribute searches for jobs on multiple custom
// atrributes.
func filterOnMultiCustomAttributes(w io.Writer, projectID string) (*talent.SearchJobsResponse, error) {
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
	req := &talent.SearchJobsRequest{
		JobQuery: &talent.JobQuery{
			CustomAttributeFilter: "(someFieldString = \"someStrVal\") AND (someFieldLong < 1000)",
		},
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talent.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		JobView: "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with multiple custom attributes: %w", err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_custom_attribute_filter_multi_attributes]
