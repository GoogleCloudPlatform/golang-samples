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

// [START compute_template_list]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// listInstanceTemplates prints a list of InstanceTemplate objects available in a project.
func listInstanceTemplates(w io.Writer, projectID string) error {
	// projectID := "your_project_id"

	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstanceTemplatesRESTClient: %w", err)
	}
	defer instanceTemplatesClient.Close()

	req := &computepb.ListInstanceTemplatesRequest{
		Project: projectID,
	}

	it := instanceTemplatesClient.List(ctx, req)
	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s %s\n", instance.GetName(), instance.GetProperties().GetMachineType())
	}

	return nil
}

// [END compute_template_list]
