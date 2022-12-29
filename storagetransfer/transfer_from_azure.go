// Copyright 2022 Google LLC
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

package storagetransfer

// [START storagetransfer_transfer_from_azure]
import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/timeofday"
)

func transferFromAzure(w io.Writer, projectID string, azureStorageAccountName string, azureSourceContainer string, gcsSinkBucket string) (*storagetransferpb.TransferJob, error) {
	// Your Google Cloud Project ID
	// projectID := "my-project-id"

	// The name of your Azure Storage account
	// azureStorageAccountName := "my-azure-storage-acc"

	// The name of the Azure container to transfer objects from
	// azureSourceContainer := "my-source-container"

	// The name of the GCS bucket to transfer objects to
	// gcsSinkBucket := "my-sink-bucket"

	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storagetransfer.NewClient: %v", err)
	}
	defer client.Close()

	// A description of this job
	jobDescription := "Transfers objects from an Azure bucket to a GCS bucket"

	// The time to start the transfer
	startTime := time.Now()

	// The Azure SAS token, should be accessed via environment variable for security
	azureSasToken := os.Getenv("AZURE_SAS_TOKEN")

	req := &storagetransferpb.CreateTransferJobRequest{
		TransferJob: &storagetransferpb.TransferJob{
			ProjectId:   projectID,
			Description: jobDescription,
			TransferSpec: &storagetransferpb.TransferSpec{
				DataSource: &storagetransferpb.TransferSpec_AzureBlobStorageDataSource{
					AzureBlobStorageDataSource: &storagetransferpb.AzureBlobStorageData{
						StorageAccount: azureStorageAccountName,
						AzureCredentials: &storagetransferpb.AzureCredentials{
							SasToken: azureSasToken,
						},
						Container: azureSourceContainer,
					},
				},
				DataSink: &storagetransferpb.TransferSpec_GcsDataSink{
					GcsDataSink: &storagetransferpb.GcsData{BucketName: gcsSinkBucket}},
			},
			Schedule: &storagetransferpb.Schedule{
				ScheduleStartDate: &date.Date{
					Year:  int32(startTime.Year()),
					Month: int32(startTime.Month()),
					Day:   int32(startTime.Day()),
				},
				ScheduleEndDate: &date.Date{
					Year:  int32(startTime.Year()),
					Month: int32(startTime.Month()),
					Day:   int32(startTime.Day()),
				},
				StartTimeOfDay: &timeofday.TimeOfDay{
					Hours:   int32(startTime.Hour()),
					Minutes: int32(startTime.Minute()),
					Seconds: int32(startTime.Second()),
				},
			},
			Status: storagetransferpb.TransferJob_ENABLED,
		},
	}
	resp, err := client.CreateTransferJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer job: %v", err)
	}
	if _, err = client.RunTransferJob(ctx, &storagetransferpb.RunTransferJobRequest{
		ProjectId: projectID,
		JobName:   resp.Name,
	}); err != nil {
		return nil, fmt.Errorf("failed to run transfer job: %v", err)
	}
	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", azureSourceContainer, gcsSinkBucket, resp.Name)
	return resp, nil
}

// [END storagetransfer_transfer_from_azure]
