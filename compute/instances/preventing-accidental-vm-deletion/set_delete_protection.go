// Copyright 2022 Google LLC
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

// [START compute_delete_protection_set]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// setDeleteProtection updates the delete protection setting of given instance.
func setDeleteProtection(w io.Writer, projectID, zone, instanceName string, deleteProtection bool) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// deleteProtection := true

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.SetDeletionProtectionInstanceRequest{
		Project:            projectID,
		Zone:               zone,
		Resource:           instanceName,
		DeletionProtection: proto.Bool(deleteProtection),
	}

	op, err := instancesClient.SetDeletionProtection(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to set deletion protection: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %v", err)
	}

	fmt.Fprintf(w, "Instance updated\n")

	return nil
}

// [END compute_delete_protection_set]
