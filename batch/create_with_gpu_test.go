// Copyright 2024 Google LLC
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
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJobWithGPU(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	jobName := fmt.Sprintf("test-job-go-%v-%v", time.Now().Format("2006-01-02"), r.Int())

	buf := &bytes.Buffer{}

	job, err := createJobWithGPU(buf, tc.ProjectID, jobName)

	if err != nil {
		t.Errorf("createJobWithGPU got err: %v", err)
	}

	instance := job.GetAllocationPolicy().GetInstances()
	accelerator := instance[0].GetPolicy().GetAccelerators()
	if !instance[0].InstallGpuDrivers || accelerator[0].Type != "nvidia-tesla-t4" || accelerator[0].Count != 1 {
		t.Errorf("Accelerator wasn't set")
	}

	if err := deleteJob(buf, tc.ProjectID, "us-central1", jobName); err != nil {
		t.Errorf("deleteJob got err: %v", err)
	}
}
