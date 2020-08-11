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

// [START cloud_game_servers_cluster_delete]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
)

// deleteCluster unregisters a game server cluster.
func deleteCluster(w io.Writer, projectID, location, realmID, clusterID string) error {
	// projectID := "my-project"
	// location := "global"
	// realmID := "myrealm"
	// clusterID := "mycluster"
	ctx := context.Background()
	client, err := gaming.NewGameServerClustersClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerClustersClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.DeleteGameServerClusterRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/realms/%s/gameServerClusters/%s", projectID, location, realmID, clusterID),
	}

	op, err := client.DeleteGameServerCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteGameServerCluster: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Cluster deleted.")
	return nil
}

// [END cloud_game_servers_cluster_delete]
