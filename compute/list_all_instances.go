// Copyright 2021 Google LLC
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

// [START compute_instances_list_all]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// listAllInstances prints all instances present in a project, grouped by their zone.
func listAllInstances(w io.Writer, projectID string) error {
	// zone := "europe-central2-b"
	ctx := context.Background()
	projectsClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer projectsClient.Close()

	req := &computepb.AggregatedListInstancesRequest{
		Project: projectID,
	}

	resp, err := projectsClient.AggregatedList(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to call AggregatedList request: %v", err)
	}

	fmt.Fprintf(w, "Instances found:\n")

	for zone := range resp.Items {
		instances := resp.Items[zone].Instances
		if len(instances) > 0 {
			fmt.Fprintf(w, "%s\n", zone)
			for _, instance := range instances {
				fmt.Fprintf(w, "- %s %s\n", *instance.Name, *instance.MachineType)
			}
		}
	}

	return nil
}

// [END compute_instances_list_all]
