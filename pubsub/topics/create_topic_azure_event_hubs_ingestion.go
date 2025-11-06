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

// [START pubsub_create_topic_with_azure_event_hubs_ingestion]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

func createTopicWithAzureEventHubsIngestion(w io.Writer, projectID, topicID, resourceGroup, namespace, eventHub, clientID, tenantID, subID, gcpSA string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"

	// // Azure Event Hubs ingestion settings.
	// resourceGroup := "resource-group"
	// namespace := "namespace"
	// eventHub := "event-hub"
	// clientID := "client-id"
	// tenantID := "tenant-id"
	// subID := "subscription-id"
	// gcpSA := "gcp-service-account"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	topicpb := &pubsubpb.Topic{
		IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{
			Source: &pubsubpb.IngestionDataSourceSettings_AzureEventHubs_{
				AzureEventHubs: &pubsubpb.IngestionDataSourceSettings_AzureEventHubs{
					ResourceGroup:     resourceGroup,
					Namespace:         namespace,
					EventHub:          eventHub,
					ClientId:          clientID,
					TenantId:          tenantID,
					SubscriptionId:    subID,
					GcpServiceAccount: gcpSA,
				},
			},
		},
	}
	topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
	if err != nil {
		return fmt.Errorf("CreateTopic: %w", err)
	}
	fmt.Fprintf(w, "Created topic with azure event hubs ingestion: %v\n", topic)
	return nil
}

// [END pubsub_create_topic_with_azure_event_hubs_ingestion]
