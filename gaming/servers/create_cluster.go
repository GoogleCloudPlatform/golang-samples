// Copyright 2020 Google LLC
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

package servers

// [START cloud_game_servers_cluster_create]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
)

// createCluster registers a game server cluster.
func createCluster(w io.Writer, projectID, location, realmID, clusterID, gkeClusterName string) error {
	// projectID := "my-project"
	// location := "global"
	// realmID := "myrealm"
	// clusterID := "mycluster"
	// gkeClusterName := "projects/1234/locations/us-central1/clusters/gke-cluster"
	ctx := context.Background()
	client, err := gaming.NewGameServerClustersClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerClustersClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.CreateGameServerClusterRequest{
		Parent:              fmt.Sprintf("projects/%s/locations/%s/realms/%s", projectID, location, realmID),
		GameServerClusterId: clusterID,
		GameServerCluster: &gamingpb.GameServerCluster{
			ConnectionInfo: &gamingpb.GameServerClusterConnectionInfo{
				Namespace: "default",
				ClusterReference: &gamingpb.GameServerClusterConnectionInfo_GkeClusterReference{
					GkeClusterReference: &gamingpb.GkeClusterReference{
						Cluster: gkeClusterName,
					},
				},
			},
			Description: "My Game Server Cluster",
		},
	}

	op, err := client.CreateGameServerCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateGameServerCluster: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Cluster created: %v", resp.Name)
	return nil
}

// [END cloud_game_servers_cluster_create]
