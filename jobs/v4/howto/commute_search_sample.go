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

	"github.com/golang/protobuf/ptypes/duration"
  talent "cloud.google.com/go/talent/apiv4beta1"
	"google.golang.org/api/iterator"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// [START job_search_commute_search]

// commuteSearch searches for jobs within commute filter.
func commuteSearch(w io.Writer, projectID, companyName string) error {
	ctx := context.Background()

	// Create a job service client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewJobClient: %v", err)
	}

	jobQuery := &talentpb.JobQuery{
		CommuteFilter: &talentpb.CommuteFilter{
			CommuteMethod:  2,
			TravelDuration: &duration.Duration{Seconds:1800},
			StartCoordinates: &latlng.LatLng{
				Latitude:  37.422408,
				Longitude: -122.085609,
			},
		},
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	req := &talentpb.SearchJobsRequest{
		Parent: "projects/" + projectID,
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talentpb.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
	}
	it := c.SearchJobs(ctx, req)
	fmt.Fprintln(w, "Jobs:")
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("it.Next: %v", err)
		}
		fmt.Printf("%v\n", resp.Job.Name)
		fmt.Fprintf(w, "\t%q\n", resp.Job.Name)
	}
}

// [END job_search_commute_search]
