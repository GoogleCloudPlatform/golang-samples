// Copyright 2023 Google LLC
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

package risk

// [START dlp_k_anonymity_with_entity_id]
import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// Uses the Data Loss Prevention API to compute the k-anonymity of a
// column set in a Google BigQuery table.
func calculateKAnonymityWithEntityId(w io.Writer, projectID, datasetId, tableId string, columnNames ...string) error {
	// projectID := "your-project-id"
	// datasetId := "your-bigquery-dataset-id"
	// tableId := "your-bigquery-table-id"
	// columnNames := "age" "job_title"

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the BigQuery table to analyze
	bigQueryTable := &dlppb.BigQueryTable{
		ProjectId: "bigquery-public-data",
		DatasetId: "samples",
		TableId:   "wikipedia",
	}

	// Configure the privacy metric for the job
	// Build the QuasiID slice.
	var q []*dlppb.FieldId
	for _, c := range columnNames {
		q = append(q, &dlppb.FieldId{Name: c})
	}

	entityId := &dlppb.EntityId{
		Field: &dlppb.FieldId{
			Name: "id",
		},
	}

	kAnonymityConfig := &dlppb.PrivacyMetric_KAnonymityConfig{
		QuasiIds: q,
		EntityId: entityId,
	}

	privacyMetric := &dlppb.PrivacyMetric{
		Type: &dlppb.PrivacyMetric_KAnonymityConfig_{
			KAnonymityConfig: kAnonymityConfig,
		},
	}

	// Specify the bigquery table to store the findings.
	// The "test_results" table in the given BigQuery dataset will be created if it doesn't
	// already exist.
	outputbigQueryTable := &dlppb.BigQueryTable{
		ProjectId: projectID,
		DatasetId: datasetId,
		TableId:   tableId,
	}

	// Create action to publish job status notifications to BigQuery table.
	outputStorageConfig := &dlppb.OutputStorageConfig{
		Type: &dlppb.OutputStorageConfig_Table{
			Table: outputbigQueryTable,
		},
	}

	findings := &dlppb.Action_SaveFindings{
		OutputConfig: outputStorageConfig,
	}

	action := &dlppb.Action{
		Action: &dlppb.Action_SaveFindings_{
			SaveFindings: findings,
		},
	}

	// Configure the risk analysis job to perform
	riskAnalysisJobConfig := &dlppb.RiskAnalysisJobConfig{
		PrivacyMetric: privacyMetric,
		SourceTable:   bigQueryTable,
		Actions: []*dlppb.Action{
			action,
		},
	}

	// Build the request to be sent by the client
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_RiskJob{
			RiskJob: riskAnalysisJobConfig,
		},
	}

	// Send the request to the API using the client
	dlpJob, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Created job: %v\n", dlpJob.GetName())

	// Build a request to get the completed job
	getDlpJobReq := &dlppb.GetDlpJobRequest{
		Name: dlpJob.Name,
	}

	timeout := 15 * time.Minute
	startTime := time.Now()

	var completedJob *dlppb.DlpJob

	// Wait for job completion
	for time.Since(startTime) <= timeout {
		completedJob, err = client.GetDlpJob(ctx, getDlpJobReq)
		if err != nil {
			return err
		}

		if completedJob.GetState() == dlppb.DlpJob_DONE {
			break
		}

		time.Sleep(30 * time.Second)

	}

	if completedJob.GetState() != dlppb.DlpJob_DONE {
		fmt.Println("Job did not complete within 15 minutes.")
	}

	// Retrieve completed job status
	fmt.Fprintf(w, "Job status: %v", completedJob.State)
	fmt.Fprintf(w, "Job name: %v", dlpJob.Name)

	// Get the result and parse through and process the information
	kanonymityResult := completedJob.GetRiskDetails().GetKAnonymityResult()

	for _, result := range kanonymityResult.GetEquivalenceClassHistogramBuckets() {
		fmt.Fprintf(w, "Bucket size range: [%d, %d]\n", result.GetEquivalenceClassSizeLowerBound(), result.GetEquivalenceClassSizeLowerBound())

		for _, bucket := range result.GetBucketValues() {
			quasiIdValues := []string{}
			for _, v := range bucket.GetQuasiIdsValues() {
				quasiIdValues = append(quasiIdValues, v.GetStringValue())
			}
			fmt.Fprintf(w, "\tQuasi-ID values: %s", strings.Join(quasiIdValues, ","))
			fmt.Fprintf(w, "\tClass size: %d", bucket.EquivalenceClassSize)
		}
	}

	return nil

}

// [END dlp_k_anonymity_with_entity_id]
