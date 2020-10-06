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

// [START cloud_game_servers_rollout_get]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
)

// getGameServerDeploymentRollout retrieves info on a game server deployment's rollout.
func getGameServerDeploymentRollout(w io.Writer, projectID, deploymentID string) error {
	// projectID := "my-project"
	// deploymentID := "mydeployment"
	ctx := context.Background()
	client, err := gaming.NewGameServerDeploymentsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerDeploymentsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.GetGameServerDeploymentRolloutRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/gameServerDeployments/%s", projectID, deploymentID),
	}

	resp, err := client.GetGameServerDeploymentRollout(ctx, req)
	if err != nil {
		return fmt.Errorf("GetGameServerDeploymentRollout: %v", err)
	}

	fmt.Fprintf(w, "Rollout default: %v\n", resp.DefaultGameServerConfig)
	for _, override := range resp.GameServerConfigOverrides {
		switch override.Selector.(type) {
		case *gamingpb.GameServerConfigOverride_RealmsSelector:
			fmt.Fprint(w, "Override these realms ", override.GetRealmsSelector().Realms)
		default:
			fmt.Fprint(w, "Unknown override selector")
		}

		switch override.Change.(type) {
		case *gamingpb.GameServerConfigOverride_ConfigVersion:
			fmt.Fprint(w, "with this config: ", override.GetConfigVersion())
		default:
			fmt.Fprint(w, "with an unknown override")
		}
		fmt.Fprintln(w)
	}
	return nil
}

// [END cloud_game_servers_rollout_get]
