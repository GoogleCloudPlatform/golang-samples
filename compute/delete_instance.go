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

// [START compute_instances_delete]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// deleteInstance sends a delete request to GCP and waits for it to complete.
func deleteInstance(w io.Writer, projectID string, zone string, instanceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %v", err)
	}

	zoneOperationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewZoneOperationsRESTClient: %v", err)
	}
	defer zoneOperationsClient.Close()

	waitReq := &computepb.WaitZoneOperationRequest{
		Operation: op.GetName(),
		Project:   projectID,
		Zone:      zone,
	}

	op, err = zoneOperationsClient.Wait(ctx, waitReq)
	if err != nil {
		return fmt.Errorf("unable to wait for the operation: %v", err)
	}

	if op.GetStatus() == computepb.Operation_DONE {
		fmt.Fprintf(w, "Instance deleted\n")
	} else {
		return fmt.Errorf("delete instance operation has status %s", op.GetStatus())
	}

	return nil
}

// [END compute_instances_delete]
