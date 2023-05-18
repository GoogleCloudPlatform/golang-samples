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

// [START job_search_commute_search]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
	"github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// commuteSearch searches for jobs within commute filter.
func commuteSearch(w io.Writer, projectID, companyID string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewJobClient: %w", err)
	}
	defer c.Close()

	// Construct a jobQuery request with a commute filter.
	jobQuery := &talentpb.JobQuery{
		CommuteFilter: &talentpb.CommuteFilter{
			CommuteMethod:  talentpb.CommuteMethod_TRANSIT,
			TravelDuration: &duration.Duration{Seconds: 1800},
			StartCoordinates: &latlng.LatLng{
				Latitude:  37.422408,
				Longitude: -122.085609,
			},
		},
	}
	if companyID != "" {
		jobQuery.Companies = []string{fmt.Sprintf("projects/%s/companies/%s", projectID, companyID)}
	}

	// Construct a searchJobs request with a jobQuery.
	req := &talentpb.SearchJobsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talentpb.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrID",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionID",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
	}

	resp, err := c.SearchJobs(ctx, req)
	if err != nil {
		return fmt.Errorf("SearchJobs: %w", err)
	}

	for _, job := range resp.GetMatchingJobs() {
		fmt.Fprintf(w, "Matcing Job: %q\n", job.GetJob().GetName())
		fmt.Fprintf(w, "Job address: %v\n", job.GetCommuteInfo().GetJobLocation().GetPostalAddress().GetAddressLines())
	}

	return nil
}

// [END job_search_commute_search]
