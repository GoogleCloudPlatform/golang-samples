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

package trigger

// [START dlp_update_trigger]

import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateTrigger updates an existing job trigger in Google Cloud Data Loss Prevention (DLP).
// It modifies the configuration of the specified job trigger with the provided updated settings.
func updateTrigger(w io.Writer, jobTriggerName string) error {
	// jobTriggerName := "your-job-trigger-name" (projects/<projectID>/locations/global/jobTriggers/my-trigger)

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

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoType := &dlppb.InfoType{
		Name: "PERSON_NAME",
	}

	// Specify the inspectConfig that represents the configuration settings for inspecting sensitive data in
	// DLP API. It includes detection types, custom info types, inspection methods, and actions
	// to be taken on detection.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			infoType,
		},
		MinLikelihood: dlppb.Likelihood_LIKELY,
	}

	// Configure the inspection job we want the service to perform.
	inspectJobConfig := &dlppb.InspectJobConfig{
		InspectConfig: inspectConfig,
	}

	// Specify the jobTrigger that represents a DLP job trigger configuration.
	// It defines the conditions, actions, and schedule for executing inspections
	// on sensitive data in the specified data storage.
	jobTrigger := &dlppb.JobTrigger{
		Job: &dlppb.JobTrigger_InspectJob{
			InspectJob: inspectJobConfig,
		},
	}

	// fieldMask represents a set of fields to be included in an update operation.
	// It is used to specify which fields of a resource should be updated.
	updateMask := &fieldmaskpb.FieldMask{
		Paths: []string{"inspect_job.inspect_config.info_types", "inspect_job.inspect_config.min_likelihood"},
	}

	// Combine configurations into a request for the service.
	req := &dlppb.UpdateJobTriggerRequest{
		Name:       jobTriggerName,
		JobTrigger: jobTrigger,
		UpdateMask: updateMask,
	}

	// Send the scan request and process the response
	resp, err := client.UpdateJobTrigger(ctx, req)
	if err != nil {
		return err
	}

	// Print the result.
	fmt.Fprintf(w, "Successfully Updated trigger: %v", resp)
	return nil

}

// [END dlp_update_trigger]
