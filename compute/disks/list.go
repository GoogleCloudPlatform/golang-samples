//  Copyright 2024 Google LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package snippets

// [START compute_disk_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// listDisks lists disks in a project zone.
func listDisks(w io.Writer, projectID, zone, filter string) error {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// filter := ""

	// Formatting for filters:
	// https://cloud.google.com/python/docs/reference/compute/latest/google.cloud.compute_v1.types.ListDisksRequest

	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.ListDisksRequest{
		Project: projectID,
		Zone:    zone,
		Filter:  &filter,
	}

	it := client.List(ctx, req)
	fmt.Fprintf(w, "Disks in zone %s:\n", zone)
	for {
		disk, err := it.Next()
		if err == context.Canceled || err == context.DeadlineExceeded {
			return err
		}
		if err != nil {
			break
		}
		fmt.Fprintf(w, "- %s\n", *disk.Name)
	}

	if err != nil {
		return fmt.Errorf("ListDisks: %w", err)
	}
	return nil
}

// [END compute_disk_list]
