// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_inspect_string]

// inspectString searches for the given infoTypes in the input.
func inspectString(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, customDictionaries []string, customRegexes []string, input string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Convert the custom dictionary word lists and custom regexes to a list of CustomInfoTypes.
	var customInfoTypes []*dlppb.CustomInfoType
	for idx, it := range customDictionaries {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_DICTIONARY_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Dictionary_{
				Dictionary: &dlppb.CustomInfoType_Dictionary{
					Source: &dlppb.CustomInfoType_Dictionary_WordList_{
						WordList: &dlppb.CustomInfoType_Dictionary_WordList{
							Words: strings.Split(it, ","),
						},
					},
				},
			},
		})
	}
	for idx, it := range customRegexes {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_REGEX_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Regex_{
				Regex: &dlppb.CustomInfoType_Regex{
					Pattern: it,
				},
			},
		})
	}
	// Create a configured request.
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:       i,
			CustomInfoTypes: customInfoTypes,
			MinLikelihood:   minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		// The item to analyze.
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}
	// Send the request.
	resp, err := client.InspectContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Print the result.
	fmt.Fprintln(w, resp.GetResult())
}

// [END dlp_inspect_string]

// [START dlp_inspect_file]

// inspectFile searches for the given info types in the given Reader (with the given bytesType).
func inspectFile(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, customDictionaries []string, customRegexes []string, bytesType dlppb.ByteContentItem_BytesType, input io.Reader) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Convert the custom dictionary word lists and custom regexes to a list of CustomInfoTypes.
	var customInfoTypes []*dlppb.CustomInfoType
	for idx, it := range customDictionaries {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_DICTIONARY_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Dictionary_{
				Dictionary: &dlppb.CustomInfoType_Dictionary{
					Source: &dlppb.CustomInfoType_Dictionary_WordList_{
						WordList: &dlppb.CustomInfoType_Dictionary_WordList{
							Words: strings.Split(it, ","),
						},
					},
				},
			},
		})
	}
	for idx, it := range customRegexes {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_REGEX_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Regex_{
				Regex: &dlppb.CustomInfoType_Regex{
					Pattern: it,
				},
			},
		})
	}
	b, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	// Create a configured request.
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:       i,
			CustomInfoTypes: customInfoTypes,
			MinLikelihood:   minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		// The item to analyze.
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_ByteItem{
				ByteItem: &dlppb.ByteContentItem{
					Type: bytesType,
					Data: b,
				},
			},
		},
	}
	// Send the request.
	resp, err := client.InspectContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Print the result.
	fmt.Fprintln(w, resp.GetResult())
}

// [END dlp_inspect_file]

// [START dlp_inspect_gcs]

// inspectGCSFile searches for the given info types in the given file.
func inspectGCSFile(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, customDictionaries []string, customRegexes []string, pubSubTopic, pubSubSub, bucketName, fileName string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Convert the custom dictionary word lists and custom regexes to a list of CustomInfoTypes.
	var customInfoTypes []*dlppb.CustomInfoType
	for idx, it := range customDictionaries {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_DICTIONARY_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Dictionary_{
				Dictionary: &dlppb.CustomInfoType_Dictionary{
					Source: &dlppb.CustomInfoType_Dictionary_WordList_{
						WordList: &dlppb.CustomInfoType_Dictionary_WordList{
							Words: strings.Split(it, ","),
						},
					},
				},
			},
		})
	}
	for idx, it := range customRegexes {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_REGEX_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Regex_{
				Regex: &dlppb.CustomInfoType_Regex{
					Pattern: it,
				},
			},
		})
	}

	ctx := context.Background()

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + project + "/topics/" + pubSubTopic

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
				// StorageConfig describes where to find the data.
				StorageConfig: &dlppb.StorageConfig{
					Type: &dlppb.StorageConfig_CloudStorageOptions{
						CloudStorageOptions: &dlppb.CloudStorageOptions{
							FileSet: &dlppb.CloudStorageOptions_FileSet{
								Url: "gs://" + bucketName + "/" + fileName,
							},
						},
					},
				},
				// InspectConfig describes what fields to look for.
				InspectConfig: &dlppb.InspectConfig{
					InfoTypes:       i,
					CustomInfoTypes: customInfoTypes,
					MinLikelihood:   minLikelihood,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: maxFindings,
					},
					IncludeQuote: includeQuote,
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
	// Create the inspect job.
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the inspect job to finish by waiting for a PubSub message.
	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()
		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			log.Fatalf("Error getting completed job: %v\n", err)
		}
		r := resp.GetInspectDetails().GetResult().GetInfoTypeStats()
		if len(r) == 0 {
			fmt.Fprintf(w, "No results")
		}
		for _, s := range r {
			fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
		}
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_inspect_gcs]

// [START dlp_inspect_datastore]

// inspectDatastore searches for the given info types in the given dataset kind.
func inspectDatastore(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, customDictionaries []string, customRegexes []string, pubSubTopic, pubSubSub, dataProject, namespaceID, kind string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Convert the custom dictionary word lists and custom regexes to a list of CustomInfoTypes.
	var customInfoTypes []*dlppb.CustomInfoType
	for idx, it := range customDictionaries {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_DICTIONARY_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Dictionary_{
				Dictionary: &dlppb.CustomInfoType_Dictionary{
					Source: &dlppb.CustomInfoType_Dictionary_WordList_{
						WordList: &dlppb.CustomInfoType_Dictionary_WordList{
							Words: strings.Split(it, ","),
						},
					},
				},
			},
		})
	}
	for idx, it := range customRegexes {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_REGEX_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Regex_{
				Regex: &dlppb.CustomInfoType_Regex{
					Pattern: it,
				},
			},
		})
	}

	ctx := context.Background()

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + project + "/topics/" + pubSubTopic

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
				// StorageConfig describes where to find the data.
				StorageConfig: &dlppb.StorageConfig{
					Type: &dlppb.StorageConfig_DatastoreOptions{
						DatastoreOptions: &dlppb.DatastoreOptions{
							PartitionId: &dlppb.PartitionId{
								ProjectId:   dataProject,
								NamespaceId: namespaceID,
							},
							Kind: &dlppb.KindExpression{
								Name: kind,
							},
						},
					},
				},
				// InspectConfig describes what fields to look for.
				InspectConfig: &dlppb.InspectConfig{
					InfoTypes:       i,
					CustomInfoTypes: customInfoTypes,
					MinLikelihood:   minLikelihood,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: maxFindings,
					},
					IncludeQuote: includeQuote,
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
	// Create the inspect job.
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the inspect job to finish by waiting for a PubSub message.
	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()
		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			log.Fatalf("Error getting completed job: %v\n", err)
		}
		r := resp.GetInspectDetails().GetResult().GetInfoTypeStats()
		if len(r) == 0 {
			fmt.Fprintf(w, "No results")
		}
		for _, s := range r {
			fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
		}
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_inspect_datastore]

// [START dlp_inspect_bigquery]

// inspectBigquery searches for the given info types in the given Bigquery dataset table.
func inspectBigquery(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, customDictionaries []string, customRegexes []string, pubSubTopic, pubSubSub, dataProject, datasetID, tableID string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Convert the custom dictionary word lists and custom regexes to a list of CustomInfoTypes.
	var customInfoTypes []*dlppb.CustomInfoType
	for idx, it := range customDictionaries {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_DICTIONARY_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Dictionary_{
				Dictionary: &dlppb.CustomInfoType_Dictionary{
					Source: &dlppb.CustomInfoType_Dictionary_WordList_{
						WordList: &dlppb.CustomInfoType_Dictionary_WordList{
							Words: strings.Split(it, ","),
						},
					},
				},
			},
		})
	}
	for idx, it := range customRegexes {
		customInfoTypes = append(customInfoTypes, &dlppb.CustomInfoType{
			InfoType: &dlppb.InfoType{
				Name: fmt.Sprintf("CUSTOM_REGEX_%d", idx),
			},
			Type: &dlppb.CustomInfoType_Regex_{
				Regex: &dlppb.CustomInfoType_Regex{
					Pattern: it,
				},
			},
		})
	}

	ctx := context.Background()

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pClient, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	defer pClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	s, err := setupPubSub(ctx, pClient, project, pubSubTopic, pubSubSub)
	if err != nil {
		log.Fatalf("Error setting up PubSub: %v\n", err)
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + project + "/topics/" + pubSubTopic

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: "projects/" + project,
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
				// StorageConfig describes where to find the data.
				StorageConfig: &dlppb.StorageConfig{
					Type: &dlppb.StorageConfig_BigQueryOptions{
						BigQueryOptions: &dlppb.BigQueryOptions{
							TableReference: &dlppb.BigQueryTable{
								ProjectId: dataProject,
								DatasetId: datasetID,
								TableId:   tableID,
							},
						},
					},
				},
				// InspectConfig describes what fields to look for.
				InspectConfig: &dlppb.InspectConfig{
					InfoTypes:       i,
					CustomInfoTypes: customInfoTypes,
					MinLikelihood:   minLikelihood,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: maxFindings,
					},
					IncludeQuote: includeQuote,
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
	// Create the inspect job.
	j, err := client.CreateDlpJob(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the inspect job to finish by waiting for a PubSub message.
	ctx, cancel := context.WithCancel(ctx)
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()
		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			log.Fatalf("Error getting completed job: %v\n", err)
		}
		r := resp.GetInspectDetails().GetResult().GetInfoTypeStats()
		if len(r) == 0 {
			fmt.Fprintf(w, "No results")
		}
		for _, s := range r {
			fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
		}
		// Stop listening for more messages.
		cancel()
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_inspect_bigquery]
