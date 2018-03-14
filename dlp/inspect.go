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
	"io/ioutil"
	"log"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_inspect_string]
func inspect(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, input string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     i,
			MinLikelihood: minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}
	r, err := client.InspectContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetResult())
}

// [END dlp_inspect_string]

// [START dlp_inspect_file]
func inspectFile(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, bytesType dlppb.ByteContentItem_BytesType, fileName string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     i,
			MinLikelihood: minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_ByteItem{
				ByteItem: &dlppb.ByteContentItem{
					Type: bytesType,
					Data: b,
				},
			},
		},
	}
	r, err := client.InspectContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetResult())
}

// [END dlp_inspect_file]

// [START dlp_inspect_gcs]
func inspectGCSFile(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, pubSubTopic, pubSubSub, bucketName, fileName string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}

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
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
				StorageConfig: &dlppb.StorageConfig{
					Type: &dlppb.StorageConfig_CloudStorageOptions{
						CloudStorageOptions: &dlppb.CloudStorageOptions{
							FileSet: &dlppb.CloudStorageOptions_FileSet{
								Url: "gs://" + bucketName + "/" + fileName,
							},
						},
					},
				},
				InspectConfig: &dlppb.InspectConfig{
					InfoTypes:     i,
					MinLikelihood: minLikelihood,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: maxFindings,
					},
					IncludeQuote: includeQuote,
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
			for _, s := range jr.GetInspectDetails().GetResult().GetInfoTypeStats() {
				fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
			}
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_inspect_gcs]

// [START dlp_inspect_datastore]
func inspectDatastore(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, pubSubTopic, pubSubSub, dataProject, namespaceID, kind string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}

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
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
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
				InspectConfig: &dlppb.InspectConfig{
					InfoTypes:     i,
					MinLikelihood: minLikelihood,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: maxFindings,
					},
					IncludeQuote: includeQuote,
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
			for _, s := range jr.GetInspectDetails().GetResult().GetInfoTypeStats() {
				fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
			}
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_inspect_datastore]

// [START dlp_inspect_bigquery]
func inspectBigquery(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, pubSubTopic, pubSubSub, dataProject, datasetID, tableID string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}

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
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: &dlppb.InspectJobConfig{
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
				InspectConfig: &dlppb.InspectConfig{
					InfoTypes:     i,
					MinLikelihood: minLikelihood,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: maxFindings,
					},
					IncludeQuote: includeQuote,
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
			for _, s := range jr.GetInspectDetails().GetResult().GetInfoTypeStats() {
				fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
			}
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Error receiving from PubSub: %v\n", err)
	}
}

// [END dlp_inspect_bigquery]
