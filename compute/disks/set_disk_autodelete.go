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

// [START compute_disk_autodelete_change]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// setDiskAutodelete sets the autodelete flag of a disk to given value.
func setDiskAutoDelete(
	w io.Writer,
	projectID, zone, instanceName, diskName string, autoDelete bool,
) error {
	// projectID := "your_project_id"
	// zone := "us-west3-b"
	// instanceName := "your_instance_name"
	// diskName := "your_disk_name"
	// autoDelete := true

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	getInstanceReq := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, getInstanceReq)
	if err != nil {
		return fmt.Errorf("unable to get instance: %w", err)
	}

	diskExists := false

	for _, disk := range instance.GetDisks() {
		if disk.GetDeviceName() == diskName {
			diskExists = true
			break
		}
	}

	if !diskExists {
		return fmt.Errorf(
			"instance %s doesn't have a disk named %s attached",
			instanceName,
			diskName,
		)
	}

	req := &computepb.SetDiskAutoDeleteInstanceRequest{
		Project:    projectID,
		Zone:       zone,
		Instance:   instanceName,
		DeviceName: diskName,
		AutoDelete: autoDelete,
	}

	op, err := instancesClient.SetDiskAutoDelete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to set disk autodelete field: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "disk autoDelete field updated.\n")

	return nil
}

// [END compute_disk_autodelete_change]
