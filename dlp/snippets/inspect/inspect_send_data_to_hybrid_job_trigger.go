// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package inspect

// [START dlp_inspect_send_data_to_hybrid_job_trigger]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectDataToHybridJobTrigger uses the Data Loss Prevention API to inspect sensitive
// information using Hybrid jobs trigger that scans payloads of data sent from
// virtually any source and stores findings in Google Cloud.
func inspectDataToHybridJobTrigger(w io.Writer, projectID, textToDeIdentify, jobTriggerName string) error {
	// projectId := "your-project-id"
	// jobTriggerName := "your-job-trigger-name"
	// textToDeIdentify := "My email is test@example.org"

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the content to be inspected.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: textToDeIdentify,
		},
	}

	// Contains metadata to associate with the content.
	container := &dlppb.Container{
		Type:         "logging_sys",
		FullPath:     "10.0.0.2:logs1:app1",
		RelativePath: "app1",
		RootPath:     "10.0.0.2:logs1",
		Version:      "1.2",
	}

	// Set the required label.
	labels := map[string]string{
		"env":                           "prod",
		"appointment-bookings-comments": "",
	}

	hybridFindingDetails := &dlppb.HybridFindingDetails{
		ContainerDetails: container,
		Labels:           labels,
	}

	hybridContentItem := &dlppb.HybridContentItem{
		Item:           contentItem,
		FindingDetails: hybridFindingDetails,
	}

	// Activate the job trigger.
	activateJobreq := &dlppb.ActivateJobTriggerRequest{
		Name: jobTriggerName,
	}
	client.ActivateJobTrigger(ctx, activateJobreq)

	// Build the hybrid inspect request.
	req := &dlppb.HybridInspectJobTriggerRequest{
		Name:       jobTriggerName,
		HybridItem: hybridContentItem,
	}

	// Send the hybrid inspect request.
	resp, err := client.HybridInspectJobTrigger(ctx, req)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "successfully inspected data using hybrid job trigger: %v", resp)
	return nil
}

// [END dlp_inspect_send_data_to_hybrid_job_trigger]
