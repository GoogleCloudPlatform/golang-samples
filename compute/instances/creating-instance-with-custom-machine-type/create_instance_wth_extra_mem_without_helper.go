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

// [START compute_custom_machine_type_extra_mem_no_helper]
import (
	"fmt"
	"io"
)

// createInstanceWithExtraMemWithoutHelper сreates new VM instances with extra memory without using a CustomMachineType struct.
func createInstanceWithExtraMemWithoutHelper(w io.Writer, projectID, zone, instanceName string, coreCount, memory int) [](error) {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// coreCount := 2 // number of CPU cores you want to use.
	// memory := 256 // the amount of memory for the VM instance, in megabytes.

	// The coreCount and memory values are not validated anywhere and can be rejected by the API.

	instances := make([](error), 3)

	instances = append(instances, createInstanceWithCustomMachineType(w, projectID, zone, fmt.Sprintf("%s-n1-extra-mem", instanceName), fmt.Sprintf("zones/%s/machineTypes/custom-%v-%v-ext", zone, coreCount, memory)))
	instances = append(instances, createInstanceWithCustomMachineType(w, projectID, zone, fmt.Sprintf("%s-n2-extra-mem", instanceName), fmt.Sprintf("zones/%s/machineTypes/n2-custom-%v-%v-ext", zone, coreCount, memory)))
	instances = append(instances, createInstanceWithCustomMachineType(w, projectID, zone, fmt.Sprintf("%s-n2d-extra-mem", instanceName), fmt.Sprintf("zones/%s/machineTypes/n2d-custom-%v-%v-ext", zone, coreCount, memory)))

	return instances
}

// [END compute_custom_machine_type_extra_mem_no_helper]
