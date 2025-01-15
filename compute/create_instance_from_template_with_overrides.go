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

// [START compute_instances_create_from_template_with_overrides]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"

	"google.golang.org/protobuf/proto"
)

// createInstanceFromTemplate creates a Compute Engine VM instance from an instance template, but overrides the disk and machine type options in the template.
func createInstanceFromTemplateWithOverrides(w io.Writer, projectID, zone, instanceName, instanceTemplateName, machineType, newDiskSourceImage string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// instanceTemplateName := "your_instance_template_name"
	// machineType := "n1-standard-2"
	// newDiskSourceImage := "projects/debian-cloud/global/images/family/debian-12"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	intanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstanceTemplatesRESTClient: %w", err)
	}
	defer intanceTemplatesClient.Close()

	// Retrieve an instance template by name.
	reqGetTemplate := &computepb.GetInstanceTemplateRequest{
		Project:          projectID,
		InstanceTemplate: instanceTemplateName,
	}

	instanceTemplate, err := intanceTemplatesClient.Get(ctx, reqGetTemplate)
	if err != nil {
		return fmt.Errorf("unable to get intance template: %w", err)
	}

	for _, disk := range instanceTemplate.Properties.Disks {
		diskType := disk.InitializeParams.GetDiskType()
		if diskType != "" {
			disk.InitializeParams.DiskType = proto.String(fmt.Sprintf(`zones/%s/diskTypes/%s`, zone, diskType))
		}
	}

	reqInsertInstance := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name:        proto.String(instanceName),
			MachineType: proto.String(fmt.Sprintf(`zones/%s/machineTypes/%s`, zone, machineType)),
			Disks: append(
				// If you override a repeated field, all repeated values
				// for that property are replaced with the
				// corresponding values provided in the request.
				// When adding a new disk to existing disks,
				// insert all existing disks as well.
				instanceTemplate.Properties.Disks,
				&computepb.AttachedDisk{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: &newDiskSourceImage,
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(false),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			),
		},
		SourceInstanceTemplate: instanceTemplate.SelfLink,
	}

	op, err := instancesClient.Insert(ctx, reqInsertInstance)
	if err != nil {
		return fmt.Errorf("unable to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Instance created\n")

	return nil
}

// [END compute_instances_create_from_template_with_overrides]
