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

// [START bigquery_create_table_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"

	"github.com/googleapis/gax-go/v2/apierror"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

// createTable demonstrates creation of a new table with a predefined schema into an existing dataset.
func createTable(client *apiv2_client.Client, w io.Writer, projectID, datasetID, tableID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()

	// Define a very simple schema for the table.
	schema := &bigquerypb.TableSchema{
		Fields: []*bigquerypb.TableFieldSchema{
			{
				Name: "name",
				Type: "STRING",
				Mode: "REQUIRED",
			},
			{
				Name: "age",
				Type: "INTEGER",
			},
			{
				Name: "weight",
				Type: "FLOAT",
			},
			{
				Name: "is_magic",
				Type: "BOOLEAN",
			},
		},
	}

	// Construct a request, populating some of the available configuration
	// settings.
	req := &bigquerypb.InsertTableRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		Table: &bigquerypb.Table{
			TableReference: &bigquerypb.TableReference{
				ProjectId: projectID,
				DatasetId: datasetID,
				TableId:   tableID,
			},
			Schema: schema,
		},
	}
	resp, err := client.InsertTable(ctx, req)
	if err != nil {
		// Examine the error structure more deeply.
		if apierr, ok := apierror.FromError(err); ok {
			if status := apierr.GRPCStatus(); status.Code() == codes.AlreadyExists {
				// The error was due to the table already existing.  For this sample
				// we don't consider that a failure, so return nil.
				return nil
			}
		}
		return fmt.Errorf("InsertTable: %w", err)
	}
	// Print the JSON representation of the response to the provided writer.
	fmt.Fprintf(w, "Response from insert: %s", protojson.Format(resp))
	return nil
}

// [END bigquery_create_table_preview]
