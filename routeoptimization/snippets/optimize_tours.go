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

// [START routeoptimization_optimize_tours]
import (
	"context"
	"fmt"

	routeoptimization "cloud.google.com/go/maps/routeoptimization/apiv1"
	"google.golang.org/genproto/googleapis/type/latlng"

	rpb "cloud.google.com/go/maps/routeoptimization/apiv1/routeoptimizationpb"
)

func optimizeTours(projectID string) (*rpb.OptimizeToursResponse, error) {
	ctx := context.Background()
	c, err := routeoptimization.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("routeoptimization client: %w", err)
	}
	defer c.Close()

	// See https://pkg.go.dev/cloud.google.com/go/maps/routeoptimization/apiv1/routeoptimizationpb#OptimizeToursRequest.
	req := &rpb.OptimizeToursRequest{
		Parent: "projects/" + projectID,
		Model: &rpb.ShipmentModel{
			Shipments: []*rpb.Shipment{
				&rpb.Shipment{
					Deliveries: []*rpb.Shipment_VisitRequest{
						{ArrivalLocation: &latlng.LatLng{Latitude: 48.880942, Longitude: 2.323866}},
					},
				},
			},
			Vehicles: []*rpb.Vehicle{
				{
					StartLocation: &latlng.LatLng{Latitude: 48.863102, Longitude: 2.341204},
					EndLocation:   &latlng.LatLng{Latitude: 48.86311, Longitude: 2.341205},
				},
			},
		},
	}
	return c.OptimizeTours(ctx, req)
}

// [END routeoptimization_optimize_tours]
