// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// dlp is an example of using the DLP API.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func inspect(w io.Writer, client *dlp.Client, project, s string) {
	rcr := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "US_SOCIAL_SECURITY_NUMBER",
				},
			},
			MinLikelihood: dlppb.Likelihood_LIKELIHOOD_UNSPECIFIED,
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.InspectContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetResult())
}

func redact(w io.Writer, client *dlp.Client, project, s string) {
	rcr := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "US_SOCIAL_SECURITY_NUMBER",
				},
			},
			MinLikelihood: dlppb.Likelihood_LIKELIHOOD_UNSPECIFIED,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_RedactConfig{
									RedactConfig: &dlppb.RedactConfig{},
								},
							},
						},
					},
				},
			},
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem())
}

func infoTypes(w io.Writer, client *dlp.Client, filter string) {
	rcr := &dlppb.ListInfoTypesRequest{
		Filter: filter,
	}
	r, err := client.ListInfoTypes(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	for _, it := range r.GetInfoTypes() {
		fmt.Fprintln(w, it.GetName())
	}
}

func mask(w io.Writer, client *dlp.Client, project, s string) {
	rcr := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CharacterMaskConfig{
									CharacterMaskConfig: &dlppb.CharacterMaskConfig{
										MaskingCharacter: "*",
									},
								},
							},
						},
					},
				},
			},
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem().GetValue())
}

func deidentifyFPE(w io.Writer, client *dlp.Client, project, s, wrappedKey, cryptoKeyName string) {
	rcr := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
									CryptoReplaceFfxFpeConfig: &dlppb.CryptoReplaceFfxFpeConfig{
										CryptoKey: &dlppb.CryptoKey{
											Source: &dlppb.CryptoKey_KmsWrapped{
												KmsWrapped: &dlppb.KmsWrappedCryptoKey{
													WrappedKey:    []byte(wrappedKey),
													CryptoKeyName: cryptoKeyName,
												},
											},
										},
										Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
											CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_ALPHA_NUMERIC,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem().GetValue())
}

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
	rcr := &dlppb.CreateDlpJobRequest{
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
	j, err := client.CreateDlpJob(context.Background(), rcr)
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
	rcr := &dlppb.CreateDlpJobRequest{
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
	j, err := client.CreateDlpJob(context.Background(), rcr)
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

	rcr := &dlppb.CreateDlpJobRequest{
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
	j, err := client.CreateDlpJob(context.Background(), rcr)
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

func main() {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	project := flag.String("project", "", "GCloud project ID")
	flag.Parse()

	if *project == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "inspect":
		inspect(os.Stdout, client, *project, flag.Arg(1))
	case "redact":
		redact(os.Stdout, client, *project, flag.Arg(1))
	case "infoTypes":
		infoTypes(os.Stdout, client, flag.Arg(1))
	case "mask":
		mask(os.Stdout, client, *project, flag.Arg(1))
	case "deidfpe":
		deidentifyFPE(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3))
	case "riskNumerical":
		// For example:
		// dlp -project my-project riskNumerical bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
		riskNumerical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskCategorical":
		// For example:
		// dlp -project my-project riskCategorical bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
		riskCategorical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	case "riskKAnonymity":
		// For example:
		// dlp -project my-project riskKAnonymity bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
		riskKAnonymity(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5), flag.Arg(6))
	default:
		fmt.Fprintf(os.Stderr, `Usage: %s CMD "string"\n`, os.Args[0])
		os.Exit(1)
	}
}
