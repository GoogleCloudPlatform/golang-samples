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
	"io"
)

// createInstanceWithCustomMachineTypeWithHelper creates a new VM instance with a custom machine type.
func createInstanceWithCustomMachineTypeWithHelper(w io.Writer, projectID, zone, instanceName, cpuSeries string, coreCount, memory int, tl TypeLimit) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// cpuSeries := "e2-custom-micro" // the type of CPU you want to use"
	// coreCount := 2 // number of CPU cores you want to use.
	// memory := 256 // the amount of memory for the VM instance, in megabytes.

	custom_mt, err := createCustomMachineType(zone, cpuSeries, memory, coreCount, tl)
	if err != nil {
		return err
	}

	return createInstanceWithCustomMachineType(w, projectID, zone, instanceName, custom_mt.ToString())
}

// [END compute_custom_machine_type_create_with_helper]
