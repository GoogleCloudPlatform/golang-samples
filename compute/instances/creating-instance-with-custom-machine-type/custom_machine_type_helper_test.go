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

	// N1
	cmt1, err := createCustomMachineType(zone, N1, 8192, 8, CPUSeries_N1_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult := fmt.Sprintf("zones/%s/machineTypes/custom-8-8192", zone)
	if cmt1.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt1.ToString(), expectedResult)
	}
	expectedResult = "custom-8-8192"
	if cmt1.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt1.ToShortString(), expectedResult)
	}

	// N2
	cmt2, err := createCustomMachineType(zone, N2, 4096, 4, CPUSeries_N2_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/n2-custom-4-4096", zone)
	if cmt2.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt2.ToString(), expectedResult)
	}
	expectedResult = "n2-custom-4-4096"
	if cmt2.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt2.ToShortString(), expectedResult)
	}

	// N2D
	cmt3, err := createCustomMachineType(zone, N2D, 8192, 4, CPUSeries_N2D_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/n2d-custom-4-8192", zone)
	if cmt3.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt3.ToString(), expectedResult)
	}
	expectedResult = "n2d-custom-4-8192"
	if cmt3.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt3.ToShortString(), expectedResult)
	}

	// E2
	cmt4, err := createCustomMachineType(zone, E2, 8192, 8, CPUSeries_E2_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/e2-custom-8-8192", zone)
	if cmt4.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt4.ToString(), expectedResult)
	}
	expectedResult = "e2-custom-8-8192"
	if cmt4.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt4.ToShortString(), expectedResult)
	}

	// E2 SMALL
	cmt5, err := createCustomMachineType(zone, E2_SMALL, 4096, 0, CPUSeries_E2_SMALL_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/e2-custom-small-4096", zone)
	if cmt5.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt5.ToString(), expectedResult)
	}
	expectedResult = "e2-custom-small-4096"
	if cmt5.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt5.ToShortString(), expectedResult)
	}

	// E2 MICRO
	cmt6, err := createCustomMachineType(zone, E2_MICRO, 2048, 0, CPUSeries_E2_MICRO_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/e2-custom-micro-2048", zone)
	if cmt6.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt6.ToString(), expectedResult)
	}
	expectedResult = "e2-custom-micro-2048"
	if cmt6.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt6.ToShortString(), expectedResult)
	}

	// E2 MEDIUM
	cmt7, err := createCustomMachineType(zone, E2_MEDIUM, 8192, 0, CPUSeries_E2_MEDIUM_Limit)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/e2-custom-medium-8192", zone)
	if cmt7.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt7.ToString(), expectedResult)
	}
	expectedResult = "e2-custom-medium-8192"
	if cmt7.ToShortString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt7.ToShortString(), expectedResult)
	}

	// bad memory 256
	_, err = createCustomMachineType(zone, N1, 8194, 8, CPUSeries_N1_Limit)
	expectedResult = "requested memory must be a multiple of 256 MB"
	if fmt.Sprint(err) != expectedResult {
		t.Errorf("createCustomMachineType should return error: %s %v", expectedResult, err)
	}

	// ext memory
	cmt8, _ := createCustomMachineType(zone, N2, 638720, 8, CPUSeries_N2_Limit)
	expectedResult = fmt.Sprintf("zones/%s/machineTypes/n2-custom-8-638720-ext", zone)
	if cmt8.ToString() != expectedResult {
		t.Errorf("createCustomMachineType got %q, want %q", cmt8.ToString(), expectedResult)
	}

	// bad cpu count
	_, err = createCustomMachineType(zone, N2, 8194, 66, CPUSeries_N2_Limit)
	expectedResult = fmt.Sprintf("invalid number of cores requested. Allowed number of cores for %v is: %v", N2, CPUSeries_N2_Limit.allowedCores)
	if fmt.Sprint(err) != expectedResult {
		t.Errorf("createCustomMachineType should return error: %s %v", expectedResult, err)
	}
}
