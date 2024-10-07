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

// [START compute_regional_template_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// deleteRegionalTemplate deletes an instance template in specific region.
func deleteRegionalTemplate(w io.Writer, projectID, templateName, region string) error {
	// projectID := "your_project_id"
	// templateName := "your_template_name"
	// region := "us-east1"

	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewRegionInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRegionInstanceTemplatesRESTClient: %w", err)
	}
	defer instanceTemplatesClient.Close()

	req := &computepb.DeleteRegionInstanceTemplateRequest{
		Project:          projectID,
		Region:           region,
		InstanceTemplate: templateName,
	}

	op, err := instanceTemplatesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance template: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Instance template deleted\n")

	return nil
}

// [END compute_regional_template_delete]
