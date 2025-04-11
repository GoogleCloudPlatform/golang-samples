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

// [START dlp_inspect_gcs]
import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/pubsub"
)

// inspectGCSFile searches for the given info types in the given file.
func inspectGCSFile(w io.Writer, projectID string, infoTypeNames []string, customDictionaries []string, customRegexes []string, pubSubTopic, pubSubSub, bucketName, fileName string) error {
	// projectID := "my-project-id"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	// customDictionaries := []string{...}
	// customRegexes := []string{...}
	// pubSubTopic := "dlp-risk-sample-topic"
	// pubSubSub := "dlp-risk-sample-sub"
	// bucketName := "my-bucket"
	// fileName := "my-file.txt"

	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
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
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer pubsubClient.Close()

	// Create a PubSub subscription we can use to listen for messages.
	// Create the Topic if it doesn't exist.
	t := pubsubClient.Topic(pubSubTopic)
	if exists, err := t.Exists(ctx); err != nil {
		return fmt.Errorf("t.Exists: %w", err)
	} else if !exists {
		if t, err = pubsubClient.CreateTopic(ctx, pubSubTopic); err != nil {
			return fmt.Errorf("CreateTopic: %w", err)
		}
	}

	// Create the Subscription if it doesn't exist.
	s := pubsubClient.Subscription(pubSubSub)
	if exists, err := s.Exists(ctx); err != nil {
		return fmt.Errorf("s.Exists: %w", err)
	} else if !exists {
		if s, err = pubsubClient.CreateSubscription(ctx, pubSubSub, pubsub.SubscriptionConfig{Topic: t}); err != nil {
			return fmt.Errorf("CreateSubscription: %w", err)
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
		return fmt.Errorf("CreateDlpJob: %w", err)
	}
	fmt.Fprintf(w, "Created job: %v\n", j.GetName())

	// Wait for the inspect job to finish by waiting for a PubSub message.
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

		// Stop listening for more messages.
		defer cancel()

		resp, err := client.GetDlpJob(ctx, &dlppb.GetDlpJobRequest{
			Name: j.GetName(),
		})
		if err != nil {
			fmt.Fprintf(w, "Cloud not get job: %v", err)
			return
		}
		r := resp.GetInspectDetails().GetResult().GetInfoTypeStats()
		if len(r) == 0 {
			fmt.Fprintf(w, "No results")
		}
		for _, s := range r {
			fmt.Fprintf(w, "  Found %v instances of infoType %v\n", s.GetCount(), s.GetInfoType().GetName())
		}
	})
	if err != nil {
		return fmt.Errorf("Receive: %w", err)
	}
	return nil
}

// [END dlp_inspect_gcs]
