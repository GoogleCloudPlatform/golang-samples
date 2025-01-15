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

// [START compute_route_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// listRoutes prints a list of routes created in given project.
func listRoutes(w io.Writer, projectID string) error {
	// projectID := "your_project_id"

	ctx := context.Background()
	client, err := compute.NewRoutesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRoutesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.ListRoutesRequest{
		Project: projectID,
	}
	routes := client.List(ctx, req)

	for {
		instance, err := routes.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s\n", *instance.Name)
	}

	return nil
}

// [END compute_route_list]
