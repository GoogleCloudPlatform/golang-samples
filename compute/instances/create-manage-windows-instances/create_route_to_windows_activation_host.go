// Copyright 2022 Google LLC
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

package snippets

// [START compute_create_route_windows_activation]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createRouteToWindowsActivationHost creates a new route to
// kms.windows.googlecloud.com (35.190.247.13) for Windows activation.
func createRouteToWindowsActivationHost(
	w io.Writer,
	projectID, routeName, networkName string,
) error {
	// projectID := "your_project_id"
	// routeName := "your_route_name"
	// networkName := "global/networks/default"

	ctx := context.Background()
	routesClient, err := compute.NewRoutesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRoutesRESTClient: %w", err)
	}
	defer routesClient.Close()

	// If you have Windows instances without external IP addresses,
	// you must also enable Private Google Access so that instances
	// with only internal IP addresses can send traffic to the external
	// IP address for kms.windows.googlecloud.com.
	// More infromation: https://cloud.google.com/vpc/docs/configure-private-google-access#enabling
	req := &computepb.InsertRouteRequest{
		Project: projectID,
		RouteResource: &computepb.Route{
			Name:      proto.String(routeName),
			DestRange: proto.String("35.190.247.13/32"),
			Network:   proto.String(networkName),
			NextHopGateway: proto.String(
				fmt.Sprintf("projects/%s/global/gateways/default-internet-gateway", projectID),
			),
		},
	}

	op, err := routesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create route: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Route created\n")

	return nil
}

// [END compute_create_route_windows_activation]
