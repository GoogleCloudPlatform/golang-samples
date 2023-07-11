// Copyright 2023 Google LLC
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

import (
	"context"
	"fmt"
	"io"
	"os"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

// [START storagetransfer_create_event_driven_aws_transfer]

func createEventDrivenAWSTransfer(w io.Writer, projectID string, s3SourceBucket string, gcsSinkBucket string, sqsQueueARN string) (*storagetransferpb.TransferJob, error) {
	// Your Google Cloud Project ID.
	// projectID := "my-project-id"

	// The name of the source AWS S3 bucket.
	// s3SourceBucket := "my-source-bucket"

	// The name of the GCS bucket to transfer objects to.
	// gcsSinkBucket := "my-sink-bucket"

	// The Amazon Resource Name (ARN) of the AWS SNS queue to subscribe the event driven transfer to.
	// sqsQueueARN := "arn:aws:sqs:us-east-1:1234567891011:s3-notification-queue"

	// The AWS access key credential, should be accessed via environment variable for security
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")

	// The AWS secret key credential, should be accessed via environment variable for security
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

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
				DataSource: &storagetransferpb.TransferSpec_AwsS3DataSource{
					AwsS3DataSource: &storagetransferpb.AwsS3Data{
						BucketName: s3SourceBucket,
						AwsAccessKey: &storagetransferpb.AwsAccessKey{
							AccessKeyId:     awsAccessKeyID,
							SecretAccessKey: awsSecretKey,
						}},
				},
				DataSink: &storagetransferpb.TransferSpec_GcsDataSink{
					GcsDataSink: &storagetransferpb.GcsData{BucketName: gcsSinkBucket}},
			},
			EventStream: &storagetransferpb.EventStream{Name: sqsQueueARN},
			Status:      storagetransferpb.TransferJob_ENABLED,
		},
	}
	resp, err := client.CreateTransferJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer job: %w", err)
	}

	fmt.Fprintf(w, "Created an event driven transfer job from %v to %v subscribed to %v with name %v", s3SourceBucket, gcsSinkBucket, sqsQueueARN, resp.Name)
	return resp, nil
}

// [END storagetransfer_create_event_driven_aws_transfer]
