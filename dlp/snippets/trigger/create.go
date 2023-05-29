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

package trigger

// [START dlp_create_trigger]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"github.com/golang/protobuf/ptypes/duration"
)

// createTrigger creates a trigger with the given configuration.
func createTrigger(w io.Writer, projectID string, triggerID, displayName, description, bucketName string, infoTypeNames []string) error {
	// projectID := "my-project-id"
	// triggerID := "my-trigger"
	// displayName := "My Trigger"
	// description := "My trigger description"
	// bucketName := "my-bucket"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}

	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Convert the info type strings to a list of InfoTypes.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	// Create a configured request.
	req := &dlppb.CreateJobTriggerRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/global", projectID),
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
									Seconds: 10 * 60 * 60 * 24, // 10 days in seconds.
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
						InfoTypes:     infoTypes,
						MinLikelihood: dlppb.Likelihood_POSSIBLE,
						Limits: &dlppb.InspectConfig_FindingLimits{
							MaxFindingsPerRequest: 10,
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
							EnableAutoPopulationOfTimespanConfig: true,
						},
					},
				},
			},
		},
	}

	// Send the request.
	resp, err := client.CreateJobTrigger(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateJobTrigger: %w", err)
	}
	fmt.Fprintf(w, "Successfully created trigger: %v", resp.GetName())
	return nil
}

// [END dlp_create_trigger]
