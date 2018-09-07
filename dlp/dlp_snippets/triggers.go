// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/duration"

	dlp "cloud.google.com/go/dlp/apiv2"
	"google.golang.org/api/iterator"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_create_trigger]

// createTrigger creates a trigger with the given configuration.
func createTrigger(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, triggerID, displayName, description, bucketName string, autoPopulateTimespan bool, scanPeriodDays int64, infoTypes []string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}

	// Create a configured request.
	req := &dlppb.CreateJobTriggerRequest{
		Parent:    "projects/" + project,
		TriggerId: triggerID,
		JobTrigger: &dlppb.JobTrigger{
			DisplayName: displayName,
			Description: description,
			Status:      dlppb.JobTrigger_HEALTHY,
			// Triggers control when the job will start.
			Triggers: []*dlppb.JobTrigger_Trigger{
				{
					Trigger: &dlppb.JobTrigger_Trigger_Schedule{
						Schedule: &dlppb.Schedule{
							Option: &dlppb.Schedule_RecurrencePeriodDuration{
								RecurrencePeriodDuration: &duration.Duration{
									Seconds: scanPeriodDays * 60 * 60 * 24, // Days to seconds.
								},
							},
						},
					},
				},
			},
			// Job configures the job to run when the trigger runs.
			Job: &dlppb.JobTrigger_InspectJob{
				InspectJob: &dlppb.InspectJobConfig{
					InspectConfig: &dlppb.InspectConfig{
						InfoTypes:     i,
						MinLikelihood: minLikelihood,
						Limits: &dlppb.InspectConfig_FindingLimits{
							MaxFindingsPerRequest: maxFindings,
						},
					},
					StorageConfig: &dlppb.StorageConfig{
						Type: &dlppb.StorageConfig_CloudStorageOptions{
							CloudStorageOptions: &dlppb.CloudStorageOptions{
								FileSet: &dlppb.CloudStorageOptions_FileSet{
									Url: "gs://" + bucketName + "/*",
								},
							},
						},
						// Time-based configuration for each storage object. See more at
						// https://cloud.google.com/dlp/docs/reference/rest/v2/InspectJobConfig#TimespanConfig
						TimespanConfig: &dlppb.StorageConfig_TimespanConfig{
							// Auto-populate start and end times in order to scan new objects only.
							EnableAutoPopulationOfTimespanConfig: autoPopulateTimespan,
						},
					},
				},
			},
		},
	}
	// Send the request.
	resp, err := client.CreateJobTrigger(context.Background(), req)
	if err != nil {
		log.Fatalf("error creating job trigger: %v", err)
	}
	fmt.Fprintf(w, "Successfully created trigger: %v", resp.GetName())
}

// [END dlp_create_trigger]

// [START dlp_list_triggers]

// listTriggers lists the triggers for the given project.
func listTriggers(w io.Writer, client *dlp.Client, project string) {
	// Create a configured request.
	req := &dlppb.ListJobTriggersRequest{
		Parent: "projects/" + project,
	}
	// Send the request and iterate over the results.
	it := client.ListJobTriggers(context.Background(), req)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error getting jobs: %v", err)
		}
		c := t.GetCreateTime()
		u := t.GetUpdateTime()
		fmt.Fprintf(w, "Trigger %v\n", t.GetName())
		fmt.Fprintf(w, "  Created: %v\n", time.Unix(c.GetSeconds(), int64(c.GetNanos())).Format(time.RFC1123))
		fmt.Fprintf(w, "  Updated: %v\n", time.Unix(u.GetSeconds(), int64(u.GetNanos())).Format(time.RFC1123))
		fmt.Fprintf(w, "  Display Name: %q\n", t.GetDisplayName())
		fmt.Fprintf(w, "  Description: %q\n", t.GetDescription())
		fmt.Fprintf(w, "  Status: %v\n", t.GetStatus())
		fmt.Fprintf(w, "  Error Count: %v\n", len(t.GetErrors()))
	}
}

// [END dlp_list_triggers]

// [START dlp_delete_trigger]

// deleteTrigger deletes the given trigger.
func deleteTrigger(w io.Writer, client *dlp.Client, triggerID string) {
	req := &dlppb.DeleteJobTriggerRequest{
		Name: triggerID,
	}
	err := client.DeleteJobTrigger(context.Background(), req)
	if err != nil {
		log.Fatalf("error deleting job: %v", err)
	}
	fmt.Fprintf(w, "Successfully deleted trigger %v", triggerID)
}

// [END dlp_delete_trigger]
