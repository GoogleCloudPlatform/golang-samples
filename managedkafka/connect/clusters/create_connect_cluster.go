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

package clusters

// [START managedkafka_create_connect_cluster]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func createConnectCluster(w io.Writer, projectID, region, clusterID, kafkaCluster string, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-connect-cluster"
	// kafkaCluster := "projects/my-project-id/locations/us-central1/clusters/my-kafka-cluster"
	ctx := context.Background()
	client, err := managedkafka.NewManagedKafkaConnectClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewManagedKafkaConnectClient got err: %w", err)
	}
	defer client.Close()

	locationPath := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	clusterPath := fmt.Sprintf("%s/connectClusters/%s", locationPath, clusterID)

	// Capacity configuration with 12 vCPU and 12 GiB RAM
	capacityConfig := &managedkafkapb.CapacityConfig{
		VcpuCount:   12,
		MemoryBytes: 12884901888, // 12 GiB in bytes
	}

	connectCluster := &managedkafkapb.ConnectCluster{
		Name:           clusterPath,
		KafkaCluster:   kafkaCluster,
		CapacityConfig: capacityConfig,
	}

	req := &managedkafkapb.CreateConnectClusterRequest{
		Parent:           locationPath,
		ConnectClusterId: clusterID,
		ConnectCluster:   connectCluster,
	}
	op, err := client.CreateConnectCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateConnectCluster got err: %w", err)
	}
	// The duration of this operation can vary considerably, typically taking 5-15 minutes.
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("op.Wait got err: %w", err)
	}
	fmt.Fprintf(w, "Created connect cluster: %s\n", resp.Name)
	return nil
}

// [END managedkafka_create_connect_cluster]
