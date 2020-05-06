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

package dataset

// [START bigquery_delete_dataset]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// deleteDataset demonstrates the deletion of an empty dataset.
func deleteDataset(projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// To recursively delete a dataset and contents, use DeleteWithContents.
	if err := client.Dataset(datasetID).Delete(ctx); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}
	return nil
}

// [END bigquery_delete_dataset]
