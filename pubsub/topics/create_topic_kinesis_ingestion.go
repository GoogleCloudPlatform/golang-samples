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

// [START pubsub_create_topic_with_kinesis_ingestion]
import (
	"context"
	"fmt"
	"io"

	pubsub "cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

func createTopicWithKinesisIngestion(w io.Writer, projectID, topic string) error {
	// projectID := "my-project-id"
	// topicID := "projects/my-project-id/topics/my-topic"
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

	topicpb := &pubsubpb.Topic{
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
	topicpb, err = client.TopicAdminClient.CreateTopic(ctx, topicpb)
	if err != nil {
		return fmt.Errorf("failed to create topic with kinesis: %w", err)
	}
	fmt.Fprintf(w, "Kinesis topic created: %v\n", topicpb)
	return nil
}

// [END pubsub_create_topic_with_kinesis_ingestion]
