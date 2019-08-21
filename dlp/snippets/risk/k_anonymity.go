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

// [START dlp_k_anonymity]
import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// riskKAnonymity computes the risk of the given columns using K Anonymity.
func riskKAnonymity(w io.Writer, projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID string, columnNames ...string) error {
	// projectID := "my-project-id"
	// dataProject := "bigquery-public-data"
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// datasetID := "nhtsa_traffic_fatalities"
	// tableID := "accident_2015"
	// columnNames := "state_number" "county"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	s, err := setupPubSub(projectID, pubSubTopic, pubSubSub)
	if err != nil {
		return fmt.Errorf("setupPubSub: %v", err)
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
		Parent: "projects/" + projectID,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				// PrivacyMetric configures what to compute.
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_KAnonymityConfig_{
						KAnonymityConfig: &dlppb.PrivacyMetric_KAnonymityConfig{
							QuasiIds: q,
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
	j, err := client.CreateDlpJob(context.Background(), req)
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
		j, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			log.Fatalf("GetDlpJob: %v", err)
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

// [END dlp_k_anonymity]
