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

// [START batch_create_gpu_job]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Creates and runs a job with custom events
func createJobWithGPU(w io.Writer, projectID, jobName string) (*batchpb.Job, error) {
	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("batchClient error: %w", err)
	}
	defer batchClient.Close()

	runn1 := &batchpb.Runnable{
		Executable: &batchpb.Runnable_Script_{
			Script: &batchpb.Runnable_Script{
				Command: &batchpb.Runnable_Script_Text{
					Text: "echo Hello world from script 1 for task ${BATCH_TASK_INDEX}",
				},
			},
		},
	}

	taskSpec := &batchpb.TaskSpec{
		ComputeResource: &batchpb.ComputeResource{
			// CpuMilli is milliseconds per cpu-second. This means the task requires 2 whole CPUs.
			CpuMilli:  2000,
			MemoryMib: 16,
		},
		MaxRunDuration: &durationpb.Duration{
			Seconds: 3600,
		},
		MaxRetryCount: 2,
		Runnables:     []*batchpb.Runnable{runn1},
	}

	taskGroups := []*batchpb.TaskGroup{
		{
			TaskCount: 4,
			TaskSpec:  taskSpec,
		},
	}

	labels := map[string]string{"env": "testing", "type": "container"}

	// Policies are used to define on what kind of virtual machines the tasks will run on.
	// In this case, we tell the system to use "n1-standard-16" machine type.
	// Read more about machine types here: https://cloud.google.com/compute/docs/machine-types
	allocationPolicy := &batchpb.AllocationPolicy{
		Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{{
			PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy{
				Policy: &batchpb.AllocationPolicy_InstancePolicy{
					MachineType: "n1-standard-16",
					Accelerators: []*batchpb.AllocationPolicy_Accelerator{{
						// Note: not every accelerator is compatible with instance type
						// Read more here: https://cloud.google.com/compute/docs/gpus#t4-gpus
						Type:  "nvidia-tesla-t4",
						Count: 1,
					}},
				},
			},
			InstallGpuDrivers: true,
		}},
	}

	// We use Cloud Logging as it's an out of the box available option
	logsPolicy := &batchpb.LogsPolicy{
		Destination: batchpb.LogsPolicy_CLOUD_LOGGING,
	}

	job := &batchpb.Job{
		Name:             jobName,
		TaskGroups:       taskGroups,
		AllocationPolicy: allocationPolicy,
		Labels:           labels,
		LogsPolicy:       logsPolicy,
	}

	request := &batchpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, "us-central1"),
		JobId:  jobName,
		Job:    job,
	}

	created_job, err := batchClient.CreateJob(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("unable to create job: %w", err)
	}

	fmt.Fprintf(w, "Job created: %v\n", created_job)
	return created_job, nil
}

// [END batch_create_gpu_job]
