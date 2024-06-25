// Copyright 2024 Google LLC
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

// [START compute_route_create]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createRoute creates a route with given name inside given project.
func createRoute(w io.Writer, projectID, name string) error {
	// projectID := "your_project_id"
	// name := "testname"

	ctx := context.Background()
	client, err := compute.NewRoutesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRoutesRESTClient: %w", err)
	}
	defer client.Close()

	route := &computepb.Route{
		Name:           proto.String(name),
		Network:        proto.String("global/networks/default"),
		DestRange:      proto.String("0.0.0.0/0"),
		NextHopGateway: proto.String("global/gateways/default-internet-gateway"),
	}

	req := &computepb.InsertRouteRequest{
		Project:       projectID,
		RouteResource: route,
	}
	op, err := client.Insert(ctx, req)

	if err != nil {
		return fmt.Errorf("unable to insert a route: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Route created\n")

	return nil
}

// [END compute_route_create]
