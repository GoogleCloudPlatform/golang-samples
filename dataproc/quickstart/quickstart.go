// Copyright 2019 Google LLC
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

// [START dataproc_quickstart]

// This quickstart shows how you can use the Cloud Dataproc Client library to create a
// Cloud Dataproc cluster, submit a PySpark job to the cluster, wait for the job to finish
// and finally delete the cluster.
//
// Usage:
//     go build
//     ./quickstart --project_id <PROJECT_ID> --region <REGION> \
//         --cluster_name <CLUSTER_NAME> --job_file_path <GCS_JOB_FILE_PATH>
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1"
)

func main() {
	var projectID, clusterName, region, jobFilePath string
	flag.StringVar(&projectID, "project_id", "", "Cloud Project ID, used for creating resources.")
	flag.StringVar(&region, "region", "", "Region that resources should be created in.")
	flag.StringVar(&clusterName, "cluster_name", "", "Name of Cloud Dataproc cluster to create.")
	flag.StringVar(&jobFilePath, "job_file_path", "", "Path to job file in GCS.")
	flag.Parse()

	ctx := context.Background()

	// Create the cluster client.
	endpoint := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("error creating the cluster client: %s\n", err)
	}

	// Create the cluster config.
	createReq := &dataprocpb.CreateClusterRequest{
		ProjectId: projectID,
		Region:    region,
		Cluster: &dataprocpb.Cluster{
			ProjectId:   projectID,
			ClusterName: clusterName,
			Config: &dataprocpb.ClusterConfig{
				MasterConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   1,
					MachineTypeUri: "n1-standard-1",
				},
				WorkerConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   2,
					MachineTypeUri: "n1-standard-1",
				},
			},
		},
	}

	// Create the cluster.
	createOp, err := clusterClient.CreateCluster(ctx, createReq)
	if err != nil {
		log.Fatalf("error submitting the cluster creation request: %v\n", err)
	}

	createResp, err := createOp.Wait(ctx)
	if err != nil {
		log.Fatalf("error creating the cluster: %v\n", err)
	}

	// Defer cluster deletion.
	defer func() {
		dReq := &dataprocpb.DeleteClusterRequest{
			ProjectId:   projectID,
			Region:      region,
			ClusterName: clusterName,
		}
		deleteOp, err := clusterClient.DeleteCluster(ctx, dReq)
		deleteOp.Wait(ctx)
		if err != nil {
			fmt.Printf("error deleting cluster %q: %v\n", clusterName, err)
			return
		}
		fmt.Printf("Cluster %q successfully deleted\n", clusterName)
	}()

	// Output a success message.
	fmt.Printf("Cluster created successfully: %q\n", createResp.ClusterName)

	// Create the job client.
	jobClient, err := dataproc.NewJobControllerClient(ctx, option.WithEndpoint(endpoint))

	// Create the job config.
	submitJobReq := &dataprocpb.SubmitJobRequest{
		ProjectId: projectID,
		Region:    region,
		Job: &dataprocpb.Job{
			Placement: &dataprocpb.JobPlacement{
				ClusterName: clusterName,
			},
			TypeJob: &dataprocpb.Job_PysparkJob{
				PysparkJob: &dataprocpb.PySparkJob{
					MainPythonFileUri: jobFilePath,
				},
			},
		},
	}

	submitJobResp, err := jobClient.SubmitJob(ctx, submitJobReq)
	if err != nil {
		fmt.Printf("error submitting job: %v\n", err)
		return
	}

	id := submitJobResp.Reference.JobId

	fmt.Printf("Submitted job %q\n", id)

	// These states all signify that a job has terminated, successfully or not.
	terminalStates := map[dataprocpb.JobStatus_State]bool{
		dataprocpb.JobStatus_ERROR:     true,
		dataprocpb.JobStatus_CANCELLED: true,
		dataprocpb.JobStatus_DONE:      true,
	}

	// We can create a timeout such that the job gets cancelled if not in a terminal state after a certain amount of time.
	timeout := 5 * time.Minute
	start := time.Now()

	var state dataprocpb.JobStatus_State
	for {
		if time.Since(start) > timeout {
			cancelReq := &dataprocpb.CancelJobRequest{
				ProjectId: projectID,
				Region:    region,
				JobId:     id,
			}

			if _, err := jobClient.CancelJob(ctx, cancelReq); err != nil {
				fmt.Printf("error cancelling job: %v\n", err)
			}
			fmt.Printf("job %q timed out after %d minutes\n", id, int64(timeout.Minutes()))
			return
		}

		getJobReq := &dataprocpb.GetJobRequest{
			ProjectId: projectID,
			Region:    region,
			JobId:     id,
		}
		getJobResp, err := jobClient.GetJob(ctx, getJobReq)
		if err != nil {
			fmt.Printf("error getting job %q with error: %v\n", id, err)
			return
		}
		state = getJobResp.Status.State
		if terminalStates[state] {
			break
		}

		// Sleep as to not excessively poll the API.
		time.Sleep(1 * time.Second)
	}

	// Cloud Dataproc job outget gets saved to a GCS bucket allocated to it.
	getCReq := &dataprocpb.GetClusterRequest{
		ProjectId:   projectID,
		Region:      region,
		ClusterName: clusterName,
	}

	resp, err := clusterClient.GetCluster(ctx, getCReq)
	if err != nil {
		fmt.Printf("error getting cluster %q: %v\n", clusterName, err)
		return
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("error creating storage client: %v\n", err)
		return
	}

	obj := fmt.Sprintf("google-cloud-dataproc-metainfo/%s/jobs/%s/driveroutput.000000000", resp.ClusterUuid, id)
	reader, err := storageClient.Bucket(resp.Config.ConfigBucket).Object(obj).NewReader(ctx)
	if err != nil {
		fmt.Printf("error reading job output: %v\n", err)
		return
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("could not read output from Dataproc Job %q\n", id)
		return
	}

	fmt.Printf("job %q finished with state %s:\n%s\n", id, state, body)
}

// [END dataproc_quickstart]
