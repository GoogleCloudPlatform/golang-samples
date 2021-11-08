package storagetransfer

// [START storagetransfer_quickstart]
import (
	"context"
	"fmt"
	"io"
	"log"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	storagetransferpb "google.golang.org/genproto/googleapis/storagetransfer/v1"
)

// quickstart creates and runs a transfer job between two GCS buckets.
func quickstart(w io.Writer, projectID string, sourceGCSBucket string, sinkGCSBucket string) error {
	// Your Google Cloud Project ID
	// projectID := "my-project-id"

	// The name of the GCS bucket to transfer data from
	// sourceGCSBucket := "my-source-bucket"

	// The name of the GCS bucket to transfer data to
	// sinkGCSBucket := "my-sink-bucket"
	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storagetransfer.NewClient: %v", err)
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
		log.Fatalf("Failed to create transfer job: %v", err)
	}
	_, err = client.RunTransferJob(ctx, &storagetransferpb.RunTransferJobRequest{
		ProjectId: projectID,
		JobName: resp.Name,
	})

	if err != nil {
		log.Fatalf("Failed to run transfer job: %v", err)
	}
	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", sourceGCSBucket, sinkGCSBucket, resp.Name)
	return nil
}
// [END storagetransfer_quickstart]