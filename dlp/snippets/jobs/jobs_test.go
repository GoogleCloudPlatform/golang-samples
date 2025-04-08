// Copyright 2019 Google LLC
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

// Package jobs contains example snippets using the DLP jobs API.
package jobs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/gofrs/uuid"
)

// setupPubSub creates a subscription to the given topic.
func setupPubSub(projectID, topic, sub string) (*pubsub.Subscription, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	// Create the Topic if it doesn't exist.
	t := client.Topic(topic)
	if exists, err := t.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking PubSub topic: %w", err)
	} else if !exists {
		if t, err = client.CreateTopic(ctx, topic); err != nil {
			return nil, fmt.Errorf("error creating PubSub topic: %w", err)
		}
	}

	// Create the Subscription if it doesn't exist.
	s := client.Subscription(sub)
	if exists, err := s.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking for subscription: %w", err)
	} else if !exists {
		if s, err = client.CreateSubscription(ctx, sub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}
	}

	return s, nil
}

// riskNumerical computes the numerical risk of the given column.
func riskNumerical(projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, columnName string) error {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	// Create a PubSub Client used to listen for when the inspect job finishes.
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Error creating PubSub client: %w", err)
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	s, err := setupPubSub(projectID, pubSubTopic, pubSubSub)
	if err != nil {
		return fmt.Errorf("setupPubSub: %w", err)
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + projectID + "/topics/" + pubSubTopic

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				// PrivacyMetric configures what to compute.
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_NumericalStatsConfig_{
						NumericalStatsConfig: &dlppb.PrivacyMetric_NumericalStatsConfig{
							Field: &dlppb.FieldId{
								Name: columnName,
							},
						},
					},
				},
				// SourceTable describes where to find the data.
				SourceTable: &dlppb.BigQueryTable{
					ProjectId: dataProject,
					DatasetId: datasetID,
					TableId:   tableID,
				},
				// Send a message to PubSub using Actions.
				Actions: []*dlppb.Action{
					{
						Action: &dlppb.Action_PubSub{
							PubSub: &dlppb.Action_PublishToPubSub{
								Topic: topic,
							},
						},
					},
				},
			},
		},
	}
	// Create the risk job.
	j, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateDlpJob: %w", err)
	}

	// Wait for the risk job to finish by waiting for a PubSub message.
	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		return fmt.Errorf("Recieve: %w", err)
	}
	return nil
}

var jobIDRegexp = regexp.MustCompile(`Job ([^ ]+) status.*`)

func createBucketForCreatJob(t *testing.T, projectID string) (string, string, error) {
	t.Helper()

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", "", err
	}
	defer client.Close()
	u := uuid.Must(uuid.NewV4()).String()[:8]
	bucketName := "dlp-job-go-lang-test" + u

	// Check if the bucket already exists.
	bucketExists := false
	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err == nil {
		bucketExists = true
	}

	// If the bucket doesn't exist, create it.
	if !bucketExists {
		if err := client.Bucket(bucketName).Create(ctx, projectID, &storage.BucketAttrs{
			StorageClass: "STANDARD",
			Location:     "us-central1",
		}); err != nil {
			log.Fatalf("---Failed to create bucket: %v", err)
		}
		fmt.Printf("---Bucket '%s' created successfully.\n", bucketName)
	} else {
		fmt.Printf("---Bucket '%s' already exists.\n", bucketName)
	}

	filePathToUpload := "testdata/test.txt"

	// Open local file.
	file, err := ioutil.ReadFile(filePathToUpload)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Get a reference to the bucket
	bucket := client.Bucket(bucketName)

	// Upload the file
	u = uuid.Must(uuid.NewV4()).String()[:8]
	fileName := "test" + u + ".txt"
	object := bucket.Object(fileName)
	writer := object.NewWriter(ctx)
	_, err = writer.Write(file)
	if err != nil {
		log.Fatalf("---Failed to write file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("---Failed to close writer: %v", err)
	}
	fmt.Printf("---File uploaded successfully: %v\n", fileName)

	// Check if the file exists in the bucket
	_, err = bucket.Object(fileName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			fmt.Printf("---File %v does not exist in bucket %v\n", fileName, bucketName)
		} else {
			log.Fatalf("---Failed to check file existence: %v", err)
		}
	} else {
		fmt.Printf("---File %v exists in bucket %v\n", fileName, bucketName)
	}

	return bucketName, fileName, nil
}

func deleteAssetsOfCreateJobTest(t *testing.T, projectID, bucketName, objectName string) error {
	t.Helper()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	o := client.Bucket(bucketName).Object(objectName)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		t.Fatal(err)
	}

	bucket := client.Bucket(bucketName)
	if err := bucket.Delete(ctx); err != nil {
		t.Fatal(err)
	}
	return nil
}
