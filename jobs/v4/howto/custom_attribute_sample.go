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

// [START job_search_create_job_custom_attributes]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"github.com/gofrs/uuid"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
	money "google.golang.org/genproto/googleapis/type/money"
)

// createJobWithCustomAttributes creates a job with custom attributes.
func createJobWithCustomAttributes(w io.Writer, projectID, companyID, jobTitle string) (*talentpb.Job, error) {
	ctx := context.Background()

	// Initialize a job service client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("talent.NewJobClient: %v", err)
	}

	// requisitionID shoud be the unique ID in your system
	requisitionID := fmt.Sprintf("job-with-custom-attribute-%s", uuid.Must(uuid.NewV4()).String())
	jobToCreate := &talentpb.Job{
		Company:       fmt.Sprintf("projects/%s/companies/%s", projectID, companyID),
		RequisitionId: requisitionID,
		Title:         jobTitle,
		ApplicationInfo: &talentpb.Job_ApplicationInfo{
			Uris: []string{"https://googlesample.com/career"},
		},
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
		LanguageCode:    "en-US",
		PromotionValue:  2,
		EmploymentTypes: []talentpb.EmploymentType{talentpb.EmploymentType_FULL_TIME},
		Addresses:       []string{"Mountain View, CA"},
		CustomAttributes: map[string]*talentpb.CustomAttribute{
			"someFieldString": {
				Filterable:   true,
				StringValues: []string{"someStrVal"},
			},
			"someFieldLong": {
				Filterable: true,
				LongValues: []int64{900},
			},
		},
		CompensationInfo: &talentpb.CompensationInfo{
			Entries: []*talentpb.CompensationInfo_CompensationEntry{
				{
					Type: talentpb.CompensationInfo_BASE,
					Unit: talentpb.CompensationInfo_HOURLY,
					CompensationAmount: &talentpb.CompensationInfo_CompensationEntry_Amount{
						Amount: &money.Money{
							CurrencyCode: "USD",
							Units:        1,
						},
					},
				},
			},
		},
	}

	// Construct a createJob request.
	req := &talentpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Job:    jobToCreate,
	}

	resp, err := c.CreateJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateJob: %v", err)
	}

	fmt.Fprintf(w, "Created job with custom attributres: %q\n", resp.GetName())
	fmt.Fprintf(w, "Custom long field has value: %v\n", resp.GetCustomAttributes()["someFieldLong"].GetLongValues())

	return resp, nil
}

// [END job_search_create_job_custom_attributes]
