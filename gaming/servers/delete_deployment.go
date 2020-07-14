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

// [START cloud_game_servers_deployment_delete]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1beta"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1beta"
)

// deleteGameServerDeployment deletes a game server deployment.
func deleteGameServerDeployment(w io.Writer, projectID, deploymentID string) error {
	// projectID := "my-project"
	// deploymentID := "mydeployment"
	ctx := context.Background()
	client, err := gaming.NewGameServerDeploymentsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerDeploymentsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.DeleteGameServerDeploymentRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/gameServerDeployments/%s", projectID, deploymentID),
	}

	op, err := client.DeleteGameServerDeployment(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteGameServerDeployment: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Deployment deleted.")
	return nil
}

// [END cloud_game_servers_deployment_delete]
