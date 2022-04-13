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

// [START compute_custom_machine_type_create_with_helper]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

func customMachineTypeHelper(zone, cpuSeries string, coreCount, memory int) (error, string) {
	const (
		N1        = "custom"
		N2        = "n2-custom"
		N2D       = "n2d-custom"
		E2        = "e2-custom"
		E2_MICRO  = "e2-custom-micro"
		E2_SMALL  = "e2-custom-small"
		E2_MEDIUM = "e2-custom-medium"
	)

	type TypeLimit struct {
		allowedCores     []int
		minMemPerCore    int
		maxMemPerCore    int
		allowExtraMemory bool
		extraMemoryLimit int
	}

	var (
		CPUSeries_E2_Limit        = TypeLimit{MakeRange(2, 33, 2), 512, 8192, false, 0}
		CPUSeries_E2_MICRO_Limit  = TypeLimit{[]int{}, 1024, 2048, false, 0}
		CPUSeries_E2_SMALL_Limit  = TypeLimit{[]int{}, 2048, 4096, false, 0}
		CPUSeries_E2_MEDIUM_Limit = TypeLimit{[]int{}, 4096, 8192, false, 0}
		CPUSeries_N2_Limit        = TypeLimit{append(MakeRange(2, 33, 2), MakeRange(36, 129, 4)...), 512, 8192, true, 624 << 10}
		CPUSeries_N2D_Limit       = TypeLimit{[]int{2, 4, 8, 16, 32, 48, 64, 80, 96}, 512, 8192, true, 768 << 10}
		CPUSeries_N1_Limit        = TypeLimit{append([]int{1}, MakeRange(2, 97, 2)...), 922, 6656, true, 624 << 10}
	)

	typeLimitsMap := map[string]TypeLimit{
		N1:        CPUSeries_N1_Limit,
		N2:        CPUSeries_N2_Limit,
		N2D:       CPUSeries_N2D_Limit,
		E2:        CPUSeries_E2_Limit,
		E2_MICRO:  CPUSeries_E2_MICRO_Limit,
		E2_SMALL:  CPUSeries_E2_SMALL_Limit,
		E2_MEDIUM: CPUSeries_E2_MEDIUM_Limit,
	}

	if !containsString([]string{E2, N1, N2, N2D}, cpuSeries) {
		return fmt.Errorf("incorrect cpu type: %v", cpuSeries), ""
	}

	typeLimit := typeLimitsMap[cpuSeries]

	// Check whether the requested parameters are allowed. Find more information about limitations of custom machine
	// types at: https://cloud.google.com/compute/docs/general-purpose-machines#custom_machine_types

	// Check the number of cores
	if len(typeLimit.allowedCores) > 0 && !containsInt(typeLimit.allowedCores, coreCount) {
		return fmt.Errorf("invalid number of cores requested. Allowed number of cores for %v is: %v", cpuSeries, typeLimit.allowedCores), ""
	}

	// Memory must be a multiple of 256 MB
	if memory%256 != 0 {
		return fmt.Errorf("requested memory must be a multiple of 256 MB"), ""
	}

	// Check if the requested memory isn't too little
	if memory < coreCount*typeLimit.minMemPerCore {
		return fmt.Errorf("requested memory is too low. Minimal memory for %v is %v MB per core", cpuSeries, typeLimit.minMemPerCore), ""
	}

	// Check if the requested memory isn't too much
	if memory > coreCount*typeLimit.maxMemPerCore && !typeLimit.allowExtraMemory {
		return fmt.Errorf("requested memory is too large.. Maximum memory allowed for %v is %v MB per core", cpuSeries, typeLimit.maxMemPerCore), ""
	}
	if memory > typeLimit.extraMemoryLimit && typeLimit.allowExtraMemory {
		return fmt.Errorf("requested memory is too large.. Maximum memory allowed for %v is %v MB", cpuSeries, typeLimit.extraMemoryLimit), ""
	}

	// Return the custom machine type in form of a string acceptable by Compute Engine API.
	if containsString([]string{E2_SMALL, E2_MICRO, E2_MEDIUM}, cpuSeries) {
		return nil, fmt.Sprintf("zones/%v/machineTypes/%v-%v", zone, cpuSeries, memory)
	}

	if memory > coreCount*typeLimit.maxMemPerCore {
		return nil, fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v-ext", zone, cpuSeries, coreCount, memory)
	}

	return nil, fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v", zone, cpuSeries, coreCount, memory)
}

// createInstanceWithCustomMachineTypeWithHelper creates a new VM instance with a custom machine type.
func createInstanceWithCustomMachineTypeWithHelper(w io.Writer, projectID, zone, instanceName, cpuSeries string, coreCount, memory int) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// cpuSeries := "e2-custom-micro" // the type of CPU you want to use"
	// coreCount := 2 // number of CPU cores you want to use.
	// memory := 256 // the amount of memory for the VM instance, in megabytes.

	err, machineType := customMachineTypeHelper(zone, cpuSeries, coreCount, memory)
	if err != nil {
		return fmt.Errorf("unable to create custom machine type string: %v", err)
	}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String("projects/debian-cloud/global/images/family/debian-10"),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(machineType),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %v", err)
	}

	fmt.Fprintf(w, "Instance created\n")

	return nil
}

// [END compute_custom_machine_type_create_with_helper]
