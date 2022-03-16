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

// [START compute_custom_machine_type_helper_class]
import (
	"fmt"
	"strings"
)

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

type CustomMachineType struct {
	zone, cpuSeries     string
	memoryMb, coreCount int
	typeLimit           TypeLimit
}

func (t CustomMachineType) _Check() error {
	// Check whether the requested parameters are allowed. Find more information about limitations of custom machine
	// types at: https://cloud.google.com/compute/docs/general-purpose-machines#custom_machine_types

	// Check the number of cores
	if len(t.typeLimit.allowedCores) > 0 && !contains_int(t.typeLimit.allowedCores, t.coreCount) {
		return fmt.Errorf("invalid number of cores requested. Allowed number of cores for %v is: %v", t.cpuSeries, t.typeLimit.allowedCores)
	}

	// Memory must be a multiple of 256 MB
	if t.memoryMb%256 != 0 {
		return fmt.Errorf("requested memory must be a multiple of 256 MB")
	}

	// Check if the requested memory isn't too little
	if t.memoryMb < t.coreCount*t.typeLimit.minMemPerCore {
		return fmt.Errorf("requested memory is too low. Minimal memory for %v is %v MB per core", t.cpuSeries, t.typeLimit.minMemPerCore)
	}

	// Check if the requested memory isn't too much
	if t.memoryMb > t.coreCount*t.typeLimit.maxMemPerCore {
		if t.typeLimit.allowExtraMemory {
			if t.memoryMb > t.typeLimit.extraMemoryLimit {
				return fmt.Errorf("requested memory is too large.. Maximum memory allowed for %v is %v MB", t.cpuSeries, t.typeLimit.extraMemoryLimit)
			}
		} else {
			return fmt.Errorf("requested memory is too large.. Maximum memory allowed for %v is %v MB per core", t.cpuSeries, t.typeLimit.maxMemPerCore)
		}
	}

	return nil
}

func (t CustomMachineType) IsExtraMemoryUsed() bool {
	return t.memoryMb > t.coreCount*t.typeLimit.maxMemPerCore
}

func (t CustomMachineType) ToString() string {
	// Return the custom machine type in form of a string acceptable by Compute Engine API.

	if contains_string([]string{E2_SMALL, E2_MICRO, E2_MEDIUM}, t.cpuSeries) {
		return fmt.Sprintf("zones/%v/machineTypes/%v-%v", t.zone, t.cpuSeries, t.memoryMb)
	}

	if t.IsExtraMemoryUsed() {
		return fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v-ext", t.zone, t.cpuSeries, t.coreCount, t.memoryMb)
	}

	return fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v", t.zone, t.cpuSeries, t.coreCount, t.memoryMb)
}

func (t CustomMachineType) ToShortString() string {
	// Return machine type in a format without the zone. For example, n2-custom-0-10240.
	// This format is used to create instance templates.
	ss := strings.Split(t.ToString(), "/")
	return ss[len(ss)-1]
}

func createCustomMachineType(zone, cpuSeries string, memoryMb, coreCount int, tl TypeLimit) (*CustomMachineType, error) {
	if contains_string([]string{E2_SMALL, E2_MICRO, E2_MEDIUM}, cpuSeries) {
		coreCount = 2
	}
	custom_mt := CustomMachineType{zone, cpuSeries, memoryMb, coreCount, tl}
	check := custom_mt._Check()

	if check != nil {
		return &CustomMachineType{}, check
	}
	return &custom_mt, nil
}

// [START compute_custom_machine_type_helper_class]
