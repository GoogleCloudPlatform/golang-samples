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

package risk

// [START dlp_k_anonymity_with_entity_id]

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/pubsub"
)

// Uses the Data Loss Prevention API to compute the k-anonymity of a
// column set in a Google BigQuery table.
func calculateKAnonymityWithEntityId(w io.Writer, projectID, bigQueryProjectID, datasetId, tableId, pubSubTopic, pubSubSub string, columnNames ...string) error {
	// projectID := "my-project-id"
	// datasetId := "your-bigquery-dataset-id"
	// tableId := "your-bigquery-table-id"
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// columnNames := "age" "job_title"

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

	// Specify the BigQuery table to analyze
	bigQueryTable := &dlppb.BigQueryTable{
		ProjectId: bigQueryProjectID,
		DatasetId: datasetId,
		TableId:   tableId,
	}

	// Configure the privacy metric for the job
	// Build the QuasiID slice.
	var q []*dlppb.FieldId
	for _, c := range columnNames {
		q = append(q, &dlppb.FieldId{Name: c})
	}

	entityId := &dlppb.EntityId{
		Field: &dlppb.FieldId{
			Name: "id",
		},
	}

	kanonymityConfig := &dlppb.PrivacyMetric_KAnonymityConfig{
		QuasiIds: q,
		EntityId: entityId,
	}

	privacyMetric := &dlppb.PrivacyMetric{
		Type: &dlppb.PrivacyMetric_KAnonymityConfig_{
			KAnonymityConfig: kanonymityConfig,
		},
	}

	// Create action to publish job status notifications over Google Cloud Pub/Sub
	// Create a PubSub Client used to listen for when the inspect job finishes.
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	// Create the Topic if it doesn't exist.
	t := pubsubClient.Topic(pubSubTopic)
	if exists, err := t.Exists(ctx); err != nil {
		fmt.Fprintf(w, "t.Exists: %v", err)
		return err
	} else if !exists {
		if t, err = pubsubClient.CreateTopic(ctx, pubSubTopic); err != nil {
			fmt.Fprintf(w, "CreateTopic: %v", err)
			return err
		}
	}

	// Create the Subscription if it doesn't exist.
	s := pubsubClient.Subscription(pubSubSub)
	if exists, err := s.Exists(ctx); err != nil {
		fmt.Fprintf(w, "s.Exits: %v", err)
		return err
	} else if !exists {
		if s, err = pubsubClient.CreateSubscription(ctx, pubSubSub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			fmt.Fprintf(w, "CreateSubscription: %v", err)
			return err
		}
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + projectID + "/topics/" + pubSubTopic

	publishToPubSub := &dlppb.Action_PublishToPubSub{
		Topic: topic,
	}

	action := &dlppb.Action{
		Action: &dlppb.Action_PubSub{
			PubSub: publishToPubSub,
		},
	}

	// Configure the risk analysis job to perform
	riskAnalysisJobConfig := &dlppb.RiskAnalysisJobConfig{
		PrivacyMetric: privacyMetric,
		SourceTable:   bigQueryTable,
		Actions: []*dlppb.Action{
			action,
		},
	}

	// Build the request to be sent by the client
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: riskAnalysisJobConfig,
		},
	}

	// Send the request to the API using the client
	j, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the risk job to finish by waiting for a PubSub message.
	// This only waits for 10 minutes. For long jobs, consider using a truly
	// asynchronous execution model such as Cloud Functions.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()
		time.Sleep(500 * time.Millisecond)
		j, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			fmt.Fprintf(w, "GetDlpJob: %v", err)
			return
		}
		h := j.GetRiskDetails().GetKAnonymityResult().GetEquivalenceClassHistogramBuckets()
		for i, b := range h {
			fmt.Fprintf(w, "Histogram bucket %v\n", i)
			fmt.Fprintf(w, "  Size range: [%v,%v]\n", b.GetEquivalenceClassSizeLowerBound(), b.GetEquivalenceClassSizeUpperBound())
			fmt.Fprintf(w, "  %v unique values total\n", b.GetBucketSize())
			for _, v := range b.GetBucketValues() {
				var qvs []string
				for _, qv := range v.GetQuasiIdsValues() {
					qvs = append(qvs, qv.String())
				}
				fmt.Fprintf(w, "    QuasiID values: %s\n", strings.Join(qvs, ", "))
				fmt.Fprintf(w, "    Class size: %v\n", v.GetEquivalenceClassSize())
			}
		}
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil

}

// [END dlp_k_anonymity_with_entity_id]
