// Copyright 2025 Google LLC
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

package table

// [START bigquery_get_table_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"google.golang.org/protobuf/encoding/protojson"
)

// getTable demonstrates fetching table information.
func getTable(client *apiv2_client.Client, w io.Writer, projectID, datasetID, tableID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()

	req := &bigquerypb.GetTableRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		TableId:   tableID,

		// BigQuery can return only a subset of the table's information, based on
		// your needs.  In particular, for large tables it is recommended to only
		// fetch storage information is necessary.  In this example, we'll only get
		// BASIC information about the table.
		View: bigquerypb.GetTableRequest_BASIC,
	}

	resp, err := client.GetTable(ctx, req)
	if err != nil {
		return fmt.Errorf("GetTable: %w", err)
	}

	// Print some of the information about the model to the provided writer.
	fmt.Fprintf(w, "Table %q has description %q\n",
		resp.GetTableReference().GetTableId(),
		resp.GetDescription())

	// Print information about the top-level schema of the table.
	if schema := resp.GetSchema(); schema != nil {
		fmt.Fprintf(w, "Table %q has %d top-level fields\n",
			resp.GetTableReference().GetTableId(),
			len(schema.GetFields()))
	}
	// Alternately, use the protojson package to print a more complete representation
	// of the table using a basic JSON mapping.
	fmt.Fprintf(w, "Table JSON representation:\n%s\n", protojson.Format(resp))
	return nil
}

// [END bigquery_get_table_preview]
