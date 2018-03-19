// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	"google.golang.org/api/iterator"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_list_jobs]
// listJobs lists jobs matching the given optional filter and optional jobType.
func listJobs(w io.Writer, client *dlp.Client, project, filter, jobType string) {
	// Create a configured request.
	req := &dlppb.ListDlpJobsRequest{
		Parent: "projects/" + project,
		Filter: filter,
		Type:   dlppb.DlpJobType(dlppb.DlpJobType_value[jobType]),
	}
	// Send the request and iterate over the results.
	it := client.ListDlpJobs(context.Background(), req)
	for {
		j, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error getting jobs: %v", err)
		}
		fmt.Fprintf(w, "Job %v status: %v\n", j.GetName(), j.GetState())
	}
}

// [END dlp_list_jobs]

// [START dlp_delete_job]
// deleteJob deletes the job with the given name.
func deleteJob(w io.Writer, client *dlp.Client, jobName string) {
	req := &dlppb.DeleteDlpJobRequest{
		Name: jobName,
	}
	err := client.DeleteDlpJob(context.Background(), req)
	if err != nil {
		log.Fatalf("error deleting job: %v", err)
	}
	fmt.Fprintf(w, "Successfully deleted job")
}

// [END dlp_delete_job]
