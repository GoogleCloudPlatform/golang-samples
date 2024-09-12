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

// [START compute_reservation_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// Get list of reservations for given project in particular zone
func listReservations(w io.Writer, projectID, zone string) error {
	// projectID := "your_project_id"
	// zone := "us-west3-a"

	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer reservationsClient.Close()

	req := &computepb.ListReservationsRequest{
		Project: projectID,
		Zone:    zone,
	}

	it := reservationsClient.List(ctx, req)
	fmt.Fprintf(w, "Instances found in zone %s:\n", zone)
	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s %d\n", instance.GetName(), instance.GetSpecificReservation().GetCount())
	}

	return nil
}

// [END compute_reservation_list]
