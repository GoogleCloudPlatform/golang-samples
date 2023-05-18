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

// Package dataproc shows how you can use the Cloud Dataproc Client library to manage
// Cloud Dataproc clusters. In this example, we'll show how to instantiate an inline
// workflow template.
package dataproc

// [START dataproc_instantiate_inline_workflow_template]
import (
	"context"
	"fmt"
	"io"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/dataproc/apiv1/dataprocpb"
	"google.golang.org/api/option"
)

func instantiateInlineWorkflowTemplate(w io.Writer, projectID, region string) error {
	// projectID := "your-project-id"
	// region := "us-central1"

	ctx := context.Background()

	// Create the cluster client.
	endpoint := region + "-dataproc.googleapis.com:443"
	workflowTemplateClient, err := dataproc.NewWorkflowTemplateClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("dataproc.NewWorkflowTemplateClient: %w", err)
	}
	defer workflowTemplateClient.Close()

	// Create jobs for the workflow.
	teragenJob := &dataprocpb.OrderedJob{
		JobType: &dataprocpb.OrderedJob_HadoopJob{
			HadoopJob: &dataprocpb.HadoopJob{
				Driver: &dataprocpb.HadoopJob_MainJarFileUri{
					MainJarFileUri: "file:///usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar",
				},
				Args: []string{
					"teragen",
					"1000",
					"hdfs:///gen/",
				},
			},
		},
		StepId: "teragen",
	}

	terasortJob := &dataprocpb.OrderedJob{
		JobType: &dataprocpb.OrderedJob_HadoopJob{
			HadoopJob: &dataprocpb.HadoopJob{
				Driver: &dataprocpb.HadoopJob_MainJarFileUri{
					MainJarFileUri: "file:///usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar",
				},
				Args: []string{
					"terasort",
					"hdfs:///gen/",
					"hdfs:///sort/",
				},
			},
		},
		StepId: "terasort",
		PrerequisiteStepIds: []string{
			"teragen",
		},
	}

	// Create the cluster placement.
	clusterPlacement := &dataprocpb.WorkflowTemplatePlacement{
		Placement: &dataprocpb.WorkflowTemplatePlacement_ManagedCluster{
			ManagedCluster: &dataprocpb.ManagedCluster{
				ClusterName: "my-managed-cluster",
				Config: &dataprocpb.ClusterConfig{
					GceClusterConfig: &dataprocpb.GceClusterConfig{
						// Leave "ZoneUri" empty for "Auto Zone Placement"
						// ZoneUri: ""
						ZoneUri: "us-central1-a",
					},
				},
			},
		},
	}

	// Create the Instantiate Inline Workflow Template Request.
	req := &dataprocpb.InstantiateInlineWorkflowTemplateRequest{
		Parent: fmt.Sprintf("projects/%s/regions/%s", projectID, region),
		Template: &dataprocpb.WorkflowTemplate{
			Jobs: []*dataprocpb.OrderedJob{
				teragenJob,
				terasortJob,
			},
			Placement: clusterPlacement,
		},
	}

	// Create the cluster.
	op, err := workflowTemplateClient.InstantiateInlineWorkflowTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("InstantiateInlineWorkflowTemplate: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("InstantiateInlineWorkflowTemplate.Wait: %w", err)
	}

	// Output a success message.
	fmt.Fprintf(w, "Workflow created successfully.")
	return nil
}

// [END dataproc_instantiate_inline_workflow_template]
