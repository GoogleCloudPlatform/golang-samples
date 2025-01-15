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

// [START batch_create_job_with_template]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Creates and runs a job that executes the specified script
func createScriptJobWithTemplate(w io.Writer, projectID, region, jobName, templateLink string) error {
	// projectID := "your_project_id"
	// region := "us-central1"
	// jobName := "some-job"
	/* A link to an existing Instance Template. Acceptable formats:
	*  "projects/{project_id}/global/instanceTemplates/{template_name}"
	*  "{template_name}" - if the template is defined in the same project as used to create the Job.
	 */
	// template_link := "my-instance-template"

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	// Define what will be done as part of the job.
	command := &batchpb.Runnable_Script_Text{
		Text: "echo Hello world! This is task ${BATCH_TASK_INDEX}. This job has a total of ${BATCH_TASK_COUNT} tasks.",
	}

	// We can specify what resources are requested by each task.
	resources := &batchpb.ComputeResource{
		// CpuMilli is milliseconds per cpu-second. This means the task requires 2 whole CPUs.
		CpuMilli:  2000,
		MemoryMib: 16,
	}

	taskSpec := &batchpb.TaskSpec{
		Runnables: []*batchpb.Runnable{{
			Executable: &batchpb.Runnable_Script_{
				Script: &batchpb.Runnable_Script{Command: command},
			},
		}},
		ComputeResource: resources,
		MaxRunDuration: &durationpb.Duration{
			Seconds: 3600,
		},
		MaxRetryCount: 2,
	}

	// Tasks are grouped inside a job using TaskGroups.
	taskGroups := []*batchpb.TaskGroup{
		{
			TaskCount: 4,
			TaskSpec:  taskSpec,
		},
	}

	// Policies are used to define on what kind of virtual machines the tasks will run on.
	// In this case, we are going to use an Instance Template that defines the VM.
	allocationPolicy := &batchpb.AllocationPolicy{
		Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{{
			PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_InstanceTemplate{
				InstanceTemplate: templateLink,
			},
		}},
	}

	// We use Cloud Logging as it's an out of the box available option
	logsPolicy := &batchpb.LogsPolicy{
		Destination: batchpb.LogsPolicy_CLOUD_LOGGING,
	}

	jobLabels := map[string]string{"env": "testing", "type": "script"}

	// The job's parent is the region in which the job will run
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)

	job := batchpb.Job{
		TaskGroups:       taskGroups,
		AllocationPolicy: allocationPolicy,
		Labels:           jobLabels,
		LogsPolicy:       logsPolicy,
	}

	req := &batchpb.CreateJobRequest{
		Parent: parent,
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

// [END batch_create_job_with_template]
