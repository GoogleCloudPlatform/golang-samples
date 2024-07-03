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
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

func customMachineTypeURI(zone, cpuSeries string, coreCount, memory int) (string, error) {
	const (
		n1       = "custom"
		n2       = "n2-custom"
		n2d      = "n2d-custom"
		e2       = "e2-custom"
		e2Micro  = "e2-custom-micro"
		e2Small  = "e2-custom-small"
		e2Medium = "e2-custom-medium"
	)

	type typeLimit struct {
		allowedCores     []int
		minMemPerCore    int
		maxMemPerCore    int
		allowExtraMemory bool
		extraMemoryLimit int
	}

	makeRange := func(start, end, step int) []int {
		if step <= 0 || end < start {
			return []int{}
		}
		s := make([]int, 0, 1+(end-start)/step)
		for start <= end {
			s = append(s, start)
			start += step
		}
		return s
	}

	containsString := func(s []string, str string) bool {
		for _, v := range s {
			if v == str {
				return true
			}
		}

		return false
	}

	containsInt := func(nums []int, n int) bool {
		for _, v := range nums {
			if v == n {
				return true
			}
		}

		return false
	}

	var (
		cpuSeriesE2Limit = typeLimit{
			allowedCores:  makeRange(2, 33, 2),
			minMemPerCore: 512,
			maxMemPerCore: 8192,
		}
		cpuSeriesE2MicroLimit  = typeLimit{minMemPerCore: 1024, maxMemPerCore: 2048}
		cpuSeriesE2SmallLimit  = typeLimit{minMemPerCore: 2048, maxMemPerCore: 4096}
		cpuSeriesE2MeidumLimit = typeLimit{minMemPerCore: 4096, maxMemPerCore: 8192}
		cpuSeriesN2Limit       = typeLimit{
			allowedCores:  append(makeRange(2, 33, 2), makeRange(36, 129, 4)...),
			minMemPerCore: 512, maxMemPerCore: 8192,
			allowExtraMemory: true,
			extraMemoryLimit: 624 << 10,
		}
		cpuSeriesN2DLimit = typeLimit{
			allowedCores:  []int{2, 4, 8, 16, 32, 48, 64, 80, 96},
			minMemPerCore: 512, maxMemPerCore: 8192,
			allowExtraMemory: true,
			extraMemoryLimit: 768 << 10,
		}
		cpuSeriesN1Limit = typeLimit{
			allowedCores:     append([]int{1}, makeRange(2, 97, 2)...),
			minMemPerCore:    922,
			maxMemPerCore:    6656,
			allowExtraMemory: true,
			extraMemoryLimit: 624 << 10,
		}
	)

	typeLimitsMap := map[string]typeLimit{
		n1:       cpuSeriesN1Limit,
		n2:       cpuSeriesN2Limit,
		n2d:      cpuSeriesN2DLimit,
		e2:       cpuSeriesE2Limit,
		e2Micro:  cpuSeriesE2MicroLimit,
		e2Small:  cpuSeriesE2SmallLimit,
		e2Medium: cpuSeriesE2MeidumLimit,
	}

	if !containsString([]string{e2, n1, n2, n2d}, cpuSeries) {
		return "", fmt.Errorf("incorrect cpu type: %v", cpuSeries)
	}

	tl := typeLimitsMap[cpuSeries]

	// Check whether the requested parameters are allowed.
	// Find more information about limitations of custom machine types at:
	// https://cloud.google.com/compute/docs/general-purpose-machines#custom_machine_types

	// Check the number of cores
	if len(tl.allowedCores) > 0 && !containsInt(tl.allowedCores, coreCount) {
		return "", fmt.Errorf(
			"invalid number of cores requested. Allowed number of cores for %v is: %v",
			cpuSeries,
			tl.allowedCores,
		)
	}

	// Memory must be a multiple of 256 MB
	if memory%256 != 0 {
		return "", fmt.Errorf("requested memory must be a multiple of 256 MB")
	}

	// Check if the requested memory isn't too little
	if memory < coreCount*tl.minMemPerCore {
		return "", fmt.Errorf(
			"requested memory is too low. Minimal memory for %v is %v MB per core",
			cpuSeries,
			tl.minMemPerCore,
		)
	}

	// Check if the requested memory isn't too much
	if memory > coreCount*tl.maxMemPerCore && !tl.allowExtraMemory {
		return "", fmt.Errorf(
			"requested memory is too large.. Maximum memory allowed for %v is %v MB per core",
			cpuSeries,
			tl.maxMemPerCore,
		)
	}
	if memory > tl.extraMemoryLimit && tl.allowExtraMemory {
		return "", fmt.Errorf(
			"requested memory is too large.. Maximum memory allowed for %v is %v MB",
			cpuSeries,
			tl.extraMemoryLimit,
		)
	}

	// Return the custom machine type in form of a string acceptable by Compute Engine API.
	if containsString([]string{e2Small, e2Micro, e2Medium}, cpuSeries) {
		return fmt.Sprintf("zones/%v/machineTypes/%v-%v", zone, cpuSeries, memory), nil
	}

	if memory > coreCount*tl.maxMemPerCore {
		return fmt.Sprintf(
			"zones/%v/machineTypes/%v-%v-%v-ext",
			zone,
			cpuSeries,
			coreCount,
			memory,
		), nil
	}

	return fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v", zone, cpuSeries, coreCount, memory), nil
}

// createInstanceWithCustomMachineTypeWithHelper creates a new VM instance with a custom machine type.
func createInstanceWithCustomMachineTypeWithHelper(
	w io.Writer,
	projectID, zone, instanceName, cpuSeries string,
	coreCount, memory int,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// cpuSeries := "e2-custom-micro" // the type of CPU you want to use"
	// coreCount := 2 // number of CPU cores you want to use.
	// memory := 256 // the amount of memory for the VM instance, in megabytes.

	machineType, err := customMachineTypeURI(zone, cpuSeries, coreCount, memory)
	if err != nil {
		return fmt.Errorf("unable to create custom machine type string: %w", err)
	}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
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
						DiskSizeGb: proto.Int64(10),
						SourceImage: proto.String(
							"projects/debian-cloud/global/images/family/debian-12",
						),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
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
		return fmt.Errorf("unable to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Instance created\n")

	return nil
}

// [END compute_custom_machine_type_create_with_helper]
