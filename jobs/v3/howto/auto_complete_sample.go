// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START auto_complete_job_title]

// jobTitleAutoComplete suggests the job titles of the given companyName based
// on query.
func jobTitleAutoComplete(w io.Writer, projectID, companyName, query string) (*talent.CompleteQueryResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
	}

	parent := "projects/" + projectID
	complete := service.Projects.Complete(parent).Query(query).LanguageCode("en-US").Type("JOB_TITLE").PageSize(10)
	if companyName != "" {
		complete.CompanyName(companyName)
	}
	resp, err := complete.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to auto complete with query %s in company %s: %v", query, companyName, err)
	}

	fmt.Fprintf(w, "Auto complete results:")
	for _, c := range resp.CompletionResults {
		fmt.Fprintf(w, "\t%v", c.Suggestion)
	}

	return resp, nil
}

// [END auto_complete_job_title]

// [START auto_complete_default]

// defaultAutoComplete suggests job titles or company display names of given companyName based on query.
func defaultAutoComplete(w io.Writer, projectID, companyName, query string) (*talent.CompleteQueryResponse, error) {
	ctx := context.Background()

	parent := "projects/" + projectID

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
	}

	complete := service.Projects.Complete(parent).Query(query).LanguageCode("en-US").Type("COMBINED").PageSize(10)
	if companyName != "" {
		complete.CompanyName(companyName)
	}
	resp, err := complete.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to auto complete with query %s in company %s: %v", query, companyName, err)
	}

	fmt.Fprintf(w, "Auto complete results:")
	for _, c := range resp.CompletionResults {
		fmt.Fprintf(w, "\t%v", c.Suggestion)
	}

	return resp, nil

}

// [END auto_complete_default]
