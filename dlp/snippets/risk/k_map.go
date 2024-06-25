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

// [START dlp_k_map]
import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/ptypes/empty"
)

// riskKMap runs K Map on the given data.
func riskKMap(w io.Writer, projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, region string, columnNames ...string) error {
	// projectID := "my-project-id"
	// dataProject := "bigquery-public-data"
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// datasetID := "san_francisco"
	// tableID := "bikeshare_trips"
	// region := "US"
	// columnNames := "zip_code"
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

	// Build the QuasiID slice.
	var q []*dlppb.PrivacyMetric_KMapEstimationConfig_TaggedField
	for _, c := range columnNames {
		q = append(q, &dlppb.PrivacyMetric_KMapEstimationConfig_TaggedField{
			Field: &dlppb.FieldId{
				Name: c,
			},
			Tag: &dlppb.PrivacyMetric_KMapEstimationConfig_TaggedField_Inferred{
				Inferred: &empty.Empty{},
			},
		})
	}

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				// PrivacyMetric configures what to compute.
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_KMapEstimationConfig_{
						KMapEstimationConfig: &dlppb.PrivacyMetric_KMapEstimationConfig{
							QuasiIds:   q,
							RegionCode: region,
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
		h := j.GetRiskDetails().GetKMapEstimationResult().GetKMapEstimationHistogram()
		for i, b := range h {
			fmt.Fprintf(w, "Histogram bucket %v\n", i)
			fmt.Fprintf(w, "  Anonymity range: [%v,%v]\n", b.GetMaxAnonymity(), b.GetMaxAnonymity())
			fmt.Fprintf(w, "  %v unique values total\n", b.GetBucketSize())
			for _, v := range b.GetBucketValues() {
				var qvs []string
				for _, qv := range v.GetQuasiIdsValues() {
					qvs = append(qvs, qv.String())
				}
				fmt.Fprintf(w, "    QuasiID values: %s\n", strings.Join(qvs, ", "))
				fmt.Fprintf(w, "    Estimated anonymity: %v\n", v.GetEstimatedAnonymity())
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

// [END dlp_k_map]
