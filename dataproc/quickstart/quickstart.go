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

// This quickstart shows how you can use the Cloud Dataproc Client library to create a
// Cloud Dataproc cluster, submit a PySpark job to the cluster, wait for the job to finish
// and finally delete the cluster.

// [START dataproc_quickstart]
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

// [END dataproc_quickstart]
var (
	projectID string
	region string
	clusterName string
	jobFilePath string
)

func init () {
	flag.StringVar(&projectID, "project_id", "", "Cloud Project ID, used for creating resources.")
	flag.StringVar(&region, "region", "", "Region that resources should be created in.")
	flag.StringVar(&clusterName, "cluster_name", "", "Name of Cloud Dataproc cluster to create.")
	flag.StringVar(&jobFilePath, "job_file_path", "", "Path to job file in GCS.")
}
// [START dataproc_quickstart]
func main() {
	// TODO (Developer): Set the following variables.
	// projectID = "your-project-id"
	// region = "us-central1"
	// clusterName = "your-cluster-name"
	// jobFilePath = "gs://your_file/location"
	// [END dataproc_quickstart]
	flag.Parse()
	// [START dataproc_quickstart]
	ctx := context.Background()

	// Create the cluster client
	endpoint := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	fmt.Printf("projectId: %s", projectID)
	fmt.Printf("region: %s", region)
	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("error creating the cluster client: %s", err)
	}

	// Create the cluster config
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

	// Create the cluster
	createOp, err := clusterClient.CreateCluster(ctx, createReq)
	if err != nil {
		log.Fatalf("error submitting the cluster creation request: %v", err)
	}

	createResp, err := createOp.Wait(ctx)
	if err != nil {
		log.Fatalf("error creating the cluster: %v", err)
	}

	// Output a success message
	fmt.Printf("Cluster created successfully: %q", createResp.ClusterName)

	// Create the job client
	jobClient, err := dataproc.NewJobControllerClient(ctx, option.WithEndpoint(endpoint))

	// Create the job config
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
		fmt.Errorf("error submitting job: %v", err)
	}

	id := submitJobResp.Reference.JobId

	fmt.Printf("Submitted job %q", id)

	// These states all signify that a job has terminated, successfully or not
	terminalStates := map[dataprocpb.JobStatus_State]bool{
		dataprocpb.JobStatus_ERROR:     true,
		dataprocpb.JobStatus_CANCELLED: true,
		dataprocpb.JobStatus_DONE:      true,
	}

	// We can create a timeout such that the job gets cancelled if not in a terminal state after a certain amount of time
	timeout := 5 * time.Minute
	start := time.Now().UTC()

	state := submitJobResp.Status.State
	for !terminalStates[state] {
		if time.Now().UTC().Sub(start) > timeout {
			cancelReq := &dataprocpb.CancelJobRequest{
				ProjectId: projectID,
				Region:    region,
				JobId:     id,
			}
			jobClient.CancelJob(ctx, cancelReq)
			fmt.Errorf("job %q timed out after threshold of %d seconds", id, int64(timeout.Seconds()))
		}

		getJobReq := &dataprocpb.GetJobRequest{
			ProjectId: projectID,
			Region:    region,
			JobId:     id,
		}
		getJobResp, err := jobClient.GetJob(ctx, getJobReq)
		if err != nil {
			fmt.Errorf("error getting job %q with error: %v", id, err)
			break
		}
		state = getJobResp.Status.State

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
		log.Fatalf("error getting cluster %q with error: %v", clusterName, err)
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating storage client with error: %v", err)
	}

	obj := fmt.Sprintf("google-cloud-dataproc-metainfo/%s/jobs/%s/driveroutput.000000000", resp.ClusterUuid, id)
	reader, err := storageClient.Bucket(resp.Config.ConfigBucket).Object(obj).NewReader(ctx)
	if err != nil {
		log.Fatalf("error reading job output: %v", err)
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("could not read output from Dataproc Job %q", id)
	}

	fmt.Printf("Job %q finished with state %s:\n%s", id, state, body)

	// Delete the cluster once the job has terminated
	dReq := &dataprocpb.DeleteClusterRequest{
		ProjectId:   projectID,
		Region:      region,
		ClusterName: clusterName,
	}
	deleteOp, err := clusterClient.DeleteCluster(ctx, dReq)
	deleteOp.Wait(ctx)
	if err != nil {
		log.Fatalf("error deleting cluster %q: %v", clusterName, err)
	}
	fmt.Printf("Cluster %q successfully deleted", clusterName)
}

// [END dataproc_quickstart]
