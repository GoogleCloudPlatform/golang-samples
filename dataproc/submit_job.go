// Copyright 2020 Google LLC
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

// Package dataproc shows how you can use the Dataproc Client library to manage
// Dataproc clusters. In this example, we'll show how to submit a Spark job a cluster.
package dataproc

// [START dataproc_submit_job]
import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/dataproc/apiv1/dataprocpb"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func submitJob(w io.Writer, projectID, region, clusterName string) error {
	// projectID := "your-project-id"
	// region := "us-central1"
	// clusterName := "your-cluster"
	ctx := context.Background()

	// Create the job client.
	endpoint := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	jobClient, err := dataproc.NewJobControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("error creating the job client: %s\n", err)
	}

	// Create the job config.
	submitJobReq := &dataprocpb.SubmitJobRequest{
		ProjectId: projectID,
		Region:    region,
		Job: &dataprocpb.Job{
			Placement: &dataprocpb.JobPlacement{
				ClusterName: clusterName,
			},
			TypeJob: &dataprocpb.Job_SparkJob{
				SparkJob: &dataprocpb.SparkJob{
					Driver: &dataprocpb.SparkJob_MainClass{
						MainClass: "org.apache.spark.examples.SparkPi",
					},
					JarFileUris: []string{"file:///usr/lib/spark/examples/jars/spark-examples.jar"},
					Args:        []string{"1000"},
				},
			},
		},
	}

	submitJobOp, err := jobClient.SubmitJobAsOperation(ctx, submitJobReq)
	if err != nil {
		return fmt.Errorf("error with request to submitting job: %w", err)
	}

	submitJobResp, err := submitJobOp.Wait(ctx)
	if err != nil {
		return fmt.Errorf("error submitting job: %w", err)
	}

	re := regexp.MustCompile("gs://(.+?)/(.+)")
	matches := re.FindStringSubmatch(submitJobResp.DriverOutputResourceUri)

	if len(matches) < 3 {
		return fmt.Errorf("regex error: %s", submitJobResp.DriverOutputResourceUri)
	}

	// Dataproc job output gets saved to a GCS bucket allocated to it.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating storage client: %w", err)
	}

	obj := fmt.Sprintf("%s.000000000", matches[2])
	reader, err := storageClient.Bucket(matches[1]).Object(obj).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("error reading job output: %w", err)
	}

	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("could not read output from Dataproc Job: %w", err)
	}

	fmt.Fprintf(w, "Job finished successfully: %s", body)

	return nil
}

// [END dataproc_submit_job]
