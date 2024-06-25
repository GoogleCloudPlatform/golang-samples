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

// [START dlp_l_diversity]
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

// riskLDiversity computes the L Diversity of the given columns.
func riskLDiversity(w io.Writer, projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, sensitiveAttribute string, columnNames ...string) error {
	// projectID := "my-project-id"
	// dataProject := "bigquery-public-data"
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// datasetID := "nhtsa_traffic_fatalities"
	// tableID := "accident_2015"
	// sensitiveAttribute := "city"
	// columnNames := "state_number", "county"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

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

	// Build the QuasiID slice.
	var q []*dlppb.FieldId
	for _, c := range columnNames {
		q = append(q, &dlppb.FieldId{Name: c})
	}

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				// PrivacyMetric configures what to compute.
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_LDiversityConfig_{
						LDiversityConfig: &dlppb.PrivacyMetric_LDiversityConfig{
							QuasiIds: q,
							SensitiveAttribute: &dlppb.FieldId{
								Name: sensitiveAttribute,
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
		j, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			fmt.Fprintf(w, "GetDlpJob: %v", err)
			return
		}
		h := j.GetRiskDetails().GetLDiversityResult().GetSensitiveValueFrequencyHistogramBuckets()
		for i, b := range h {
			fmt.Fprintf(w, "Histogram bucket %v\n", i)
			fmt.Fprintf(w, "  Size range: [%v,%v]\n", b.GetSensitiveValueFrequencyLowerBound(), b.GetSensitiveValueFrequencyUpperBound())
			fmt.Fprintf(w, "  %v unique values total\n", b.GetBucketSize())
			for _, v := range b.GetBucketValues() {
				var qvs []string
				for _, qv := range v.GetQuasiIdsValues() {
					qvs = append(qvs, qv.String())
				}
				fmt.Fprintf(w, "    QuasiID values: %s\n", strings.Join(qvs, ", "))
				fmt.Fprintf(w, "    Class size: %v\n", v.GetEquivalenceClassSize())
				for _, sv := range v.GetTopSensitiveValues() {
					fmt.Fprintf(w, "    Sensitive value %v occurs %v times\n", sv.GetValue(), sv.GetCount())
				}
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

// [END dlp_l_diversity]
