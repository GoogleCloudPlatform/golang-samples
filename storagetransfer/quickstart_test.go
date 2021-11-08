package storagetransfer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	storagetransferpb "google.golang.org/genproto/googleapis/storagetransfer/v1"
)

var str *storage.Client
var sts *storagetransfer.Client

func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	str = c
	defer str.Close()

	s, err := storagetransfer.NewClient(ctx)
	if err != nil {
		log.Fatalf("storagetransfer.NewClient: %v", err)
	}
	sts = s
	defer sts.Close()

	sinkBucketName, err := testutil.CreateTestBucket(ctx, t, str, tc.ProjectID, "sts-go-sink")
	if err != nil {
		t.Fatalf("Bucket creation failed: %v", err)
	}
	defer testutil.DeleteBucketIfExists(ctx, str, sinkBucketName)

	sourceBucketName, err := testutil.CreateTestBucket(ctx, t, str, tc.ProjectID, "sts-go-source")
	if err != nil {
		t.Fatalf("Bucket creation failed: %v", err)
	}
	defer testutil.DeleteBucketIfExists(ctx, str, sourceBucketName)

	grantStsPermissions(sinkBucketName, tc.ProjectID)
	grantStsPermissions(sourceBucketName, tc.ProjectID)

	buf := new(bytes.Buffer)
	if err := quickstart(buf, tc.ProjectID, sourceBucketName, sinkBucketName); err != nil {
		t.Errorf("quickstart: %#v", err)
	}

	got := buf.String()
	if want := "transferJobs/"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	re := regexp.MustCompile("transferJobs/\\d+")
	name := re.FindString(got)

	tj := &storagetransferpb.TransferJob{
		Name: name,
		Status: storagetransferpb.TransferJob_DELETED,
	}
	sts.UpdateTransferJob(ctx, &storagetransferpb.UpdateTransferJobRequest{
		JobName: name,
		ProjectId: tc.ProjectID,
		TransferJob: tj,
	})
}

func grantStsPermissions(bucketName string, projectID string) {
	ctx := context.Background()

	req := &storagetransferpb.GetGoogleServiceAccountRequest{
		ProjectId: projectID,
	}

	resp, err := sts.GetGoogleServiceAccount(ctx, req)
	if err != nil {
		fmt.Print("Error getting service account")
	}
	email := resp.AccountEmail

	identity := "serviceAccount:" + email

	bucket := str.Bucket(bucketName)
	policy, err := bucket.IAM().Policy(ctx)
	if err != nil {
		log.Fatalf("Bucket(%q).IAM().Policy: %v", bucketName, err)
	}

	var objectViewer iam.RoleName = "roles/storage.objectViewer"
	var bucketReader iam.RoleName = "roles/storage.legacyBucketReader"
	var bucketWriter iam.RoleName = "roles/storage.legacyBucketWriter"

	policy.Add(identity, objectViewer)
	policy.Add(identity, bucketReader)
	policy.Add(identity, bucketWriter)

	if err := bucket.IAM().SetPolicy(ctx, policy); err != nil {
		log.Fatalf("Bucket(%q).IAM().SetPolicy: %v", bucketName, err)
	}
}
