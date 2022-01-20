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

// [START cloud_game_servers_config_create]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
)

// fleet is the spec portion of an agones Fleet.  It must be in JSON format.
// See https://agones.dev/site/docs/reference/fleet/ for more on fleets.
const fleet = `
{
   "replicas": 10,
   "scheduling": "Packed",
   "strategy": {
      "type": "RollingUpdate",
      "rollingUpdate": {
         "maxSurge": "25%",
         "maxUnavailable": "25%"
      }
   },
   "template": {
      "metadata": {
         "labels": {
            "gameName": "udp-server"
         }
      },
      "spec": {
         "ports": [
            {
               "name": "default",
               "portPolicy": "Dynamic",
               "containerPort": 2156,
               "protocol": "TCP"
            }
         ],
         "health": {
            "initialDelaySeconds": 30,
            "periodSeconds": 60
         },
         "sdkServer": {
            "logLevel": "Info",
            "grpcPort": 9357,
            "httpPort": 9358
         },
         "template": {
            "spec": {
               "containers": [
                  {
                     "name": "dedicated",
                     "image": "gcr.io/agones-images/udp-server:0.17",
                     "imagePullPolicy": "Always",
                     "resources": {
                        "requests": {
                           "memory": "200Mi",
                           "cpu": "500m"
                        },
                        "limits": {
                           "memory": "200Mi",
                           "cpu": "500m"
                        }
                     }
                  }
               ]
            }
         }
      }
   }
}
`

// createGameServerConfig creates a game server config.
func createGameServerConfig(w io.Writer, projectID, deploymentID, configID string) error {
	// projectID := "my-project"
	// deploymentID := "mydeployment"
	// configID := "mydeployment"
	ctx := context.Background()
	client, err := gaming.NewGameServerConfigsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerConfigsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.CreateGameServerConfigRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/global/gameServerDeployments/%s", projectID, deploymentID),
		ConfigId: configID,
		GameServerConfig: &gamingpb.GameServerConfig{
			FleetConfigs: []*gamingpb.FleetConfig{
				{
					Name:      "fleet-spec-1",
					FleetSpec: fleet,
				},
			},
			Description: "My Game Server Config",
		},
	}

	op, err := client.CreateGameServerConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateGameServerConfig: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Config created: %v", resp.Name)
	return nil
}

// [END cloud_game_servers_config_create]
