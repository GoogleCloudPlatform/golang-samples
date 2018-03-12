package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func inspect(w io.Writer, client *dlp.Client, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, project, input string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	rcr := &dlppb.InspectContentRequest{
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
	r, err := client.InspectContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetResult())
}

func inspectFile(w io.Writer, client *dlp.Client, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, project string, bytesType dlppb.ByteContentItem_BytesType, fileName string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	rcr := &dlppb.InspectContentRequest{
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
	r, err := client.InspectContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetResult())
}

func inspectGCSFile(w io.Writer, client *dlp.Client, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, project, pubSubTopic, pubSubSub, bucketName, fileName string) {
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

	rcr := &dlppb.CreateDlpJobRequest{
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

func inspectDatastore(w io.Writer, client *dlp.Client, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, project, pubSubTopic, pubSubSub, dataProject, namespaceID, kind string) {
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

	rcr := &dlppb.CreateDlpJobRequest{
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

func inspectBigquery(w io.Writer, client *dlp.Client, minLikelihood dlppb.Likelihood, maxFindings int32, includeQuote bool, infoTypes []string, project, pubSubTopic, pubSubSub, dataProject, datasetID, tableID string) {
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

	rcr := &dlppb.CreateDlpJobRequest{
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
