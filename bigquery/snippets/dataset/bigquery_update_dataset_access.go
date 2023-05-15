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

// [START bigquery_update_dataset_access]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// updateDatasetAccessControl demonstrates how the access control policy of a dataset
// can be amended by adding an additional entry corresponding to a specific user identity.
func updateDatasetAccessControl(projectID, datasetID string) error {
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
	// Append a new access control entry to the existing access list.
	update := bigquery.DatasetMetadataToUpdate{
		Access: append(meta.Access, &bigquery.AccessEntry{
			Role:       bigquery.ReaderRole,
			EntityType: bigquery.UserEmailEntity,
			Entity:     "sample.bigquery.dev@gmail.com"},
		),
	}

	// Leverage the ETag for the update to assert there's been no modifications to the
	// dataset since the metadata was originally read.
	if _, err := ds.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	return nil
}

// [END bigquery_update_dataset_access]
