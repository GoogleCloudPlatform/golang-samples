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

func TestCustomMachineTypeSnippets(t *testing.T) {
	zone := "europe-central2-b"

	var customMachineTypeTests = []struct {
		cpuSeries string
		memory    int
		cpu       int
		cpuLimit  typeLimit
		out       string
		outMt     string
	}{
		{
			cpuSeries: n1,
			memory:    8192,
			cpu:       8,
			cpuLimit:  cpuSeriesN1Limit,
			out:       "zones/europe-central2-b/machineTypes/custom-8-8192",
			outMt:     "custom-8-8192",
		},
		{
			cpuSeries: n2,
			memory:    4096,
			cpu:       4,
			cpuLimit:  cpuSeriesN2Limit,
			out:       "zones/europe-central2-b/machineTypes/n2-custom-4-4096",
			outMt:     "n2-custom-4-4096",
		},
		{
			cpuSeries: n2d,
			memory:    8192,
			cpu:       4,
			cpuLimit:  cpuSeriesN2Limit,
			out:       "zones/europe-central2-b/machineTypes/n2d-custom-4-8192",
			outMt:     "n2d-custom-4-8192",
		},
		{
			cpuSeries: e2,
			memory:    8192,
			cpu:       8,
			cpuLimit:  cpuSeriesN2Limit,
			out:       "zones/europe-central2-b/machineTypes/e2-custom-8-8192",
			outMt:     "e2-custom-8-8192",
		},
		{
			cpuSeries: e2Small,
			memory:    4096,
			cpu:       0,
			cpuLimit:  cpuSeriesE2SmallLimit,
			out:       "zones/europe-central2-b/machineTypes/e2-custom-small-4096",
			outMt:     "e2-custom-small-4096",
		},
		{
			cpuSeries: e2Micro,
			memory:    2048,
			cpu:       0,
			cpuLimit:  cpuSeriesE2MicroLimit,
			out:       "zones/europe-central2-b/machineTypes/e2-custom-micro-2048",
			outMt:     "e2-custom-micro-2048",
		},
		{
			cpuSeries: n2,
			memory:    638720,
			cpu:       8,
			cpuLimit:  cpuSeriesN2Limit,
			out:       "zones/europe-central2-b/machineTypes/n2-custom-8-638720-ext",
			outMt:     "n2-custom-8-638720-ext",
		},
	}

	for _, tt := range customMachineTypeTests {
		t.Run(tt.outMt, func(t *testing.T) {
			cmt, err := createCustomMachineType(zone, tt.cpuSeries, tt.memory, tt.cpu, tt.cpuLimit)
			if err != nil {
				t.Errorf("createCustomMachineType return error: %v", err)
			}
			if cmt.String() != tt.out {
				t.Errorf("got %q, want %q", cmt.String(), tt.out)
			}
			if cmt.machineType() != tt.outMt {
				t.Errorf("got %q, want %q", cmt.machineType(), tt.outMt)
			}
		})
	}
}

func TestCustomMachineTypeErrorsSnippets(t *testing.T) {
	zone := "europe-central2-b"

	// bad memory 256
	_, err := createCustomMachineType(zone, n1, 8194, 8, cpuSeriesN1Limit)
	want := "requested memory must be a multiple of 256 MB"
	if err.Error() != want {
		t.Errorf("createCustomMachineType should return error: %v %v", want, err)
	}

	// wrong cpu count
	_, err = createCustomMachineType(zone, n2, 8194, 66, cpuSeriesN2Limit)
	want = fmt.Sprintf("invalid number of cores requested. Allowed number of cores for %v is: %v", n2, cpuSeriesN2Limit.allowedCores)
	if err.Error() != want {
		t.Errorf("createCustomMachineType should return error: %v %v", want, err)
	}
}
