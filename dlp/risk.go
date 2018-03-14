/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/ptypes/empty"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func setupPubSub(ctx context.Context, client *pubsub.Client, project, topic, sub string) (*pubsub.Subscription, error) {
	t := client.Topic(topic)
	if exists, err := t.Exists(ctx); err != nil {
		return nil, fmt.Errorf("error checking PubSub topic: %v", err)
	} else if !exists {
		if t, err = client.CreateTopic(ctx, topic); err != nil {
			return nil, fmt.Errorf("error creating PubSub topic: %v", err)
		}
	}

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

// [START dlp_numerical_stats]
func riskNumerical(w io.Writer, client *dlp.Client, project, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, columnName string) {
	ctx := context.Background()
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}
	topic := "projects/" + project + "/topics/" + pubSubTopic
	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_NumericalStatsConfig_{
						NumericalStatsConfig: &dlppb.PrivacyMetric_NumericalStatsConfig{
							Field: &dlppb.FieldId{
								Name: columnName,
							},
						},
					},
				},
				SourceTable: &dlppb.BigQueryTable{
					ProjectId: dataProject,
					DatasetId: datasetID,
					TableId:   tableID,
				},
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
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j)

	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		if msg.Attributes["DlpJobName"] == j.GetName() {
			jr, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
				Name: j.GetName(),
			})
			if err != nil {
				log.Fatalf("Error getting completed job: %v\n", err)
			}
			n := jr.GetRiskDetails().GetNumericalStatsResult()
			fmt.Fprintf(w, "Value range: [%v, %v]\n", n.GetMinValue(), n.GetMaxValue())
			var tmp string
			for p, v := range n.GetQuantileValues() {
				if v.String() != tmp {
					fmt.Fprintf(w, "Value at %v quantile: %v\n", p, v)
					tmp = v.String()
				}
			}
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_numerical_stats]

// [START dlp_categorical_stats]
func riskCategorical(w io.Writer, client *dlp.Client, project, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, columnName string) {
	ctx := context.Background()
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}
	topic := "projects/" + project + "/topics/" + pubSubTopic
	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_CategoricalStatsConfig_{
						CategoricalStatsConfig: &dlppb.PrivacyMetric_CategoricalStatsConfig{
							Field: &dlppb.FieldId{
								Name: columnName,
							},
						},
					},
				},
				SourceTable: &dlppb.BigQueryTable{
					ProjectId: dataProject,
					DatasetId: datasetID,
					TableId:   tableID,
				},
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
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j)

	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		if msg.Attributes["DlpJobName"] == j.GetName() {
			jr, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
				Name: j.GetName(),
			})
			if err != nil {
				log.Fatalf("Error getting completed job: %v\n", err)
			}
			h := jr.GetRiskDetails().GetCategoricalStatsResult().GetValueFrequencyHistogramBuckets()
			for i, b := range h {
				fmt.Fprintf(w, "Histogram bucket %v\n", i)
				fmt.Fprintf(w, "  Most common value occurs %v times\n", b.GetValueFrequencyUpperBound())
				fmt.Fprintf(w, "  Least common value occurs %v times\n", b.GetValueFrequencyLowerBound())
				fmt.Fprintf(w, "  %v unique values total\n", b.GetBucketSize())
				for _, v := range b.GetBucketValues() {
					fmt.Fprintf(w, "    Value %v occurs %v times\n", v.GetValue(), v.GetCount())
				}
			}
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_categorical_stats]

// [START dlp_k_anonymity]
func riskKAnonymity(w io.Writer, client *dlp.Client, project, dataProject, pubSubTopic, pubSubSub, datasetID, tableID string, columnNames ...string) {
	ctx := context.Background()
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}
	topic := "projects/" + project + "/topics/" + pubSubTopic

	// Build the QuasiID slice.
	var q []*dlppb.FieldId
	for _, c := range columnNames {
		q = append(q, &dlppb.FieldId{Name: c})
	}

	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_KAnonymityConfig_{
						KAnonymityConfig: &dlppb.PrivacyMetric_KAnonymityConfig{
							QuasiIds: q,
						},
					},
				},
				SourceTable: &dlppb.BigQueryTable{
					ProjectId: dataProject,
					DatasetId: datasetID,
					TableId:   tableID,
				},
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
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j)

	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		if msg.Attributes["DlpJobName"] == j.GetName() {
			j, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
				Name: j.GetName(),
			})
			if err != nil {
				log.Fatalf("Error getting completed job: %v\n", err)
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
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_k_anonymity]

// [START dlp_l_diversity]
func riskLDiversity(w io.Writer, client *dlp.Client, project, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, sensitiveAttribute string, columnNames ...string) {
	ctx := context.Background()
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}
	topic := "projects/" + project + "/topics/" + pubSubTopic

	// Build the QuasiID slice.
	var q []*dlppb.FieldId
	for _, c := range columnNames {
		q = append(q, &dlppb.FieldId{Name: c})
	}

	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
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
				SourceTable: &dlppb.BigQueryTable{
					ProjectId: dataProject,
					DatasetId: datasetID,
					TableId:   tableID,
				},
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
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j)

	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		if msg.Attributes["DlpJobName"] == j.GetName() {
			j, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
				Name: j.GetName(),
			})
			if err != nil {
				log.Fatalf("Error getting completed job: %v\n", err)
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
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_l_diversity]

// [START k_map]
func riskKMap(w io.Writer, client *dlp.Client, project, dataProject, pubSubTopic, pubSubSub, datasetID, tableID, region string, columnNames ...string) {
	ctx := context.Background()
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}
	topic := "projects/" + project + "/topics/" + pubSubTopic

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

	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: &dlppb.RiskAnalysisJobConfig{
				PrivacyMetric: &dlppb.PrivacyMetric{
					Type: &dlppb.PrivacyMetric_KMapEstimationConfig_{
						KMapEstimationConfig: &dlppb.PrivacyMetric_KMapEstimationConfig{
							QuasiIds:   q,
							RegionCode: region,
						},
					},
				},
				SourceTable: &dlppb.BigQueryTable{
					ProjectId: dataProject,
					DatasetId: datasetID,
					TableId:   tableID,
				},
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
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j)

	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		if msg.Attributes["DlpJobName"] == j.GetName() {
			j, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
				Name: j.GetName(),
			})
			if err != nil {
				log.Fatalf("Error getting completed job: %v\n", err)
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
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END k_map]
