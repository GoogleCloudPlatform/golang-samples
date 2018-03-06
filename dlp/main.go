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

	dlp "cloud.google.com/go/dlp/apiv2"
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

func riskNumerical(w io.Writer, client *dlp.Client, project, dataProject, datasetID, tableID, columnName string) {
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
			},
		},
	}
	j, err := client.CreateDlpJob(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, j)
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
		riskNumerical(os.Stdout, client, *project, flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
	default:
		fmt.Fprintf(os.Stderr, `Usage: %s CMD "string"\n`, os.Args[0])
		os.Exit(1)
	}
}
