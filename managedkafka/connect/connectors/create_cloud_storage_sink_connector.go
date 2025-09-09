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

package connectors

// [START managedkafka_create_cloud_storage_sink_connector]
import (
	"context"
	"fmt"
	"io"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"
)

// createCloudStorageSinkConnector creates a Cloud Storage Sink connector.
func createCloudStorageSinkConnector(w io.Writer, projectID, region, connectClusterID, connectorID, topics, gcsBucketName, tasksMax, formatOutputType, valueConverter, valueConverterSchemasEnable, keyConverter string, opts ...option.ClientOption) error {
	// TODO(developer): Update with your config values. Here is a sample configuration:
	// projectID := "my-project-id"
	// region := "us-central1"
	// connectClusterID := "my-connect-cluster"
	// connectorID := "GCS_SINK_CONNECTOR_ID"
	// topics := "GMK_TOPIC_ID"
	// gcsBucketName := "GCS_BUCKET_NAME"
	// tasksMax := "3"
	// formatOutputType := "json"
	// valueConverter := "org.apache.kafka.connect.json.JsonConverter"
	// valueConverterSchemasEnable := "false"
	// keyConverter := "org.apache.kafka.connect.storage.StringConverter"
	ctx := context.Background()
	client, err := managedkafka.NewManagedKafkaConnectClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewManagedKafkaConnectClient got err: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/connectClusters/%s", projectID, region, connectClusterID)

	config := map[string]string{
		"connector.class":                "io.aiven.kafka.connect.gcs.GcsSinkConnector",
		"tasks.max":                      tasksMax,
		"topics":                         topics,
		"gcs.bucket.name":                gcsBucketName,
		"gcs.credentials.default":        "true",
		"format.output.type":             formatOutputType,
		"name":                           connectorID,
		"value.converter":                valueConverter,
		"value.converter.schemas.enable": valueConverterSchemasEnable,
		"key.converter":                  keyConverter,
	}

	connector := &managedkafkapb.Connector{
		Name:    fmt.Sprintf("%s/connectors/%s", parent, connectorID),
		Configs: config,
	}

	req := &managedkafkapb.CreateConnectorRequest{
		Parent:      parent,
		ConnectorId: connectorID,
		Connector:   connector,
	}

	resp, err := client.CreateConnector(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateConnector got err: %w", err)
	}
	fmt.Fprintf(w, "Created Cloud Storage sink connector: %s\n", resp.Name)
	return nil
}

// [END managedkafka_create_cloud_storage_sink_connector]
