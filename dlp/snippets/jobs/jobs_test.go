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
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// setupPubSub creates a subscription to the given topic.
func setupPubSub(projectID, topic, sub string) (*pubsub.Subscription, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	// Create the Topic if it doesn't exist.
	t := client.Topic(topic)
	if exists, err := t.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking PubSub topic: %v", err)
	} else if !exists {
		if t, err = client.CreateTopic(ctx, topic); err != nil {
			return nil, fmt.Errorf("error creating PubSub topic: %v", err)
		}
	}

	// Create the Subscription if it doesn't exist.
	s := client.Subscription(sub)
	if exists, err := s.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking for subscription: %v", err)
	} else if !exists {
		if s, err = client.CreateSubscription(ctx, sub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return nil, fmt.Errorf("failed to create subscription: %v", err)
		}
	}

	return s, nil
}

// riskNumerical computes the numerical risk of the given column.
func riskNumerical(projectID, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, columnName string) error {
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
		return fmt.Errorf("CreateDlpJob: %v", err)
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
		return fmt.Errorf("Recieve: %v", err)
	}
	return nil
}

func TestListJobs(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	listJobs(buf, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(tc.ProjectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		err := listJobs(buf, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
		if err != nil {
			t.Errorf("listJobs(%s, %s, %s) = error %q, want nil", buf, tc.ProjectID, "", err)
		}
		s = buf.String()
	}
	if !strings.Contains(buf.String(), "Job") {
		t.Errorf("%q not found in listJobs output: %q", "Job", s)
	}
}

var jobIDRegexp = regexp.MustCompile(`Job ([^ ]+) status.*`)

func TestDeleteJob(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	listJobs(buf, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
	s := buf.String()
	if len(s) == 0 {
		// Create job.
		riskNumerical(tc.ProjectID, "bigquery-public-data", "risk-topic", "risk-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		buf.Reset()
		listJobs(buf, tc.ProjectID, "", "RISK_ANALYSIS_JOB")
		s = buf.String()
	}
	jobName := string(jobIDRegexp.FindSubmatch([]byte(s))[1])
	buf.Reset()
	deleteJob(buf, jobName)
	if got := buf.String(); got != "Successfully deleted job" {
		t.Errorf("unable to delete job: %s", s)
	}
}
