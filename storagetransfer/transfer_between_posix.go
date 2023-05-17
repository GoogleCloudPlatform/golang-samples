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

// [START storagetransfer_transfer_posix_to_posix]

import (
	"context"
	"fmt"
	"io"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

func transferBetweenPosix(w io.Writer, projectID string, sourceAgentPoolName string, sinkAgentPoolName string, rootDirectory string, destinationDirectory string, bucketName string) (*storagetransferpb.TransferJob, error) {
	// Your project id
	// projectId := "my-project-id"

	// The agent pool associated with the POSIX data source. If not provided, defaults to the default agent
	// sourceAgentPoolName := "projects/my-project/agentPools/transfer_service_default"

	// The agent pool associated with the POSIX data sink. If not provided, defaults to the default agent
	// sinkAgentPoolName := "projects/my-project/agentPools/transfer_service_default"

	// The root directory path on the source filesystem
	// rootDirectory := "/directory/to/transfer/source"

	// The root directory path on the sink filesystem
	// destinationDirectory := "/directory/to/transfer/sink"

	// The ID of the GCS bucket for intermediate storage
	// bucketName := "my-intermediate-bucket";

	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storagetransfer.NewClient: %w", err)
	}
	defer client.Close()

	req := &storagetransferpb.CreateTransferJobRequest{
		TransferJob: &storagetransferpb.TransferJob{
			ProjectId: projectID,
			TransferSpec: &storagetransferpb.TransferSpec{
				SourceAgentPoolName: sourceAgentPoolName,
				SinkAgentPoolName:   sinkAgentPoolName,
				DataSource: &storagetransferpb.TransferSpec_PosixDataSource{
					PosixDataSource: &storagetransferpb.PosixFilesystem{RootDirectory: rootDirectory},
				},
				DataSink: &storagetransferpb.TransferSpec_PosixDataSink{
					PosixDataSink: &storagetransferpb.PosixFilesystem{RootDirectory: destinationDirectory},
				},
				IntermediateDataLocation: &storagetransferpb.TransferSpec_GcsIntermediateDataLocation{
					GcsIntermediateDataLocation: &storagetransferpb.GcsData{BucketName: bucketName},
				},
			},
			Status: storagetransferpb.TransferJob_ENABLED,
		},
	}

	resp, err := client.CreateTransferJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer job: %w", err)
	}
	if _, err = client.RunTransferJob(ctx, &storagetransferpb.RunTransferJobRequest{
		ProjectId: projectID,
		JobName:   resp.Name,
	}); err != nil {
		return nil, fmt.Errorf("failed to run transfer job: %w", err)
	}
	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", rootDirectory, destinationDirectory, resp.Name)
	return resp, nil
}

// [END storagetransfer_transfer_posix_to_posix]
