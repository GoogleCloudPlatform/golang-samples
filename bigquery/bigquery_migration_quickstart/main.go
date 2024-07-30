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

// [START bigquerymigration_quickstart]

// The bigquery_migration_quickstart application demonstrates basic usage of the
// BigQuery migration API by executing a workflow that performs a batch SQL
// translation task.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	migration "cloud.google.com/go/bigquery/migration/apiv2"
	"cloud.google.com/go/bigquery/migration/apiv2/migrationpb"
)

func main() {
	// Define command line flags for controlling the behavior of this quickstart.
	projectID := flag.String("project_id", "", "Cloud Project ID.")
	location := flag.String("location", "us", "BigQuery Migration location used for interactions.")
	outputPath := flag.String("output", "", "Cloud Storage path for translated resources.")
	// Parse flags and do some minimal validation.
	flag.Parse()
	if *projectID == "" {
		log.Fatal("empty --project_id specified, please provide a valid project ID")
	}
	if *location == "" {
		log.Fatal("empty --location specified, please provide a valid location")
	}
	if *outputPath == "" {
		log.Fatalf("empty --output specified, please provide a valid cloud storage path")
	}

	ctx := context.Background()
	migClient, err := migration.NewClient(ctx)
	if err != nil {
		log.Fatalf("migration.NewClient: %v", err)
	}
	defer migClient.Close()

	workflow, err := executeTranslationWorkflow(ctx, migClient, *projectID, *location, *outputPath)
	if err != nil {
		log.Fatalf("workflow execution failed: %v\n", err)
	}

	reportWorkflowStatus(workflow)
}

// executeTranslationWorkflow constructs a migration workflow that performs batch SQL translation.
func executeTranslationWorkflow(ctx context.Context, client *migration.Client, projectID, location, outPath string) (*migrationpb.MigrationWorkflow, error) {

	// Construct the workflow creation request.  In this workflow, we have only a single translation task present.
	req := &migrationpb.CreateMigrationWorkflowRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		MigrationWorkflow: &migrationpb.MigrationWorkflow{
			DisplayName: "example SQL conversion",
			Tasks: map[string]*migrationpb.MigrationTask{
				"example_conversion": {
					Type: "Translation_Teradata2BQ",
					TaskDetails: &migrationpb.MigrationTask_TranslationConfigDetails{
						TranslationConfigDetails: &migrationpb.TranslationConfigDetails{
							SourceLocation: &migrationpb.TranslationConfigDetails_GcsSourcePath{
								GcsSourcePath: "gs://cloud-samples-data/bigquery/migration/translation/input/",
							},
							TargetLocation: &migrationpb.TranslationConfigDetails_GcsTargetPath{
								GcsTargetPath: outPath,
							},
							SourceDialect: &migrationpb.Dialect{
								DialectValue: &migrationpb.Dialect_TeradataDialect{
									TeradataDialect: &migrationpb.TeradataDialect{
										Mode: migrationpb.TeradataDialect_SQL,
									},
								},
							},
							TargetDialect: &migrationpb.Dialect{
								DialectValue: &migrationpb.Dialect_BigqueryDialect{},
							},
						},
					},
				},
			},
		},
	}

	// Create the workflow using the request.
	workflow, err := client.CreateMigrationWorkflow(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateMigrationWorkflow: %w", err)
	}
	fmt.Printf("workflow created: %s", workflow.GetName())

	// This is an asyncronous process, so we now poll the workflow
	// until completion or a suitable timeout has elapsed.
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	for {
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("task %s didn't complete due to context expiring", workflow.GetName())
		default:
			polledWorkflow, err := client.GetMigrationWorkflow(timeoutCtx, &migrationpb.GetMigrationWorkflowRequest{
				Name: workflow.GetName(),
			})
			if err != nil {
				return nil, fmt.Errorf("polling ended in error: %w", err)
			}
			if polledWorkflow.GetState() == migrationpb.MigrationWorkflow_COMPLETED {
				// polledWorkflow contains the most recent metadata about the workflow, so we return that.
				return polledWorkflow, nil
			}
			// workflow still isn't complete, so sleep briefly before polling again.
			time.Sleep(5 * time.Second)
		}
	}
}

// reportWorkflowStatus prints information about the workflow execution in a more human readable format.
func reportWorkflowStatus(workflow *migrationpb.MigrationWorkflow) {
	fmt.Printf("Migration workflow %s ended in state %s.\n", workflow.GetName(), workflow.GetState().String())
	for k, task := range workflow.GetTasks() {
		fmt.Printf(" - Task %s had id %s", k, task.GetId())
		if task.GetProcessingError() != nil {
			fmt.Printf(" with processing error: %s", task.GetProcessingError().GetReason())
		}
		fmt.Println()
	}
}

// [END bigquerymigration_quickstart]
