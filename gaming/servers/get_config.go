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

// [START cloud_game_servers_config_get]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1beta"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1beta"
)

// getGameServerConfig retrieves info on a game server config.
func getGameServerConfig(w io.Writer, projectID, deploymentID, configID string) error {
	// projectID := "my-project"
	// deploymentID := "mydeployment"
	// configID := "myconfig"
	ctx := context.Background()
	client, err := gaming.NewGameServerConfigsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerConfigsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.GetGameServerConfigRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/gameServerDeployments/%s/configs/%s", projectID, deploymentID, configID),
	}

	resp, err := client.GetGameServerConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("GetGameServerConfig: %v", err)
	}

	fmt.Fprintf(w, "Config retrieved: %v", resp.Name)
	return nil
}

// [END cloud_game_servers_config_get]
