// Copyright 2019 Google LLC
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

// Package dataproc shows how you can use the Dataproc Client library to manage
// Dataproc clusters. In this example, we'll show how to create a cluster.
package dataproc

// [START dataproc_create_cluster]
import (
	"context"
	"fmt"
	"io"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/dataproc/apiv1/dataprocpb"
	"google.golang.org/api/option"
)

func createCluster(w io.Writer, projectID, region, clusterName string) error {
	// projectID := "your-project-id"
	// region := "us-central1"
	// clusterName := "your-cluster"
	ctx := context.Background()

	// Create the cluster client.
	endpoint := region + "-dataproc.googleapis.com:443"
	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("dataproc.NewClusterControllerClient: %w", err)
	}
	defer clusterClient.Close()

	// Create the cluster config.
	req := &dataprocpb.CreateClusterRequest{
		ProjectId: projectID,
		Region:    region,
		Cluster: &dataprocpb.Cluster{
			ProjectId:   projectID,
			ClusterName: clusterName,
			Config: &dataprocpb.ClusterConfig{
				MasterConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   1,
					MachineTypeUri: "n1-standard-2",
				},
				WorkerConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   2,
					MachineTypeUri: "n1-standard-2",
				},
			},
		},
	}

	// Create the cluster.
	op, err := clusterClient.CreateCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCluster: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("CreateCluster.Wait: %w", err)
	}

	// Output a success message.
	fmt.Fprintf(w, "Cluster created successfully: %s", resp.ClusterName)
	return nil
}

// [END dataproc_create_cluster]
