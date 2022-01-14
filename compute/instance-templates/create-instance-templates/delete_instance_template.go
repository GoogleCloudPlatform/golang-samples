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

// [START compute_template_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// deleteInstanceTemplate deletes an instance template.
func deleteInstanceTemplate(w io.Writer, projectID, templateName string) error {
	// projectID := "your_project_id"
	// templateName := "your_template_name"

	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstanceTemplatesRESTClient: %v", err)
	}
	defer instanceTemplatesClient.Close()

	req := &computepb.DeleteInstanceTemplateRequest{
		Project:          projectID,
		InstanceTemplate: templateName,
	}

	op, err := instanceTemplatesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance template: %v", err)
	}

	globalOperationsClient, err := compute.NewGlobalOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGlobalOperationsRESTClient: %v", err)
	}
	defer globalOperationsClient.Close()

	for {
		waitReq := &computepb.WaitGlobalOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   projectID,
		}
		zoneOp, err := globalOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}

		if *zoneOp.Status.Enum() == computepb.Operation_DONE {
			fmt.Fprintf(w, "Instance template deleted\n")
			break
		}
	}

	return nil
}

// [END compute_template_delete]
