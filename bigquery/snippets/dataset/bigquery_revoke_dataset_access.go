// Copyright 2021 Google LLC
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

// [START bigquery_revoke_dataset_access]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// revokeDatasetAccess updates the access control on a dataset to remove all
// access entries that reference a specific entity.
func revokeDatasetAccess(projectID, datasetID, entity string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// entity := "user@mydomain.com"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}

	var newAccessList []*bigquery.AccessEntry
	for _, entry := range meta.Access {
		if entry.Entity != entity {
			newAccessList = append(newAccessList, entry)
		}
	}

	// Only proceed with update if something in the access list was removed.
	// Additionally, we use the ETag from the initial metadata to ensure no
	// other changes were made to the access list in the interim.
	if len(newAccessList) < len(meta.Access) {

		update := bigquery.DatasetMetadataToUpdate{
			Access: newAccessList,
		}
		if _, err := ds.Update(ctx, update, meta.ETag); err != nil {
			return err
		}
	}
	return nil
}

// [END bigquery_revoke_dataset_access]
