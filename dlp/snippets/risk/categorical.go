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

package risk

// [START dlp_categorical_stats]
import (
	"context"
	"fmt"
	"io"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// riskCategorical computes the categorical risk of the given data.
func riskCategorical(w io.Writer, projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, columnName string) error {
	// projectID := "my-project-id"
	// dataProject := "bigquery-public-data"
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// datasetID := "nhtsa_traffic_fatalities"
	// tableID := "accident_2015"
	// columnName := "state_number"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}
	// Create a PubSub Client used to listen for when the inspect job finishes.
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Error creating PubSub client: %v", err)
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	s, err := setupPubSub(projectID, pubSubTopic, pubSubSub)
	if err != nil {
		return fmt.Errorf("setupPubSub: %v", err)
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + projectID + "/topics/" + pubSubTopic

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + projectID,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				// PrivacyMetric configures what to compute.
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_CategoricalStatsConfig_{
						CategoricalStatsConfig: &dlppb.PrivacyMetric_CategoricalStatsConfig{
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
		return fmt.Errorf("CreateDlpJob: %v", err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the risk job to finish by waiting for a PubSub message.
	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()
		time.Sleep(500 * time.Millisecond)
		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			fmt.Fprintf(w, "GetDlpJob: %v", err)
			return
		}
		h := resp.GetRiskDetails().GetCategoricalStatsResult().GetValueFrequencyHistogramBuckets()
		for i, b := range h {
			fmt.Fprintf(w, "Histogram bucket %v\n", i)
			fmt.Fprintf(w, "  Most common value occurs %v times\n", b.GetValueFrequencyUpperBound())
			fmt.Fprintf(w, "  Least common value occurs %v times\n", b.GetValueFrequencyLowerBound())
			fmt.Fprintf(w, "  %v unique values total\n", b.GetBucketSize())
			for _, v := range b.GetBucketValues() {
				fmt.Fprintf(w, "    Value %v occurs %v times\n", v.GetValue(), v.GetCount())
			}
		}
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		return fmt.Errorf("Recieve: %v", err)
	}
	return nil
}

// [END dlp_categorical_stats]
