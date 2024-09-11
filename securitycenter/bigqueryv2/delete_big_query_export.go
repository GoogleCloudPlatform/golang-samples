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

// [START securitycenter_delete_big_query_export_v2]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	securitycenterpb "cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

func deleteBigQueryExport(w io.Writer, parent string, bigQueryExportId string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}/locations/{location_id}
	//             - folders/{folder_id}/locations/{location_id}
	//             - projects/{project_id}/locations/{location_id}
	// bigQueryExportId := "your-bqexport-id"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.DeleteBigQueryExportRequest{
		Name: fmt.Sprintf("%s/bigQueryExports/%s", parent, bigQueryExportId),
	}

	if err := client.DeleteBigQueryExport(ctx, req); err != nil {
		return fmt.Errorf("failed to delete BigQueryExport: %w", err)
	}
	fmt.Fprintf(w, "BigQueryExport deleted successfully: %s", bigQueryExportId)
	return nil
}

// [END securitycenter_delete_big_query_export_v2]
