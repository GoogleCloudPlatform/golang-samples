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

// [START batch_create_local_ssd_job]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Creates and runs a job with local SSD
// Note: local SSD does not guarantee Local SSD data persistence.
// More details here: https://cloud.google.com/compute/docs/disks/local-ssd#data_persistence
func createJobWithSSD(w io.Writer, projectID, jobName, ssdName string) error {
	// jobName := job-name
	// ssdName := disk-name
	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("batchClient error: %w", err)
	}
	defer batchClient.Close()

	runn := &batchpb.Runnable{
		Executable: &batchpb.Runnable_Script_{
			Script: &batchpb.Runnable_Script{
				Command: &batchpb.Runnable_Script_Text{
					Text: "echo Hello world from script 1 for task ${BATCH_TASK_INDEX}",
				},
			},
		},
	}
	volume := &batchpb.Volume{
		MountPath: fmt.Sprintf("/mnt/disks/%v", ssdName),
		Source: &batchpb.Volume_DeviceName{
			DeviceName: ssdName,
		},
	}

	// The size of all the local SSDs in GB. Each local SSD is 375 GB,
	// so this value must be a multiple of 375 GB.
	// For example, for 2 local SSDs, set this value to 750 GB.
	disk := &batchpb.AllocationPolicy_Disk{
		Type:   "local-ssd",
		SizeGb: 375,
	}

	taskSpec := &batchpb.TaskSpec{
		ComputeResource: &batchpb.ComputeResource{
			// CpuMilli is milliseconds per cpu-second. This means the task requires 1 CPU.
			CpuMilli:  1000,
			MemoryMib: 16,
		},
		MaxRunDuration: &durationpb.Duration{
			Seconds: 3600,
		},
		MaxRetryCount: 2,
		Runnables:     []*batchpb.Runnable{runn},
		Volumes:       []*batchpb.Volume{volume},
	}

	taskGroups := []*batchpb.TaskGroup{
		{
			TaskCount: 4,
			TaskSpec:  taskSpec,
		},
	}

	labels := map[string]string{"env": "testing", "type": "container"}

	allocationPolicy := &batchpb.AllocationPolicy{
		Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{{
			PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy{
				Policy: &batchpb.AllocationPolicy_InstancePolicy{
					// The allowed number of local SSDs depends on the machine type for your job's VMs.
					// In this case, we tell the system to use "n1-standard-1" machine type, which require to attach local ssd manually.
					// Read more about local disks here: https://cloud.google.com/compute/docs/disks/local-ssd#lssd_disk_options
					MachineType: "n1-standard-1",
					Disks: []*batchpb.AllocationPolicy_AttachedDisk{
						{
							Attached: &batchpb.AllocationPolicy_AttachedDisk_NewDisk{
								NewDisk: disk,
							},
							DeviceName: ssdName,
						},
					},
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
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, "us-central1"),
		JobId:  jobName,
		Job:    job,
	}

	created_job, err := batchClient.CreateJob(ctx, request)
	if err != nil {
		return fmt.Errorf("unable to create job: %w", err)
	}

	fmt.Fprintf(w, "Job created: %v\n", created_job)
	return nil
}

// [END batch_create_local_ssd_job]
