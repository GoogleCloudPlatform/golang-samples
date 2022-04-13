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

// [START compute_custom_machine_type_create_shared_with_helper]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

func customMachineTypeSharedCoreHelper(zone, cpuSeries string, memory int) (error, string) {
	const (
		n1       = "custom"
		n2       = "n2-custom"
		n2d      = "n2d-custom"
		e2       = "e2-custom"
		e2Micro  = "e2-custom-micro"
		e2Small  = "e2-custom-small"
		e2Medium = "e2-custom-medium"
	)

	type TypeLimit struct {
		allowedCores     []int
		minMemPerCore    int
		maxMemPerCore    int
		allowExtraMemory bool
		extraMemoryLimit int
	}

	var (
		CPUSeriesE2Limit       = TypeLimit{MakeRange(2, 33, 2), 512, 8192, false, 0}
		CPUSeriesE2MicroLimit  = TypeLimit{[]int{}, 1024, 2048, false, 0}
		CPUSeriesE2SmallLimit  = TypeLimit{[]int{}, 2048, 4096, false, 0}
		CPUSeriesE2MediumLimit = TypeLimit{[]int{}, 4096, 8192, false, 0}
		CPUSeriesN2Limit       = TypeLimit{append(MakeRange(2, 33, 2), MakeRange(36, 129, 4)...), 512, 8192, true, 624 << 10}
		CPUSeriesN2DLimit      = TypeLimit{[]int{2, 4, 8, 16, 32, 48, 64, 80, 96}, 512, 8192, true, 768 << 10}
		CPUSeriesN1Limit       = TypeLimit{append([]int{1}, MakeRange(2, 97, 2)...), 922, 6656, true, 624 << 10}
	)

	typeLimitsMap := map[string]TypeLimit{
		n1:       CPUSeriesN1Limit,
		n2:       CPUSeriesN2Limit,
		n2d:      CPUSeriesN2DLimit,
		e2:       CPUSeriesE2Limit,
		e2Micro:  CPUSeriesE2MicroLimit,
		e2Small:  CPUSeriesE2SmallLimit,
		e2Medium: CPUSeriesE2MediumLimit,
	}

	if !containsString([]string{e2Micro, e2Small, e2Medium}, cpuSeries) {
		return fmt.Errorf("incorrect cpu type: %v", cpuSeries), ""
	}

	coreCount := 2
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
	if containsString([]string{e2Small, e2Micro, e2Medium}, cpuSeries) {
		return nil, fmt.Sprintf("zones/%v/machineTypes/%v-%v", zone, cpuSeries, memory)
	}

	if memory > coreCount*typeLimit.maxMemPerCore {
		return nil, fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v-ext", zone, cpuSeries, coreCount, memory)
	}

	return nil, fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v", zone, cpuSeries, coreCount, memory)
}

// createInstanceWithCustomSharedCore creates a new VM instance with a custom type using shared CPUs.
func createInstanceWithCustomSharedCore(w io.Writer, projectID, zone, instanceName, cpuSeries string, memory int) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// cpuSeries := "e2-custom-micro" // the type of CPU you want to use"
	// memory := 256 // the amount of memory for the VM instance, in megabytes.

	err, machineType := customMachineTypeSharedCoreHelper(zone, cpuSeries, memory)
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

// [END compute_custom_machine_type_create_shared_with_helper]
