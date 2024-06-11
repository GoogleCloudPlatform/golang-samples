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

// [START managedkafka_update_cluster]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func updateCluster(w io.Writer, projectID, region, clusterID string, memory int64, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-cluster"
	// memoryBytes := 4221225472
	ctx := context.Background()
	client, err := managedkafka.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewClient got err: %w", err)
	}
	defer client.Close()

	clusterPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, clusterID)
	capacityConfig := &managedkafkapb.CapacityConfig{
		MemoryBytes: memory,
	}
	cluster := &managedkafkapb.Cluster{
		Name:           clusterPath,
		CapacityConfig: capacityConfig,
	}
	paths := []string{"capacity_config.memory_bytes"}
	updateMask := &fieldmaskpb.FieldMask{
		Paths: paths,
	}

	req := &managedkafkapb.UpdateClusterRequest{
		UpdateMask: updateMask,
		Cluster:    cluster,
	}
	op, err := client.UpdateCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("client.UpdateCluster got err: %w", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("op.Wait got err: %w", err)
	}
	fmt.Fprintf(w, "Updated cluster: %#v\n", resp)
	return nil
}

// [END managedkafka_update_cluster]
