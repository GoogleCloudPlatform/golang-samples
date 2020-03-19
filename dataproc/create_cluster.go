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

// Package dataproc shows how you can use the Cloud Dataproc Client library to manage
// Cloud Dataproc clusters. In this example, we'll show how to create a cluster.
package dataproc

// [START dataproc_create_cluster]
import (
	"context"
	"fmt"
	"io"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1"
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
		return fmt.Errorf("dataproc.NewClusterControllerClient: %v", err)
	}

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
					MachineTypeUri: "n1-standard-1",
				},
				WorkerConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   2,
					MachineTypeUri: "n1-standard-1",
				},
			},
		},
	}

	// Create the cluster.
	op, err := clusterClient.CreateCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCluster: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("CreateCluster.Wait: %v", err)
	}

	// Output a success message.
	fmt.Fprintf(w, "Cluster created successfully: %s", resp.ClusterName)
	return nil
}

// [END dataproc_create_cluster]
