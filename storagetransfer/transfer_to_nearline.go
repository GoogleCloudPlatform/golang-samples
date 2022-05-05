package storagetransfer

// [START storagetransfer_transfer_to_nearline]
import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/timeofday"
	"google.golang.org/protobuf/types/known/durationpb"
	"io"
	"time"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	storagetransferpb "google.golang.org/genproto/googleapis/storagetransfer/v1"
)

func transfer_to_nearline(w io.Writer, projectID string, jobDescription string, gcsSourceBucket string, gcsNearlineSinkBucket string, startTime time.Time) (*storagetransferpb.TransferJob, error) {
	// Your Google Cloud Project ID
	// projectID := "my-project-id"

	// A description of this job
	// jobDescription = "Transfers objects that haven't been modified in 30 days to a Nearline bucket"

	// The name of the GCS bucket to transfer objects from
	// gcsSourceBucket := "my-source-bucket"

	// The name of the Nearline GCS bucket to transfer objects to
	// gcsNearlineSinkBucket := "my-sink-bucket"

	// The time to start the transfer
	// startTime := time.Now()

	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storagetransfer.NewClient: %v", err)
	}
	defer client.Close()

	req := &storagetransferpb.CreateTransferJobRequest{
		TransferJob: &storagetransferpb.TransferJob{
			ProjectId:   projectID,
			Description: jobDescription,
			TransferSpec: &storagetransferpb.TransferSpec{
				DataSink: &storagetransferpb.TransferSpec_GcsDataSink{
					GcsDataSink: &storagetransferpb.GcsData{BucketName: gcsNearlineSinkBucket}},
				DataSource: &storagetransferpb.TransferSpec_GcsDataSource{
					GcsDataSource: &storagetransferpb.GcsData{BucketName: gcsSourceBucket},
				},
				ObjectConditions: &storagetransferpb.ObjectConditions{
					MinTimeElapsedSinceLastModification: &durationpb.Duration{Seconds: 2592000 /*30 days */},
				},
				TransferOptions: &storagetransferpb.TransferOptions{DeleteObjectsFromSourceAfterTransfer: true},
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
		return nil, fmt.Errorf("Failed to create transfer job: %v", err)
	}
	if _, err = client.RunTransferJob(ctx, &storagetransferpb.RunTransferJobRequest{
		ProjectId: projectID,
		JobName:   resp.Name,
	}); err != nil {
		return nil, fmt.Errorf("Failed to run transfer job: %v", err)
	}
	fmt.Fprintf(w, "Created and ran transfer job from %v to %v with name %v", gcsSourceBucket, gcsNearlineSinkBucket, resp.Name)
	return resp, nil
}

// [END storagetransfer_transfer_to_nearline]
