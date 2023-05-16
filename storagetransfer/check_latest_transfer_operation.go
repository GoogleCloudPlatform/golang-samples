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

// [START storagetransfer_get_latest_transfer_operation]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
)

func checkLatestTransferOperation(w io.Writer, projectID string, jobName string) (*storagetransferpb.TransferOperation, error) {
	// Your Google Cloud Project ID
	// projectID := "my-project-id"

	// The name of the job whose latest operation to check
	// jobName := "transferJobs/1234567890"
	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storagetransfer.NewClient: %w", err)
	}
	defer client.Close()

	job, err := client.GetTransferJob(ctx, &storagetransferpb.GetTransferJobRequest{
		JobName:   jobName,
		ProjectId: projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer job: %w", err)
	}

	latestOpName := job.LatestOperationName
	if latestOpName != "" {
		lro, err := client.LROClient.GetOperation(ctx, &longrunningpb.GetOperationRequest{
			Name: latestOpName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get transfer operation: %w", err)
		}
		latestOp := &storagetransferpb.TransferOperation{}
		lro.Metadata.UnmarshalTo(latestOp)

		fmt.Fprintf(w, "the latest transfer operation for job %q is: \n%v", jobName, latestOp)
		return latestOp, nil
	} else {
		fmt.Fprintf(w, "Transfer job %q hasn't run yet, try again later", jobName)
		return nil, nil
	}
}

// [END storagetransfer_get_latest_transfer_operation]
