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

// [START bigquery_update_view_query]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// updateView demonstrates updating the query metadata that defines a logical view.
func updateView(projectID, datasetID, viewID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// viewID := "myview"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	view := client.Dataset(datasetID).Table(viewID)
	meta, err := view.Metadata(ctx)
	if err != nil {
		return err
	}

	newMeta := bigquery.TableMetadataToUpdate{
		// This example updates a view into the shakespeare dataset to exclude works named after kings.
		ViewQuery: "SELECT word, word_count, corpus, corpus_date FROM `bigquery-public-data.samples.shakespeare` WHERE corpus NOT LIKE '%king%'",
	}

	if _, err := view.Update(ctx, newMeta, meta.ETag); err != nil {
		return err
	}
	return nil
}

// [END bigquery_update_view_query]
