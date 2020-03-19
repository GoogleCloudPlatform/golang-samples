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

package job

// [START bigquery_copy_table_multiple_source]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// copyMultiTable demonstrates using a copy job to copy multiple source tables into a single destination table.
func copyMultiTable(projectID, srcDatasetID string, srcTableIDs []string, dstDatasetID, dstTableID string) error {
	// projectID := "my-project-id"
	// srcDatasetID := "sourcedataset"
	// srcTableIDs := []string{"table1","table2"}
	// dstDatasetID = "destinationdataset"
	// dstTableID = "destinationtable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	srcDataset := client.Dataset(srcDatasetID)
	dstDataset := client.Dataset(dstDatasetID)
	var tableRefs []*bigquery.Table
	for _, v := range srcTableIDs {
		tableRefs = append(tableRefs, srcDataset.Table(v))
	}
	copier := dstDataset.Table(dstTableID).CopierFrom(tableRefs...)
	copier.WriteDisposition = bigquery.WriteTruncate
	job, err := copier.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// [END bigquery_copy_table_multiple_source]
