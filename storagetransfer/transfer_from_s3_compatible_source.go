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

// [START storagetransfer_transfer_from_s3_compatible_source]

import (
	"context"
	"fmt"
	"io"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

func transferFromS3CompatibleSource(w io.Writer, projectID string, sourceAgentPoolName string, sourceBucketName string, sourcePath string, gcsSinkBucket string, gcsPath string) (*storagetransferpb.TransferJob, error) {
	// Your project id.
	// projectId := "my-project-id"

	// The agent pool associated with the S3 compatible data source. If not provided, defaults to the default agent.
	// sourceAgentPoolName := "projects/my-project/agentPools/transfer_service_default"

	// The S3 compatible bucket name to transfer data from.
	//sourceBucketName = "my-bucket-name"

	// The S3 compatible path (object prefix) to transfer data from.
	//sourcePath = "path/to/data"

	// The ID of the GCS bucket to transfer data to.
	//gcsSinkBucket = "my-sink-bucket"

	// The GCS path (object prefix) to transfer data to.
	//gcsPath = "path/to/data"

	// The S3 region of the source bucket.
	region := "us-east-1"

	// The S3 compatible endpoint.
	endpoint := "us-east-1.example.com"

	// The S3 compatible network protocol.
	protocol := storagetransferpb.S3CompatibleMetadata_NETWORK_PROTOCOL_HTTPS

	// The S3 compatible request model.
	requestModel := storagetransferpb.S3CompatibleMetadata_REQUEST_MODEL_VIRTUAL_HOSTED_STYLE

	// The S3 Compatible auth method.
	authMethod := storagetransferpb.S3CompatibleMetadata_AUTH_METHOD_AWS_SIGNATURE_V4

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
				DataSource: &storagetransferpb.TransferSpec_AwsS3CompatibleDataSource{
					AwsS3CompatibleDataSource: &storagetransferpb.AwsS3CompatibleData{
						BucketName: sourceBucketName,
						Path:       sourcePath,
						Endpoint:   endpoint,
						Region:     region,
						DataProvider: &storagetransferpb.AwsS3CompatibleData_S3Metadata{
							S3Metadata: &storagetransferpb.S3CompatibleMetadata{
								AuthMethod:   authMethod,
								RequestModel: requestModel,
								Protocol:     protocol,
							},
						},
					}},
				DataSink: &storagetransferpb.TransferSpec_GcsDataSink{
					GcsDataSink: &storagetransferpb.GcsData{
						BucketName: gcsSinkBucket,
						Path:       gcsPath,
					},
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
	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", sourceBucketName, gcsSinkBucket, resp.Name)
	return resp, nil
}

// [END storagetransfer_transfer_from_s3_compatible_source]
