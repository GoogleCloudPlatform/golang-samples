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

// [START bigquery_grant_access_to_dataset]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// grantAccessToDataset creates a new ACL conceding the READER role to the group "example-analyst-group@google.com"
// For more information on the types of ACLs available see:
// https://cloud.google.com/storage/docs/access-control/lists
func grantAccessToDataset(w io.Writer, projectID, datasetID string) error {
	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// datasetID := "mydataset"

	ctx := context.Background()

	// Create BigQuery handler.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Create dataset handler
	dataset := client.Dataset(datasetID)

	// Get metadata
	meta, err := dataset.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Metadata: %w", err)
	}

	// Find more details about BigQuery Entity Types here:
	// https://pkg.go.dev/cloud.google.com/go/bigquery#EntityType
	//
	// Find more details about BigQuery Access Roles here:
	// https://pkg.go.dev/cloud.google.com/go/bigquery#AccessRole

	entityType := bigquery.GroupEmailEntity
	entityID := "example-analyst-group@google.com"
	roleType := bigquery.ReaderRole

	// Append a new access control entry to the existing access list.
	update := bigquery.DatasetMetadataToUpdate{
		Access: append(meta.Access, &bigquery.AccessEntry{
			Role:       roleType,
			EntityType: entityType,
			Entity:     entityID,
		}),
	}

	// Leverage the ETag for the update to assert there's been no modifications to the
	// dataset since the metadata was originally read.
	meta, err = dataset.Update(ctx, update, meta.ETag)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Details for Access entries in dataset %v.\n", datasetID)
	for _, access := range meta.Access {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Role: %s\n", access.Role)
		fmt.Fprintf(w, "Entities: %v\n", access.Entity)
	}

	return nil
}

// [END bigquery_grant_access_to_dataset]
