// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package inspect

// [START dlp_inspect_bigquery_with_sampling]
import (
	"context"
	"fmt"
	"io"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/pubsub"
)

// inspectBigQueryTableWithSampling inspect bigQueries for sensitive data with sampling
func inspectBigQueryTableWithSampling(w io.Writer, projectID, topicID, subscriptionID string) error {
	// projectId := "your-project-id"
	// topicID := "your-pubsub-topic-id"
	// or provide a topicID name to create one
	// subscriptionID := "your-pubsub-subscription-id"
	// or provide a subscription name to create one

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the BigQuery table to be inspected.
	tableReference := &dlppb.BigQueryTable{
		ProjectId: "bigquery-public-data",
		DatasetId: "usa_names",
		TableId:   "usa_1910_current",
	}

	bigQueryOptions := &dlppb.BigQueryOptions{
		TableReference: tableReference,
		RowsLimit:      int64(10000),
		SampleMethod:   dlppb.BigQueryOptions_RANDOM_START,
		IdentifyingFields: []*dlppb.FieldId{
			{Name: "name"},
		},
	}

	// Provide storage config with BigqueryOptions
	storageConfig := &dlppb.StorageConfig{
		Type: &dlppb.StorageConfig_BigQueryOptions{
			BigQueryOptions: bigQueryOptions,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoTypes := []*dlppb.InfoType{
		{Name: "PERSON_NAME"},
	}

	// Specify how the content should be inspected.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes:    infoTypes,
		IncludeQuote: true,
	}

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	// Create the Topic if it doesn't exist.
	t := pubsubClient.Topic(topicID)
	if exists, err := t.Exists(ctx); err != nil {
		return err
	} else if !exists {
		if t, err = pubsubClient.CreateTopic(ctx, topicID); err != nil {
			return err
		}
	}

	// Create the Subscription if it doesn't exist.
	s := pubsubClient.Subscription(subscriptionID)
	if exists, err := s.Exists(ctx); err != nil {
		return err
	} else if !exists {
		if s, err = pubsubClient.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return err
		}
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)

	action := &dlppb.Action{
		Action: &dlppb.Action_PubSub{
			PubSub: &dlppb.Action_PublishToPubSub{
				Topic: topic,
			},
		},
	}

	// Configure the long running job we want the service to perform.
	inspectJobConfig := &dlppb.InspectJobConfig{
		StorageConfig: storageConfig,
		InspectConfig: inspectConfig,
		Actions: []*dlppb.Action{
			action,
		},
	}

	// Create the request for the job configured above.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: inspectJobConfig,
		},
	}

	// Use the client to send the request.
	j, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Job Created: %v", j.GetName())

	// Wait for the inspect job to finish by waiting for a PubSub message.
	// This only waits for 10 minutes. For long jobs, consider using a truly
	// asynchronous execution model such as Cloud Functions.
	c, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	err = s.Receive(c, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()

		// Stop listening for more messages.
		defer cancel()
	})
	if err != nil {
		return err
	}

	resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
		Name: j.GetName(),
	})
	if err != nil {
		return err
	}
	r := resp.GetInspectDetails().GetResult().GetInfoTypeStats()
	if len(r) == 0 {
		fmt.Fprintf(w, "No results")
		return err
	}
	for _, s := range r {
		fmt.Fprintf(w, "\nFound %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
	}
	return nil

}

// [END dlp_inspect_bigquery_with_sampling]
