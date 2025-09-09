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

// [START managedkafka_create_mirrormaker2_source_connector]
import (
	"context"
	"fmt"
	"io"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"
)

// createMirrorMaker2SourceConnector creates a MirrorMaker 2.0 Source connector.
func createMirrorMaker2SourceConnector(w io.Writer, projectID, region, connectClusterID, connectorID, sourceBootstrapServers, targetBootstrapServers, tasksMax, sourceClusterAlias, targetClusterAlias, topics, topicsExclude string, opts ...option.ClientOption) error {
	// TODO(developer): Update with your config values. Here is a sample configuration:
	// projectID := "my-project-id"
	// region := "us-central1"
	// connectClusterID := "my-connect-cluster"
	// connectorID := "mm2-source-to-target-connector-id"
	// sourceBootstrapServers := "source_cluster_dns"
	// targetBootstrapServers := "target_cluster_dns"
	// tasksMax := "3"
	// sourceClusterAlias := "source"
	// targetClusterAlias := "target"
	// topics := ".*"
	// topicsExclude := "mm2.*.internal,.*.replica,__.*"
	ctx := context.Background()
	client, err := managedkafka.NewManagedKafkaConnectClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewManagedKafkaConnectClient got err: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/connectClusters/%s", projectID, region, connectClusterID)

	config := map[string]string{
		"connector.class":      "org.apache.kafka.connect.mirror.MirrorSourceConnector",
		"name":                 connectorID,
		"tasks.max":            tasksMax,
		"source.cluster.alias": sourceClusterAlias,
		"target.cluster.alias": targetClusterAlias, // This is usually the primary cluster.
		// Replicate all topics from the source
		"topics": topics,
		// The value for bootstrap.servers is a hostname:port pair for the Kafka broker in
		// the source/target cluster.
		// For example: "kafka-broker:9092"
		"source.cluster.bootstrap.servers": sourceBootstrapServers,
		"target.cluster.bootstrap.servers": targetBootstrapServers,
		// You can define an exclusion policy for topics as follows:
		// To exclude internal MirrorMaker 2 topics, internal topics and replicated topics.
		// topicsExclude := "mm2.*.internal,.*.replica,__.*"
		"topics.exclude": topicsExclude,
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
	fmt.Fprintf(w, "Created MirrorMaker 2.0 Source connector: %s\n", resp.Name)
	return nil
}

// [END managedkafka_create_mirrormaker2_source_connector]
