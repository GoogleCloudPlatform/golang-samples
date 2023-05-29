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
package jobs

// [START dlp_create_job]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// createJob creates an inspection job
func createJob(w io.Writer, projectID, gcsPath string, infoTypeNames []string) error {
	// projectID := "my-project-id"
	// gcsPath := "gs://" + "your-bucket-name" + "path/to/file.txt";
	// infoTypeNames := []string{"EMAIL_ADDRESS", "PERSON_NAME", "LOCATION", "PHONE_NUMBER"}

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

	// Specify the GCS file to be inspected.
	storageConfig := &dlppb.StorageConfig{
		Type: &dlppb.StorageConfig_CloudStorageOptions{
			CloudStorageOptions: &dlppb.CloudStorageOptions{
				FileSet: &dlppb.CloudStorageOptions_FileSet{
					Url: gcsPath,
				},
			},
		},

		// Set autoPopulateTimespan to true to scan only new content.
		TimespanConfig: &dlppb.StorageConfig_TimespanConfig{
			EnableAutoPopulationOfTimespanConfig: true,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types.
	var infoTypes []*dlppb.InfoType
	for _, c := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: c})
	}

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes:    infoTypes,
		IncludeQuote: true,

		// The minimum likelihood required before returning a match:
		// See: https://cloud.google.com/dlp/docs/likelihood
		MinLikelihood: dlppb.Likelihood_UNLIKELY,

		// The maximum number of findings to report (0 = server maximum)
		Limits: &dlppb.InspectConfig_FindingLimits{
			MaxFindingsPerItem: 100,
		},
	}

	// Create and send the request.
	req := dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
				InspectConfig: inspectConfig,
				StorageConfig: storageConfig,
			},
		},
	}

	// Send the request.
	response, err := client.CreateDlpJob(ctx, &req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "Created a Dlp Job %v and Status is: %v", response.Name, response.State)
	return nil
}

// [END dlp_create_job]
