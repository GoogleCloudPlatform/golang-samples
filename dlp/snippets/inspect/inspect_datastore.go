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

package inspect

// [START dlp_inspect_datastore]
import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/pubsub"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// inspectDatastore searches for the given info types in the given dataset kind.
func inspectDatastore(w io.Writer, projectID string, infoTypeNames []string, customDictionaries []string, customRegexes []string, pubSubTopic, pubSubSub, dataProject, namespaceID, kind string) error {
	// projectID := "my-project-id"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	// customDictionaries := []string{...}
	// customRegexes := []string{...}
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// namespaceID := "namespace-id"
	// kind := "MyKind"

	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}

	// Convert the info type strings to a list of InfoTypes.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
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

	// Create a PubSub Client used to listen for when the inspect job finishes.
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	// Create the Topic if it doesn't exist.
	t := pubsubClient.Topic(pubSubTopic)
	if exists, err := t.Exists(ctx); err != nil {
		return fmt.Errorf("t.Exists: %v", err)
	} else if !exists {
		if t, err = pubsubClient.CreateTopic(ctx, pubSubTopic); err != nil {
			return fmt.Errorf("CreateTopic: %v", err)
		}
	}

	// Create the Subscription if it doesn't exist.
	s := pubsubClient.Subscription(pubSubSub)
	if exists, err := s.Exists(ctx); err != nil {
		return fmt.Errorf("s.Exists: %v", err)
	} else if !exists {
		if s, err = pubsubClient.CreateSubscription(ctx, pubSubSub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return fmt.Errorf("CreateSubscription: %v", err)
		}
	}

	// topic is the PubSub topic string where messages should be sent.
	topic := "projects/" + projectID + "/topics/" + pubSubTopic

	// Create a configured request.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
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
					InfoTypes:       infoTypes,
					CustomInfoTypes: customInfoTypes,
					MinLikelihood:   dlppb.Likelihood_POSSIBLE,
					Limits: &dlppb.InspectConfig_FindingLimits{
						MaxFindingsPerRequest: 10,
					},
					IncludeQuote: true,
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
	j, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateDlpJob: %v", err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the inspect job to finish by waiting for a PubSub message.
	// This only waits for 1 minute. For long jobs, consider using a truly
	// asynchronous execution model such as Cloud Functions.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	err = s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// If this is the wrong job, do not process the result.
		if msg.Attributes["DlpJobName"] != j.GetName() {
			msg.Nack()
			return
		}
		msg.Ack()

		// Stop listening for more messages.
		defer cancel()

		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			fmt.Fprintf(w, "Error getting completed job: %v\n", err)
			return
		}
		r := resp.GetInspectDetails().GetResult().GetInfoTypeStats()
		if len(r) == 0 {
			fmt.Fprintf(w, "No results")
			return
		}
		for _, s := range r {
			fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
		}
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

// [END dlp_inspect_datastore]
