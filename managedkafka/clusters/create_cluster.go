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

package clusters

// [START managedkafka_create_cluster]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func createCluster(w io.Writer, projectID, region, clusterID, subnet string, cpu, memoryBytes int64, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-cluster"
	// subnet := "projects/my-project-id/regions/us-central1/subnetworks/default"
	// cpu := 3
	// memoryBytes := 3221225472
	ctx := context.Background()
	client, err := managedkafka.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewClient got err: %w", err)
	}
	defer client.Close()

	locationPath := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	clusterPath := fmt.Sprintf("%s/clusters/%s", locationPath, clusterID)

	// Memory must be between 1 GiB and 8 GiB per CPU.
	capacityConfig := &managedkafkapb.CapacityConfig{
		VcpuCount:   cpu,
		MemoryBytes: memoryBytes,
	}
	var networkConfig []*managedkafkapb.NetworkConfig
	networkConfig = append(networkConfig, &managedkafkapb.NetworkConfig{
		Subnet: subnet,
	})
	platformConfig := &managedkafkapb.Cluster_GcpConfig{
		GcpConfig: &managedkafkapb.GcpConfig{
			AccessConfig: &managedkafkapb.AccessConfig{
				NetworkConfigs: networkConfig,
			},
		},
	}
	rebalanceConfig := &managedkafkapb.RebalanceConfig{
		Mode: managedkafkapb.RebalanceConfig_AUTO_REBALANCE_ON_SCALE_UP,
	}
	cluster := &managedkafkapb.Cluster{
		Name:            clusterPath,
		CapacityConfig:  capacityConfig,
		PlatformConfig:  platformConfig,
		RebalanceConfig: rebalanceConfig,
	}

	req := &managedkafkapb.CreateClusterRequest{
		Parent:    locationPath,
		ClusterId: clusterID,
		Cluster:   cluster,
	}
	op, err := client.CreateCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateCluster got err: %w", err)
	}
	// The duration of this operation can vary considerably, typically taking 10-40 minutes.
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("op.Wait got err: %w", err)
	}
	fmt.Fprintf(w, "Created cluster: %s\n", resp.Name)
	return nil
}

// [END managedkafka_create_cluster]
