// Copyright 2025 Google LLC
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

// [START pubsub_old_version_create_topic_with_aws_msk_ingestion]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func createTopicWithAWSMSKIngestion(w io.Writer, projectID, topicID, clusterARN, mskTopic, awsRoleARN, gcpSA string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"

	// // AWS MSK ingestion settings.
	// clusterARN := "cluster-arn"
	// mskTopic := "msk-topic"
	// awsRoleARN := "aws-role-arn"
	// gcpSA := "gcp-service-account"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	cfg := &pubsub.TopicConfig{
		IngestionDataSourceSettings: &pubsub.IngestionDataSourceSettings{
			Source: &pubsub.IngestionDataSourceAmazonMSK{
				ClusterARN:        clusterARN,
				Topic:             mskTopic,
				AWSRoleARN:        awsRoleARN,
				GCPServiceAccount: gcpSA,
			},
		},
	}
	t, err := client.CreateTopicWithConfig(ctx, topicID, cfg)
	if err != nil {
		return fmt.Errorf("CreateTopic: %w", err)
	}
	fmt.Fprintf(w, "Created topic with AWS MSK ingestion settings: %v\n", t)
	return nil
}

// [END pubsub_create_topic_with_aws_msk_ingestion]
