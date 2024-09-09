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

// [START securitycenter_update_big_query_export_v2]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	securitycenterpb "cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Updates an existing BigQuery export.
func updateBigQueryExport(w io.Writer, parent string, bigQueryExportId string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}/locations/{location_id}
	//             - folders/{folder_id}/locations/{location_id}
	//             - projects/{project_id}/locations/{location_id}
	// bigQueryExportId := "your-bq-export-id"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	bigQueryExport := &securitycenterpb.BigQueryExport{
		Name:        fmt.Sprintf("%s/bigQueryExports/%s", parent, bigQueryExportId),
		Description: "BigQueryExport that receives all HIGH severity Findings, updated description",
	}

	req := &securitycenterpb.UpdateBigQueryExportRequest{
		BigQueryExport: bigQueryExport,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"description",
			},
		},
	}

	bigqueryconfig, err := client.UpdateBigQueryExport(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to update BigQueryConfig: %w", err)
	}
	fmt.Fprintf(w, "BigQueryConfig Name: %s ", bigqueryconfig.Name)
	return nil
}

// [END securitycenter_update_big_query_export_v2]
