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

// [START pubsub_old_version_update_topic_type]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func updateTopicType(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
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

	updateCfg := pubsub.TopicConfigToUpdate{
		// If wanting to clear ingestion settings, set this to zero value: &pubsub.IngestionDataSourceSettings{}
		IngestionDataSourceSettings: &pubsub.IngestionDataSourceSettings{
			Source: &pubsub.IngestionDataSourceAWSKinesis{
				StreamARN:         streamARN,
				ConsumerARN:       consumerARN,
				AWSRoleARN:        awsRoleARN,
				GCPServiceAccount: gcpServiceAccount,
			},
		},
	}
	topicCfg, err := client.Topic(topicID).Update(ctx, updateCfg)
	if err != nil {
		return fmt.Errorf("topic.Update: %w", err)
	}
	fmt.Fprintf(w, "Topic updated with kinesis source: %v\n", topicCfg)
	return nil
}

// [END pubsub_old_version_update_topic_type]
