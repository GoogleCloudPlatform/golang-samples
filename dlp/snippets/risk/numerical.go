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

// [START dlp_numerical_stats]
import (
	"context"
	"fmt"
	"io"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/pubsub"
)

// riskNumerical computes the numerical risk of the given column.
func riskNumerical(w io.Writer, projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, columnName string) error {
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
		return fmt.Errorf("dlp.NewClient: %w", err)
	}

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	// Create the Topic if it doesn't exist.
	t := pubsubClient.Topic(pubSubTopic)
	topicExists, err := t.Exists(ctx)
	if err != nil {
		return err
	}
	if !topicExists {
		if t, err = pubsubClient.CreateTopic(ctx, pubSubTopic); err != nil {
			return err
		}
	}

	// Create the Subscription if it doesn't exist.
	s := pubsubClient.Subscription(pubSubSub)
	subExists, err := s.Exists(ctx)
	if err != nil {
		return err
	}
	if !subExists {
		if s, err = pubsubClient.CreateSubscription(ctx, pubSubSub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return err
		}
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
		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			fmt.Fprintf(w, "GetDlpJob: %v", err)
			return
		}
		n := resp.GetRiskDetails().GetNumericalStatsResult()
		fmt.Fprintf(w, "Value range: [%v, %v]\n", n.GetMinValue(), n.GetMaxValue())
		var tmp string
		for p, v := range n.GetQuantileValues() {
			if v.String() != tmp {
				fmt.Fprintf(w, "Value at %v quantile: %v\n", p, v)
				tmp = v.String()
			}
		}
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		return fmt.Errorf("Recieve: %w", err)
	}
	return nil
}

// [END dlp_numerical_stats]
