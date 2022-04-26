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

func makeRange(start, end, step int) []int {
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

type customMachineType struct {
	zone, cpuSeries     string
	memoryMb, coreCount int
	typeLimit
}

// Validate whether the requested parameters are allowed.
// Find more information about limitations of custom machine types at:
// https://cloud.google.com/compute/docs/general-purpose-machines#custom_machine_types
func validate(cmt *customMachineType) error {
	// Check the number of cores
	if len(cmt.typeLimit.allowedCores) > 0 {
		coreExists := false
		for _, v := range cmt.typeLimit.allowedCores {
			if v == cmt.coreCount {
				coreExists = true
			}
		}
		if !coreExists {
			return fmt.Errorf("invalid number of cores requested. Allowed number of cores for %v is: %v", cmt.cpuSeries, cmt.typeLimit.allowedCores)
		}
	}

	// Memory must be a multiple of 256 MB
	if cmt.memoryMb%256 != 0 {
		return fmt.Errorf("requested memory must be a multiple of 256 MB")
	}

	// Check if the requested memory isn't too little
	if cmt.memoryMb < cmt.coreCount*cmt.typeLimit.minMemPerCore {
		return fmt.Errorf("requested memory is too low. Minimal memory for %v is %v MB per core", cmt.cpuSeries, cmt.typeLimit.minMemPerCore)
	}

	// Check if the requested memory isn't too much
	if cmt.memoryMb > cmt.coreCount*cmt.typeLimit.maxMemPerCore && !cmt.typeLimit.allowExtraMemory {
		return fmt.Errorf("requested memory is too large. Maximum memory allowed for %v is %v MB per core", cmt.cpuSeries, cmt.typeLimit.maxMemPerCore)
	}
	if cmt.memoryMb > cmt.typeLimit.extraMemoryLimit && cmt.typeLimit.allowExtraMemory {
		return fmt.Errorf("requested memory is too large. Maximum memory allowed for %v is %v MB", cmt.cpuSeries, cmt.typeLimit.extraMemoryLimit)
	}

	return nil
}

// Srring returns the custom machine type in form of a string acceptable by Compute Engine API.
func (t customMachineType) String() string {
	containsString := func(s []string, str string) bool {
		for _, v := range s {
			if v == str {
				return true
			}
		}

		return false
	}

	if containsString([]string{e2Small, e2Micro, e2Medium}, t.cpuSeries) {
		return fmt.Sprintf("zones/%v/machineTypes/%v-%v", t.zone, t.cpuSeries, t.memoryMb)
	}

	if t.memoryMb > t.coreCount*t.typeLimit.maxMemPerCore {
		return fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v-ext", t.zone, t.cpuSeries, t.coreCount, t.memoryMb)
	}

	return fmt.Sprintf("zones/%v/machineTypes/%v-%v-%v", t.zone, t.cpuSeries, t.coreCount, t.memoryMb)
}

// Returns machine type in a format without the zone. For example, n2-custom-0-10240.
// This format is used to create instance templates.
func (t customMachineType) machineType() string {
	// Return machine type in a format without the zone. For example, n2-custom-0-10240.
	// This format is used to create instance templates.
	ss := strings.Split(t.String(), "/")
	return ss[len(ss)-1]
}

func createCustomMachineType(zone, cpuSeries string, memoryMb, coreCount int, tl typeLimit) (*customMachineType, error) {
	containsString := func(s []string, str string) bool {
		for _, v := range s {
			if v == str {
				return true
			}
		}

		return false
	}

	if containsString([]string{e2Small, e2Micro, e2Medium}, cpuSeries) {
		coreCount = 2
	}
	cmt := &customMachineType{
		zone:      zone,
		cpuSeries: cpuSeries,
		memoryMb:  memoryMb,
		coreCount: coreCount,
		typeLimit: tl,
	}

	if err := validate(cmt); err != nil {
		return &customMachineType{}, err
	}
	return cmt, nil
}

// [END compute_custom_machine_type_helper_class]
