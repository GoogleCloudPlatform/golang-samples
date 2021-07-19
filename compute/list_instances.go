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

// [START compute_instances_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// listInstances gets a list of instances created in given project in given zone.
func listInstances(w io.Writer, projectID string, zone string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}

	defer client.Close()

	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
	}

	resp, err := client.List(ctx, req)
	if err != nil {
		return fmt.Errorf("List instances request: %v", err)
	}

	fmt.Fprintf(w, "Instances found in zone %s:\n", zone)

	for _, instance := range resp.Items {
		fmt.Fprintf(w, "- %s %s\n", *instance.Name, *instance.MachineType)
	}

	return nil
}

// [END compute_instances_list]
