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

// [START storagetransfer_download_to_posix]

import (
	"context"
	"fmt"
	"io"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

func downloadToPosix(w io.Writer, projectID string, sinkAgentPoolName string, gcsSourceBucket string, gcsSourcePath string, rootDirectory string) (*storagetransferpb.TransferJob, error) {
	// Your project id
	// projectId := "my-project-id"

	// The agent pool associated with the POSIX data sink. Defaults to the default agent if not specified
	// Make sure that the pub/sub resources are set up for this agent pool, or you'll get an error (See the
	// "agent pools" tab in the Data Transfer Console)
	// sinkAgentPoolName := "projects/my-project/agentPools/transfer_service_default"

	// Your GCS source bucket name
	// gcsSourceBucket := "my-gcs-source-bucket"

	// A directory prefix on the Google Cloud Storage bucket to download from
	// gcsSourcePath := "foo/bar/"

	// The root directory path on the source filesystem
	// rootDirectory := "/path/to/transfer/source"

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
				SinkAgentPoolName: sinkAgentPoolName,
				DataSource: &storagetransferpb.TransferSpec_GcsDataSource{
					GcsDataSource: &storagetransferpb.GcsData{
						BucketName: gcsSourceBucket,
						Path:       gcsSourcePath,
					},
				},
				DataSink: &storagetransferpb.TransferSpec_PosixDataSink{
					PosixDataSink: &storagetransferpb.PosixFilesystem{RootDirectory: rootDirectory},
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
	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", gcsSourceBucket, rootDirectory, resp.Name)
	return resp, nil
}

// [END storagetransfer_download_to_posix]
