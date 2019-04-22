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

// [START job_search_custom_ranking_search]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"google.golang.org/api/iterator"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// customRankingSearch searches for jobs based on custom ranking.
func customRankingSearch(w io.Writer, projectID, companyID string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("taleng.NewJobClient: %v", err)
	}

	// Construct a searchJobs request.
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
		JobQuery: &talentpb.JobQuery{
			Companies: []string{fmt.Sprintf("projects/%s/companies/%s", projectID, companyID)},
		},
		// More info on customRankingInfo.
		// https://godoc.org/google.golang.org/genproto/googleapis/cloud/talent/v4beta1#SearchJobsRequest_CustomRankingInfo
		CustomRankingInfo: &talentpb.SearchJobsRequest_CustomRankingInfo{
			ImportanceLevel:   talentpb.SearchJobsRequest_CustomRankingInfo_EXTREME,
			RankingExpression: "(someFieldLong + 25) * 0.25",
		},
		OrderBy: "custom_ranking desc",
	}

	it := c.SearchJobs(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("it.Next: %v", err)
		}
		fmt.Fprintf(w, "Job: %q\n", resp.Job.GetName())
	}
}

// [END job_search_custom_ranking_search]
