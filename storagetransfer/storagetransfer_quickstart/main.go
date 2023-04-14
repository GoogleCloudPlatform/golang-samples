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

package storagetransfer_quickstart

// [START storagetransfer_quickstart]
import (
	"context"
	"fmt"
	"log"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

// quickstart creates and runs a transfer job between two GCS buckets.
func main() {
	// Your Google Cloud Project ID
	projectID := "my-project-id"

	// The name of the GCS bucket to transfer data from
	sourceGCSBucket := "my-source-bucket"

	// The name of the GCS bucket to transfer data to
	sinkGCSBucket := "my-sink-bucket"

	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		log.Fatalf("storagetransfer.NewClient: %v", err)
	}
	defer client.Close()

	req := &storagetransferpb.CreateTransferJobRequest{
		TransferJob: &storagetransferpb.TransferJob{
			ProjectId: projectID,
			TransferSpec: &storagetransferpb.TransferSpec{
				DataSource: &storagetransferpb.TransferSpec_GcsDataSource{
					GcsDataSource: &storagetransferpb.GcsData{BucketName: sourceGCSBucket}},
				DataSink: &storagetransferpb.TransferSpec_GcsDataSink{
					GcsDataSink: &storagetransferpb.GcsData{BucketName: sinkGCSBucket}},
			},
			Status: storagetransferpb.TransferJob_ENABLED,
		},
	}
	resp, err := client.CreateTransferJob(ctx, req)
	if err != nil {
		log.Fatalf("failed to create transfer job: %v", err)
	}
	if _, err = client.RunTransferJob(ctx, &storagetransferpb.RunTransferJobRequest{
		ProjectId: projectID,
		JobName:   resp.Name,
	}); err != nil {
		log.Fatalf("failed to run transfer job: %v", err)
	}

	fmt.Printf("Created and ran transfer job from %v to %v with name %v", sourceGCSBucket, sinkGCSBucket, resp.Name)
}

// [END storagetransfer_quickstart]
