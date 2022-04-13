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

import (
	"fmt"
	"testing"
)

var customMachineTypeTests = []struct {
	cpuSeries string
	memory    int
	cpu       int
	cpuLimit  TypeLimit
	out       string
	outShort  string
}{
	{N1, 8192, 8, CPUSeriesN1Limit, "zones/europe-central2-b/machineTypes/custom-8-8192", "custom-8-8192"},
	{N2, 4096, 4, CPUSeriesN2Limit, "zones/europe-central2-b/machineTypes/n2-custom-4-4096", "n2-custom-4-4096"},
	{N2D, 8192, 4, CPUSeriesN2DLimit, "zones/europe-central2-b/machineTypes/n2d-custom-4-8192", "n2d-custom-4-8192"},
	{E2, 8192, 8, CPUSeriesE2Limit, "zones/europe-central2-b/machineTypes/e2-custom-8-8192", "e2-custom-8-8192"},
	{E2Small, 4096, 0, CPUSeriesE2SmallLimit, "zones/europe-central2-b/machineTypes/e2-custom-small-4096", "e2-custom-small-4096"},
	{E2Micro, 2048, 0, CPUSeriesE2MicroLimit, "zones/europe-central2-b/machineTypes/e2-custom-micro-2048", "e2-custom-micro-2048"},
	{N2, 638720, 8, CPUSeriesN2Limit, "zones/europe-central2-b/machineTypes/n2-custom-8-638720-ext", "n2-custom-8-638720-ext"},
}

func TestCustomMachineTypeSnippets(t *testing.T) {
	zone := "europe-central2-b"

	for _, tt := range customMachineTypeTests {
		t.Run(tt.outShort, func(t *testing.T) {
			cmt, err := createCustomMachineType(zone, tt.cpuSeries, tt.memory, tt.cpu, tt.cpuLimit)
			if err != nil {
				t.Errorf("createCustomMachineType return error: %v", err)
			}
			if cmt.String() != tt.out {
				t.Errorf("got %q, want %q", cmt.String(), tt.out)
			}
			if cmt.ShortString() != tt.outShort {
				t.Errorf("got %q, want %q", cmt.ShortString(), tt.outShort)
			}
		})
	}
}

func TestCustomMachineTypeErrorsSnippets(t *testing.T) {
	zone := "europe-central2-b"

	// bad memory 256
	_, err := createCustomMachineType(zone, N1, 8194, 8, CPUSeriesN1Limit)
	expectedResult := "requested memory must be a multiple of 256 MB"
	if fmt.Sprint(err) != expectedResult {
		t.Errorf("createCustomMachineType should return error: %s %v", expectedResult, err)
	}

	// wrong cpu count
	_, err = createCustomMachineType(zone, N2, 8194, 66, CPUSeriesN2Limit)
	expectedResult = fmt.Sprintf("invalid number of cores requested. Allowed number of cores for %v is: %v", N2, CPUSeriesN2Limit.allowedCores)
	if fmt.Sprint(err) != expectedResult {
		t.Errorf("createCustomMachineType should return error: %s %v", expectedResult, err)
	}
}
