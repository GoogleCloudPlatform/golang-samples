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

// This sample shows how you can use the Cloud Dataproc Client library to create a
// Cloud Dataproc cluster

// [START dataproc_quickstart]
package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1"
)

func quickstart(w io.Writer, projectID, region, clusterName, jobFilePath string) error {
	ctx := context.Background()

	// Create the cluster client
	endpoint := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("Error creating the cluster client: %s", err)
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
		return fmt.Errorf("Error submitting the cluster creation request: %v", err)
	}

	createResp, err := createOp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error creating the cluster: %v", err)
	}

	// Output a success message
	fmt.Fprintf(w, "Cluster created successfully: %q", createResp.ClusterName)

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
		return fmt.Errorf("Error submitting job: %v", err)
	}

	id := submitJobResp.Reference.JobId

	fmt.Fprintf(w, "Submitted job %q", id)

	// These states all signify that a job has terminated, successfully or not
	ts := map[dataprocpb.JobStatus_State]bool{
		dataprocpb.JobStatus_ERROR:     true,
		dataprocpb.JobStatus_CANCELLED: true,
		dataprocpb.JobStatus_DONE:      true,
	}

	// We can create a timeout such that the job gets cancelled if not in a terminal state after a certain amount of time
	timeout := 5 * time.Minute
	start := time.Now().UTC()

	state := submitJobResp.Status.State
	for _, ok := ts[state]; !ok; _, ok = ts[state] {
		// fmt.Fprintf(w, "-----\n")
		// fmt.Fprintf(w, "State: %s", state)
		// fmt.Fprintf(w, "out: %t, ok: %t", out, ok)
		// fmt.Fprintf(w, "-----\n")
		if time.Now().UTC().Sub(start) > timeout {
			cancelReq := &dataprocpb.CancelJobRequest{ProjectId: projectID, Region: region, JobId: id}
			jobClient.CancelJob(ctx, cancelReq)
			return fmt.Errorf("Job %q timed out after threshold of %d seconds", id, int64(timeout.Seconds()))
		}

		time.Sleep(1 * time.Second)

		getJobReq := &dataprocpb.GetJobRequest{
			ProjectId: projectID,
			Region:    region,
			JobId:     id,
		}
		getJobResp, err := jobClient.GetJob(ctx, getJobReq)
		if err != nil {
			return fmt.Errorf("Error getting job %q with error: %v", id, err)
		}

		state = getJobResp.Status.State
	}

	// Cloud Dataproc job outget gets saved to a GCS bucket allocated to it.
	getCReq := &dataprocpb.GetClusterRequest{
		ProjectId:   projectID,
		Region:      region,
		ClusterName: clusterName,
	}

	resp, err := clusterClient.GetCluster(ctx, getCReq)
	if err != nil {
		return fmt.Errorf("Error getting cluster %q with error: %v", clusterName, err)
	}

	sc, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Error creating storage client with error: %v", err)
	}

	obj := fmt.Sprintf("google-cloud-dataproc-metainfo/%s/jobs/%s/driveroutput.000000000", resp.ClusterUuid, id)
	rc, err := sc.Bucket(resp.Config.ConfigBucket).Object(obj).NewReader(ctx)
	if err != nil {
		fmt.Fprintf(w, "UUID: "+resp.ClusterUuid)
		fmt.Fprintf(w, "id: "+id)
		return fmt.Errorf("Error reading Dataproc output: %v", err)
	}

	defer rc.Close()

	body, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("Could not read output from Dataproc Job %q", id)
	}

	fmt.Fprintf(w, "Job %q finished with state %s:\n%s", id, state, body)

	// Delete the cluster once the job has terminated
	dReq := &dataprocpb.DeleteClusterRequest{
		ProjectId:   projectID,
		Region:      region,
		ClusterName: clusterName,
	}
	deleteOp, err := clusterClient.DeleteCluster(ctx, dReq)
	deleteOp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting cluster %q: %v", clusterName, err)
	}
	fmt.Fprintf(w, "Cluster %q successfully deleted", clusterName)

	return nil
}

// [END dataproc_quickstart]
