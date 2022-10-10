// Copyright 2021 Google LLC
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

// [START batch_create_script_job]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	batchpb "google.golang.org/genproto/googleapis/cloud/batch/v1"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Creates and runs a job that executes the specified script
func createScriptJob(w io.Writer, projectID, region, jobName string) error {
	// projectID := "your_project_id"
	// region := "us-central1"
	// jobName := "some-job"

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer batchClient.Close()

	// Define what will be done as part of the job.
	task := new(batchpb.TaskSpec)
	runnable := new(batchpb.Runnable)
	// TODO: flatten this
	runnable.Executable = &batchpb.Runnable_Script_{
		Script: &batchpb.Runnable_Script{
			Command: &batchpb.Runnable_Script_Text{
				Text: "echo Hello world! This is task ${BATCH_TASK_INDEX}. This job has a total of ${BATCH_TASK_COUNT} tasks.",
			},
			// You can also run a script from a file. Just remember, that needs to be a script that's
			// already on the VM that will be running the job. Using runnable.script.text and runnable.script.path is mutually exclusive.
			// Command: &batchpb.Runnable_Script_Path{
			// 	Path: "/tmp/test.sh",
			// },
		},
	}
	task.Runnables = []*batchpb.Runnable{runnable}

	// We can specify what resources are requested by each task.
	task.ComputeResource = &batchpb.ComputeResource {
		CpuMilli:  2000, // in milliseconds per cpu-second. This means the task requires 2 whole CPUs.
		MemoryMib: 16,
	}

	task.MaxRunDuration = &durationpb.Duration{
		Seconds: 3600,
	}
	task.MaxRetryCount = 2

	// Tasks are grouped inside a job using TaskGroups.
	group := new(batchpb.TaskGroup)
	group.TaskCount = 4;
	group.TaskSpec = task;

	// Policies are used to define on what kind of virtual machines the tasks will run on.
	// In this case, we tell the system to use "e2-standard-4" machine type.
	// Read more about machine types here: https://cloud.google.com/compute/docs/machine-types
	allocationPolicy := new(batchpb.AllocationPolicy)
	policy := new(batchpb.AllocationPolicy_InstancePolicy)
	policy.MachineType = "e2-standard-4"
	policyTemplate := new(batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy)
	policyTemplate.Policy = policy
	instances := new(batchpb.AllocationPolicy_InstancePolicyOrTemplate)
	instances.PolicyTemplate = policyTemplate
	allocationPolicy.Instances = []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{instances}

	job := new(batchpb.Job)
	job.TaskGroups = []*batchpb.TaskGroup{group}
	job.AllocationPolicy = allocationPolicy
	job.Labels = map[string]string{"env": "testing", "type": "script"}
	// We use Cloud Logging as it's an out of the box available option
	job.LogsPolicy = new(batchpb.LogsPolicy)
	job.LogsPolicy.Destination = batchpb.LogsPolicy_CLOUD_LOGGING

	// The job's parent is the region in which the job will run
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)

	req := &batchpb.CreateJobRequest{
		Parent: parent,
		JobId:  jobName,
		Job: job,
	}

	created_job, err := batchClient.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create job: %v", err)
	}

	fmt.Fprintf(w, "Job created: %v\n", created_job)

	return nil
}

// [END batch_create_script_job]
