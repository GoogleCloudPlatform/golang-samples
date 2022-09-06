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

package table

// [START bigquery_create_view]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createView demonstrates creation of a BigQuery logical view.
func createView(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// tableID := "mytableid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	meta := &bigquery.TableMetadata{
		// This example shows how to create a view of the shakespeare sample dataset, which
		// provides word frequency information.  This view restricts the results to only contain
		// results for works that contain the "king" in the title, e.g. King Lear, King Henry V, etc.
		ViewQuery: "SELECT word, word_count, corpus, corpus_date FROM `bigquery-public-data.samples.shakespeare` WHERE corpus LIKE '%king%'",
	}
	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, meta); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_view]
