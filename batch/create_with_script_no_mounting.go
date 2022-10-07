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

// [START compute_instances_create]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	batchpb "google.golang.org/genproto/googleapis/cloud/batch/v1"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	//"google.golang.org/protobuf/proto"
)

// TODO: documentation
func createInstance(w io.Writer, projectID, region, jobName string) error {
	// projectID := "your_project_id"
	// region := "us-central1"
	// jobName := TODO

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer batchClient.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)

	req := &batchpb.CreateJobRequest{
		Parent: parent,
		JobId: jobName,
		Job: &batchpb.Job{
			TaskGroups: []*batchpb.TaskGroup{
				{
					TaskCount: 4,
					TaskSpec: &batchpb.TaskSpec{
						Runnables: []*batchpb.Runnable{{
							Executable: &batchpb.Runnable_Script_{
								Script: &batchpb.Runnable_Script{
									Command: &batchpb.Runnable_Script_Text{
										Text: "echo Hello world! This is task ${BATCH_TASK_INDEX}. This job has a total of ${BATCH_TASK_COUNT} tasks.",
									},
								},
							},
						}},
						ComputeResource:   &batchpb.ComputeResource{
							CpuMilli:    2000, // in milliseconds per cpu-second. This means the task requires 2 whole CPUs.
							MemoryMib:   16,
						},
						MaxRunDuration:    &durationpb.Duration{
							Seconds: 3600,
						},
						MaxRetryCount:     2,
					},
				},
			},
			AllocationPolicy: &batchpb.AllocationPolicy{
				Location:  &batchpb.AllocationPolicy_LocationPolicy{},
				Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{{
					PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy{
						Policy: &batchpb.AllocationPolicy_InstancePolicy{
							MachineType:       "e2-standard-4",
						},
					},
				}},
			},
			Labels:           map[string]string{"env": "testing", "type": "script"},
			LogsPolicy:       &batchpb.LogsPolicy{
				Destination: batchpb.LogsPolicy_CLOUD_LOGGING,
			},
		},
		// InstanceResource: &computepb.Instance{
		// 	Name: proto.String(instanceName),
		// 	Disks: []*computepb.AttachedDisk{
		// 		{
		// 			InitializeParams: &computepb.AttachedDiskInitializeParams{
		// 				DiskSizeGb:  proto.Int64(10),
		// 				SourceImage: proto.String(sourceImage),
		// 			},
		// 			AutoDelete: proto.Bool(true),
		// 			Boot:       proto.Bool(true),
		// 			Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
		// 		},
		// 	},
		// 	MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
		// 	NetworkInterfaces: []*computepb.NetworkInterface{
		// 		{
		// 			Name: proto.String(networkName),
		// 		},
		// 	},
		// },
	}

	job, err := batchClient.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %v", err)
	}

	fmt.Fprintf(w, "Job created: %v\n", job)

	return nil
}

// [END compute_instances_create]
