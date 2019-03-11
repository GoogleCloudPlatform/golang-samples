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

// [START job_search_create_job]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// createJob create a job as given.
func createJob(w io.Writer, projectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode string) (*talentpb.Job, error) {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewJobClient: %v\n", err)
		return nil, err
	}

	jobToCreate := &talentpb.Job{
		CompanyName:   fmt.Sprintf("projects/%s/companies/%s", projectID, companyID),
		RequisitionId: requisitionID,
		Title:         title,
		ApplicationInfo: &talentpb.Job_ApplicationInfo{
			Uris: []string{URI},
		},
		Description:  description,
		Addresses:    []string{address1, address2},
		LanguageCode: languageCode,
	}

	// Construct a createJob request.
	req := &talentpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Job:    jobToCreate,
	}

	resp, err := c.CreateJob(ctx, req)
	if err != nil {
		fmt.Printf("Failed to create job: %v\n", err)
		return nil, err
	}

	fmt.Printf("Created job: %q\n", resp.GetName())

	return resp, nil
}

// [END job_search_create_job]
