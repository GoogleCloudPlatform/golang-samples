package storagetransfer

// [START storagetransfer_quickstart]
import (
	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"context"
	"fmt"
	storagetransferpb "google.golang.org/genproto/googleapis/storagetransfer/v1"
	"io"
)

// quickstart creates and runs a transfer job between two GCS buckets.
func quickstart(w io.Writer, projectID string, sourceGCSBucket string, sinkGCSBucket string) (*storagetransferpb.TransferJob, error) {
	// Your Google Cloud Project ID
	// projectID := "my-project-id"

	// The name of the GCS bucket to transfer data from
	// sourceGCSBucket := "my-source-bucket"

	// The name of the GCS bucket to transfer data to
	// sinkGCSBucket := "my-sink-bucket"
	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storagetransfer.NewClient: %v", err)
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
		return nil, fmt.Errorf("Failed to create transfer job: %v", err)
	}
	if _, err = client.RunTransferJob(ctx, &storagetransferpb.RunTransferJobRequest{
		ProjectId: projectID,
		JobName: resp.Name,
	}); err != nil {
		return nil, fmt.Errorf("Failed to run transfer job: %v", err)
	}

	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", sourceGCSBucket, sinkGCSBucket, resp.Name)
	return resp, nil
}
// [END storagetransfer_quickstart]