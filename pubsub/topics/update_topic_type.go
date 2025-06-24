// Copyright 2024 Google LLC
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

package topics

// [START pubsub_update_topic_type]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func updateTopicType(w io.Writer, projectID, topic string) error {
	// projectID := "my-project-id"
	// topic := "projects/my-project-id/topics/my-topic"
	streamARN := "stream-arn"
	consumerARN := "consumer-arn"
	awsRoleARN := "aws-role-arn"
	gcpServiceAccount := "gcp-service-account"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	pbTopic := &pubsubpb.Topic{
		Name: topic,
		IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{
			Source: &pubsubpb.IngestionDataSourceSettings_AwsKinesis_{
				AwsKinesis: &pubsubpb.IngestionDataSourceSettings_AwsKinesis{
					StreamArn:         streamARN,
					ConsumerArn:       consumerARN,
					AwsRoleArn:        awsRoleARN,
					GcpServiceAccount: gcpServiceAccount,
				},
			},
		},
	}
	updateReq := &pubsubpb.UpdateTopicRequest{
		Topic: pbTopic,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"ingestion_data_source_settings"},
		},
	}
	topicCfg, err := client.TopicAdminClient.UpdateTopic(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("topic.Update: %w", err)
	}
	fmt.Fprintf(w, "Topic updated with kinesis source: %v\n", topicCfg)
	return nil
}

// [END pubsub_update_topic_type]
