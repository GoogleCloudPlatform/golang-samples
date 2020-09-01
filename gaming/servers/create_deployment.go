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

// [START cloud_game_servers_deployment_create]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
)

// createGameServerDeployment creates a game server deployment.
func createGameServerDeployment(w io.Writer, projectID, deploymentID string) error {
	// projectID := "my-project"
	// deploymentID := "mydeployment"
	ctx := context.Background()
	client, err := gaming.NewGameServerDeploymentsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerDeploymentsClient: %v", err)
	}
	defer client.Close()

	// Location is hard coded as global, as Game Server Deployments can only be
	// created in global.  This is done for all operations on Game Server
	// Deployments, as well as for its child resource types.
	req := &gamingpb.CreateGameServerDeploymentRequest{
		Parent:       fmt.Sprintf("projects/%s/locations/global", projectID),
		DeploymentId: deploymentID,
		GameServerDeployment: &gamingpb.GameServerDeployment{
			Description: "My Game Server Deployment",
		},
	}

	op, err := client.CreateGameServerDeployment(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateGameServerDeployment: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Deployment created: %v", resp.Name)
	return nil
}

// [END cloud_game_servers_deployment_create]
