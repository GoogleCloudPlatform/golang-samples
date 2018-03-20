/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"io"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	"google.golang.org/api/iterator"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func listJobs(w io.Writer, client *dlp.Client, project, filter, jobType string) {
	rcr := &dlppb.ListDlpJobsRequest{
		Parent: "projects/" + project,
		Filter: filter,
		Type:   dlppb.DlpJobType(dlppb.DlpJobType_value[jobType]),
	}
	it := client.ListDlpJobs(context.Background(), rcr)
	for {
		j, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error getting jobs: %v", err)
		}
		fmt.Fprintf(w, "Job %v status: %v", j.GetName(), j.GetState())
	}
}

func deleteJob(w io.Writer, client *dlp.Client, jobName string) {
	rcr := &dlppb.DeleteDlpJobRequest{
		Name: jobName,
	}
	err := client.DeleteDlpJob(context.Background(), rcr)
	if err != nil {
		log.Fatalf("error deleting job: %v", err)
	}
}
