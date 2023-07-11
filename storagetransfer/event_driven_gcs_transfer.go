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

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

// [START storagetransfer_create_event_driven_gcs_transfer]

func createEventDrivenGCSTransfer(w io.Writer, projectID string, gcsSourceBucket string, gcsSinkBucket string, pubSubId string) (*storagetransferpb.TransferJob, error) {
	// Your Google Cloud Project ID.
	// projectID := "my-project-id"

	// The name of the source GCS bucket.
	// gcsSourceBucket := "my-source-bucket"

	// The name of the GCS bucket to transfer objects to.
	// gcsSinkBucket := "my-sink-bucket"

	// The Pub/Sub topic to subscribe the event driven transfer to.
	// pubSubID := "projects/PROJECT_NAME/subscriptions/SUBSCRIPTION_ID"

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
				DataSource: &storagetransferpb.TransferSpec_GcsDataSource{
					GcsDataSource: &storagetransferpb.GcsData{BucketName: gcsSourceBucket}},
				DataSink: &storagetransferpb.TransferSpec_GcsDataSink{
					GcsDataSink: &storagetransferpb.GcsData{BucketName: gcsSinkBucket}},
			},
			EventStream: &storagetransferpb.EventStream{Name: pubSubId},
			Status:      storagetransferpb.TransferJob_ENABLED,
		},
	}
	resp, err := client.CreateTransferJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer job: %w", err)
	}

	fmt.Fprintf(w, "Created an event driven transfer job from %v to %v subscribed to %v with name %v", gcsSourceBucket, gcsSinkBucket, pubSubId, resp.Name)
	return resp, nil
}

// [END storagetransfer_create_event_driven_gcs_transfer]
