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

package bigqueryaccessquickstart

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// [START bigquery_revoke_dataset_access]
func revokeAccessToDataset(w io.Writer, projectID, datasetID, entity string) error {

	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// entity := "user@mydomain.com"

	ctx := context.Background()

	// Creates BigQuery client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Gets dataset handler
	dataset := client.Dataset(datasetID)

	// Gets dataset metadata
	meta, err := dataset.Metadata(ctx)
	if err != nil {
		return err
	}

	// Creates new access entry list by copying the existing and omiting the access entry entity value
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
		meta, err = dataset.Update(ctx, update, meta.ETag)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("any access entry was revoked")
	}

	fmt.Fprintf(w, "Details for Access entries in dataset %v.\n", datasetID)

	for _, access := range meta.Access {
		fmt.Fprintf(w, "Role %s : %s\n", access.Role, access.Entity)
	}

	return nil
}

// [END bigquery_revoke_dataset_access]
