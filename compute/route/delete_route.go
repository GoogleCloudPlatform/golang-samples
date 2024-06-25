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

// [START compute_route_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// deleteRoute deletes a route by name in given project.
func deleteRoute(w io.Writer, projectID, name string) error {
	// projectID := "your_project_id"
	// name := "testname"

	ctx := context.Background()
	client, err := compute.NewRoutesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRoutesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.DeleteRouteRequest{
		Project: projectID,
		Route:   name,
	}
	op, err := client.Delete(ctx, req)

	if err != nil {
		return fmt.Errorf("unable to delete a route: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Route deleted\n")

	return nil
}

// [END compute_route_delete]
