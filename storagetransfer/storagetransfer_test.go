// Copyright 2021 Google LLC
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
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sc *storage.Client
var sts *storagetransfer.Client
var sess *session.Session
var s3Bucket string
var azureContainer string
var gcsSourceBucket string
var gcsSinkBucket string
var stsServiceAccountEmail string

func TestMain(m *testing.M) {
	// Initialize global vars
	tc, _ := testutil.ContextMain(m)

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	sc = c
	defer sc.Close()

	gcsSourceBucket = testutil.UniqueBucketName("gcssourcebucket")
	source := sc.Bucket(gcsSourceBucket)
	err = source.Create(ctx, tc.ProjectID, nil)
	if err != nil {
		log.Fatalf("couldn't create GCS Source bucket: %v", err)
	}

	gcsSinkBucket = testutil.UniqueBucketName("gcssinkbucket")
	sink := sc.Bucket(gcsSinkBucket)
	err = sink.Create(ctx, tc.ProjectID, nil)
	if err != nil {
		log.Fatalf("couldn't create GCS Sink bucket: %v", err)
	}

	sts, err = storagetransfer.NewClient(ctx)
	if err != nil {
		log.Fatalf("storagetransfer.NewClient: %v", err)
	}
	defer sts.Close()

	req := &storagetransferpb.GetGoogleServiceAccountRequest{
		ProjectId: tc.ProjectID,
	}

	resp, err := sts.GetGoogleServiceAccount(ctx, req)
	if err != nil {
		log.Fatalf("error getting service account: %v", err)
	}
	stsServiceAccountEmail = resp.AccountEmail

	grantSTSPermissions(gcsSourceBucket, sc)
	grantSTSPermissions(gcsSinkBucket, sc)

	s3Bucket = testutil.UniqueBucketName("stss3bucket")
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	s3c := s3.New(sess)
	_, err = s3c.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s3Bucket),
	})
	if err != nil {
		log.Fatalf("couldn't create S3 bucket: %v", err)
	}

	connectionString := os.Getenv("AZURE_CONNECTION_STRING") +
		";" + "AccountName=" + os.Getenv("AZURE_STORAGE_ACCOUNT")
	azClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatal("Couldn't create Azure client: " + err.Error())
	}
	azureContainer = testutil.UniqueBucketName("azurebucket")

	azClient.CreateContainer(ctx, azureContainer, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Run tests
	exit := m.Run()

	err = sink.Delete(ctx)
	if err != nil {
		log.Printf("couldn't delete GCS Sink bucket: %v", err)
	}

	err = source.Delete(ctx)
	if err != nil {
		log.Printf("couldn't delete GCS Source bucket: %v", err)
	}
	s3manager.NewDeleteListIterator(s3c, &s3.ListObjectsInput{
		Bucket: aws.String(s3Bucket),
	})
	_, err = s3c.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(s3Bucket),
	})
	if err != nil {
		log.Printf("couldn't delete S3 bucket: %v", err)
	}

	_, err = azClient.DeleteContainer(ctx, azureContainer, nil)
	if err != nil {
		log.Printf("couldn't delete Azure bucket: %v", err)
	}

	os.Exit(exit)
}

func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := quickstart(buf, tc.ProjectID, gcsSourceBucket, gcsSinkBucket)
		defer cleanupSTSJob(resp, tc.ProjectID)

		if err != nil {
			r.Errorf("quickstart: %#v", err)
		}

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("quickstart: got %q, want %q", got, want)
		}
	})
}

func TestTransferFromAws(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferFromAws(buf, tc.ProjectID, s3Bucket, gcsSinkBucket)
		defer cleanupSTSJob(resp, tc.ProjectID)

		if err != nil {
			r.Errorf("transfer_from_aws: %#v", err)
		}

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_from_aws: got %q, want %q", got, want)
		}
	})
}

func TestTransferToNearline(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferToNearline(buf, tc.ProjectID, gcsSourceBucket, gcsSinkBucket)
		defer cleanupSTSJob(resp, tc.ProjectID)

		if err != nil {
			r.Errorf("transfer_from_aws: %#v", err)
		}

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_to_nearline: got %q, want %q", got, want)
		}
	})
}

func TestGetLatestTransferOperation(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		job, err := transferToNearline(buf, tc.ProjectID, gcsSourceBucket, gcsSinkBucket)
		defer cleanupSTSJob(job, tc.ProjectID)

		op, err := checkLatestTransferOperation(buf, tc.ProjectID, job.Name)

		if err != nil {
			r.Errorf("check_latest_transfer_operation: %#v", err)
		}
		if !strings.Contains(op.Name, "transferOperations/") {
			r.Errorf("check_latest_transfer_operation: Operation returned didn't have a valid operation name: %q", op.Name)
		}
		got := buf.String()
		if want := op.Name; !strings.Contains(got, want) {
			r.Errorf("check_latest_transfer_operation: got %q, want %q", got, want)
		}
	})
}

func TestDownloadToPosix(t *testing.T) {
	tc := testutil.SystemTest(t)

	rootDirectory, err := ioutil.TempDir("", "download-to-posix-test")
	if err != nil {
		t.Fatalf("download_to_posix: %#v", err)
	}
	defer os.RemoveAll(rootDirectory)

	sinkAgentPoolName := "" //use default agent pool
	gcsSourcePath := rootDirectory + "/"

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := downloadToPosix(buf, tc.ProjectID, sinkAgentPoolName, gcsSinkBucket, gcsSourcePath, rootDirectory)
		defer cleanupSTSJob(resp, tc.ProjectID)

		if err != nil {
			r.Errorf("download_to_posix: %#v", err)
		}

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("download_to_posix: got %q, want %q", got, want)
		}
	})
}

func TestTransferFromPosix(t *testing.T) {
	tc := testutil.SystemTest(t)

	rootDirectory, err := ioutil.TempDir("", "transfer-from-posix-test")
	if err != nil {
		t.Fatalf("transfer_from_posix: %#v", err)
	}
	defer os.RemoveAll(rootDirectory)

	sourceAgentPoolName := "" //use default agent pool

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferFromPosix(buf, tc.ProjectID, sourceAgentPoolName, rootDirectory, gcsSinkBucket)
		defer cleanupSTSJob(resp, tc.ProjectID)

		if err != nil {
			r.Errorf("transfer_from_posix: %#v", err)
		}

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_from_posix: got %q, want %q", got, want)
		}
	})
}

func TestTransferBetweenPosix(t *testing.T) {
	tc := testutil.SystemTest(t)

	rootDirectory, err := ioutil.TempDir("", "transfer-between-posix-test-source")
	if err != nil {
		t.Fatalf("transfer_between_posix: %#v", err)
	}
	defer os.RemoveAll(rootDirectory)

	destinationDirectory, err := ioutil.TempDir("", "transfer-between-posix-test-sink")
	if err != nil {
		t.Fatalf("transfer_between_posix: %#v", err)
	}
	defer os.RemoveAll(destinationDirectory)

	sourceAgentPoolName := "" //use default agent pool
	sinkAgentPoolName := ""   //use default agent pool

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferBetweenPosix(buf, tc.ProjectID, sourceAgentPoolName, sinkAgentPoolName, rootDirectory, destinationDirectory, gcsSinkBucket)
		if err != nil {
			r.Errorf("transfer_between_posix: %#v", err)
		}
		defer cleanupSTSJob(resp, tc.ProjectID)

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_between_posix: got %q, want %q", got, want)
		}
	})
}

func TestTransferUsingManifest(t *testing.T) {
	tc := testutil.SystemTest(t)

	rootDirectory, err := ioutil.TempDir("", "transfer-using-manifest-test")
	if err != nil {
		t.Fatalf("transfer_using_manifest: %#v", err)
	}
	defer os.RemoveAll(rootDirectory)

	sourceAgentPoolName := "" //use default agent pool
	object := sc.Bucket(gcsSourceBucket).Object("manifest.csv")
	defer object.Delete(context.Background())

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferUsingManifest(buf, tc.ProjectID, sourceAgentPoolName, rootDirectory, gcsSinkBucket, gcsSourceBucket, "manifest.csv")
		defer cleanupSTSJob(resp, tc.ProjectID)

		if err != nil {
			r.Errorf("transfer_using_manifest: %#v", err)
		}

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_using_manifest: got %q, want %q", got, want)
		}
	})
}

func TestTransferFromS3CompatibleSource(t *testing.T) {
	tc := testutil.SystemTest(t)

	sourceAgentPoolName := "" //use default agent pool
	sourcePath := ""          //use root directory
	gcsPath := ""             //use root directory

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferFromS3CompatibleSource(buf, tc.ProjectID, sourceAgentPoolName, s3Bucket, sourcePath, gcsSinkBucket, gcsPath)

		if err != nil {
			r.Errorf("transfer_from_s3_compatible_source: %#v", err)
		}
		defer cleanupSTSJob(resp, tc.ProjectID)

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_from_s3_compatible_source: got %q, want %q", got, want)
		}
	})
}

func TestTransferFromAzure(t *testing.T) {
	tc := testutil.SystemTest(t)

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := transferFromAzure(buf, tc.ProjectID, accountName, azureContainer, gcsSinkBucket)
		if err != nil {
			r.Errorf("transfer_from_azure: %#v", err)
		}
		defer cleanupSTSJob(resp, tc.ProjectID)

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("transfer_from_azure: got %q, want %q", got, want)
		}
	})
}

func TestCreateEventDrivenGCSTransfer(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	pubSubTopicId := testutil.UniqueBucketName("pubsubtopic")

	pubsubClient, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
        log.Fatalf("Couldn't create pubsub client: %v", err)
	}
	defer pubsubClient.Close()

	topic, err := pubsubClient.CreateTopic(ctx, pubSubTopicId)
	if err != nil {
		log.Fatal("Coudln't create pubsub topic: " + err.Error())
	}
	defer topic.Delete(ctx)

	policy, err := topic.IAM().Policy(ctx)
	if err != nil {
		log.Fatal("Couldn't get pubsub topic policy: " + err.Error())
	}
	policy.Add("serviceAccount:"+stsServiceAccountEmail, "roles/pubsub.subscriber")
	if err := topic.IAM().SetPolicy(ctx, policy); err != nil {
		log.Fatal("Couldn't set pubsub topic policy: " + err.Error())
	}

	subId := testutil.UniqueBucketName("pubsubsubscription")

	sub, err := pubsubClient.CreateSubscription(ctx, subId, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		log.Fatal("Couldn't create pubsub subscription: " + err.Error())
	}

	pubSubSubscriptionID := sub.String()

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := createEventDrivenGCSTransfer(buf, tc.ProjectID, gcsSourceBucket, gcsSinkBucket, pubSubSubscriptionID)
		if err != nil {
			r.Errorf("create_event_driven_gcs_transfer: %#v", err)
		}
		defer cleanupSTSJob(resp, tc.ProjectID)

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("create_event_driven_gcs_transfer: got %q, want %q", got, want)
		}
	})
}

func TestCreateEventDrivenAWSTransfer(t *testing.T) {
	tc := testutil.SystemTest(t)

	queue := testutil.UniqueBucketName("stssqsqueue")
	sqsClient := sqs.New(sess)
	result, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queue,
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	if err != nil {
		log.Fatalf("couldn't create SQS queue: %v", err)
	}
	defer sqsClient.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: result.QueueUrl,
	})

	attributes, err := sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("QueueArn")},
		QueueUrl:       result.QueueUrl,
	})
	if err != nil {
		log.Fatalf("couldn't get SQS queue attributes: %v", err)
	}

	sqsQueueARN := *attributes.Attributes["QueueArn"]

	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		resp, err := createEventDrivenAWSTransfer(buf, tc.ProjectID, s3Bucket, gcsSinkBucket, sqsQueueARN)
		if err != nil {
			r.Errorf("create_event_driven_aws_transfer: %#v", err)
		}
		defer cleanupSTSJob(resp, tc.ProjectID)

		got := buf.String()
		if want := "transferJobs/"; !strings.Contains(got, want) {
			r.Errorf("create_event_driven_aws_transfer: got %q, want %q", got, want)
		}
	})
}

func grantSTSPermissions(bucketName string, str *storage.Client) {
	ctx := context.Background()

	identity := "serviceAccount:" + stsServiceAccountEmail

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
		log.Fatalf("bucket(%q).IAM().SetPolicy: %v", bucketName, err)
	}
}

func cleanupSTSJob(job *storagetransferpb.TransferJob, projectID string) {
	if job == nil {
		return
	}

	ctx := context.Background()

	tj := &storagetransferpb.TransferJob{
		Name:   job.Name,
		Status: storagetransferpb.TransferJob_DELETED,
	}
	sts.UpdateTransferJob(ctx, &storagetransferpb.UpdateTransferJobRequest{
		JobName:     job.Name,
		ProjectId:   projectID,
		TransferJob: tj,
	})
}
