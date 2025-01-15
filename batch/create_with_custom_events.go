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

// [START batch_custom_events]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Creates and runs a job with custom events
func createJobWithCustomEvents(w io.Writer, projectID, jobName string) (*batchpb.Job, error) {
	region := "us-central1"
	displayName1 := "script 1"
	displayName2 := "barrier 1"
	displayName3 := "script 2"

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
		DisplayName: displayName1,
	}

	runn2 := &batchpb.Runnable{
		Executable: &batchpb.Runnable_Barrier_{
			Barrier: &batchpb.Runnable_Barrier{},
		},
		DisplayName: displayName2,
	}

	runn3 := &batchpb.Runnable{
		Executable: &batchpb.Runnable_Script_{
			Script: &batchpb.Runnable_Script{
				Command: &batchpb.Runnable_Script_Text{
					Text: "echo Hello world from script 2 for task ${BATCH_TASK_INDEX}",
				},
			},
		},
		DisplayName: displayName3,
	}

	runn4 := &batchpb.Runnable{
		Executable: &batchpb.Runnable_Script_{
			Script: &batchpb.Runnable_Script{
				Command: &batchpb.Runnable_Script_Text{
					Text: "sleep 30; echo '{\"batch/custom/event\": \"DESCRIPTION\"}'; sleep 30",
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
		Runnables:     []*batchpb.Runnable{runn1, runn2, runn3, runn4},
	}

	taskGroups := []*batchpb.TaskGroup{
		{
			TaskCount: 4,
			TaskSpec:  taskSpec,
		},
	}

	labels := map[string]string{"env": "testing", "type": "container"}

	// Policies are used to define on what kind of virtual machines the tasks will run on.
	// In this case, we tell the system to use "e2-standard-4" machine type.
	// Read more about machine types here: https://cloud.google.com/compute/docs/machine-types
	allocationPolicy := &batchpb.AllocationPolicy{
		Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{{
			PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy{
				Policy: &batchpb.AllocationPolicy_InstancePolicy{
					MachineType: "e2-standard-4",
				},
			},
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
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, region),
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

// [END batch_custom_events]
