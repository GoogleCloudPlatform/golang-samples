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

// [START pubsub_create_topic_with_confluent_cloud_ingestion]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

func createTopicWithConfluentCloudIngestion(w io.Writer, projectID, topicID, bootstrapServer, clusterID, confluentTopic, poolID, gcpSA string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"

	// // Confluent Cloud ingestion settings.
	// bootstrapServer := "bootstrap-server"
	// clusterID := "cluster-id"
	// confluentTopic := "confluent-topic"
	// poolID := "identity-pool-id"
	// gcpSA := "gcp-service-account"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	topicpb := &pubsubpb.Topic{
		Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
		IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{
			Source: &pubsubpb.IngestionDataSourceSettings_ConfluentCloud_{
				ConfluentCloud: &pubsubpb.IngestionDataSourceSettings_ConfluentCloud{
					BootstrapServer:   bootstrapServer,
					ClusterId:         clusterID,
					Topic:             confluentTopic,
					IdentityPoolId:    poolID,
					GcpServiceAccount: gcpSA,
				},
			},
		},
	}
	topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
	if err != nil {
		return fmt.Errorf("CreateTopic: %w", err)
	}
	fmt.Fprintf(w, "Created topic with Confluent Cloud ingestion: %v\n", topic)
	return nil
}

// [END pubsub_create_topic_with_confluent_cloud_ingestion]
