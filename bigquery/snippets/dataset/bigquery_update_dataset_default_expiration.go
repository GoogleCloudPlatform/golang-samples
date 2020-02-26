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

// [START bigquery_update_dataset_expiration]
import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

// updateDatasetDefaultExpiration demonstrats setting the default expiration of a dataset
// to a specific retention period.
func updateDatasetDefaultExpiration(projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.DatasetMetadataToUpdate{
		DefaultTableExpiration: 24 * time.Hour,
	}
	if _, err := client.Dataset(datasetID).Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	return nil
}

// [END bigquery_update_dataset_expiration]
