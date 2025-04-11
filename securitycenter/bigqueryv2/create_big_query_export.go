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

package bigqueryv2

// [START securitycenter_create_big_query_export_v2]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	securitycenterpb "cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

// Create export configuration to export findings to a BigQuery dataset.
// Optionally specify filter to export certain findings only.
func createBigQueryExport(w io.Writer, parent string, bigQueryExportId string, projectId string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}/locations/{location_id}
	//             - folders/{folder_id}/locations/{location_id}
	//             - projects/{project_id}/locations/{location_id}
	// bigQueryExportId := "random-bqexport-id-" + uuid.New().String()
	// bigQueryDataSetName := fmt.Sprintf("projects/%s/datasets/%s", "your-google-cloud-project-id", "your-big-query-dataset-id")

	// Hard-coded dataset name
	bigQueryDataSetName := fmt.Sprintf("projects/%s/datasets/sampledataset", projectId)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	bigQueryExport := &securitycenterpb.BigQueryExport{
		Description: "BigQueryExport that receives all HIGH severity Findings",
		Filter:      "severity=\"HIGH\"",
		Dataset:     bigQueryDataSetName,
	}

	req := &securitycenterpb.CreateBigQueryExportRequest{
		Parent:           parent,
		BigQueryExport:   bigQueryExport,
		BigQueryExportId: bigQueryExportId,
	}

	bigqueryconfig, err := client.CreateBigQueryExport(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to create BigQueryConfig: %w", err)
	}
	fmt.Fprintf(w, "BigQueryConfig Name: %s ", bigqueryconfig.Name)
	return nil
}

// [END securitycenter_create_big_query_export_v2]
