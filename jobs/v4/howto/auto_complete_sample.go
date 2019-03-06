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
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// [START job_search_autocomplete_job_title]

// jobTitleAutoComplete suggests the job titles of the given companyName based
// on query.
func jobTitleAutocomplete(w io.Writer, projectId, companyName, query string) (*talentpb.CompleteQueryResponse, error) {
	ctx := context.Background()

	// Initialize a completionService client.
	c, err := talent.NewCompletionClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompletionClient: %v", err)
		return nil, err
	}

	// Construct a completeQuery request.
	req := &talentpb.CompleteQueryRequest{
		Name: "projects/" + projectId,
		Query: query,
		PageSize: 1,
		CompanyName: companyName,
	}

	resp, err := c.CompleteQuery(ctx, req)
	if err != nil {
		fmt.Printf("failed to auto complete with query %s in %s: %v", query, companyName, err)
		return nil, err
	}

	fmt.Fprintf(w, "Auto complete results:")
	for _, c := range resp.GetCompletionResults() {
		fmt.Fprintf(w, "\t%v", c.Suggestion)
	}

	return resp, nil
}

// [END job_search_autocomplete_job_title]
