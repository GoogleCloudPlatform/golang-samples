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

// [START batch_create_custom_service_account]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Creates and runs a job with service account
func createJobWithSA(w io.Writer, projectID, region, jobName, serviceAccountEmail string) error {
	// projectID := "your-project-id"
	// region := "us-central1"
	// jobName := "some-job"
	// serviceAccountEmail := "your-service-account@your-project-id.iam.gserviceaccount.com"

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	runn := &batchpb.Runnable{
		Executable: &batchpb.Runnable_Script_{
			Script: &batchpb.Runnable_Script{
				Command: &batchpb.Runnable_Script_Text{
					Text: "echo Hello world! This is task ${BATCH_TASK_INDEX}. This job has a total of ${BATCH_TASK_COUNT} tasks.",
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
		Runnables:     []*batchpb.Runnable{runn},
	}

	taskGroups := []*batchpb.TaskGroup{
		{
			TaskCount: 3,
			TaskSpec:  taskSpec,
		},
	}

	// Policies are used to define on what kind of virtual machines the tasks will run on.
	// Read more about machine types here: https://cloud.google.com/compute/docs/machine-types
	allocationPolicy := &batchpb.AllocationPolicy{
		Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{{
			PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy{
				Policy: &batchpb.AllocationPolicy_InstancePolicy{
					MachineType: "e2-standard-4",
				},
			},
		}},
		ServiceAccount: &batchpb.ServiceAccount{
			Email: serviceAccountEmail,
		},
	}

	// We use Cloud Logging as it's an out of the box available option
	job := batchpb.Job{
		TaskGroups:       taskGroups,
		AllocationPolicy: allocationPolicy,
		Labels:           map[string]string{"env": "testing", "type": "script"},
		LogsPolicy:       &batchpb.LogsPolicy{Destination: batchpb.LogsPolicy_CLOUD_LOGGING},
	}

	req := &batchpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, region),
		JobId:  jobName,
		Job:    &job,
	}

	created_job, err := batchClient.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create job: %w", err)
	}

	fmt.Fprintf(w, "Job created: %v\n", created_job)
	return nil
}

// [END batch_create_custom_service_account]
