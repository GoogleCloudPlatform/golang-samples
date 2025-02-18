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

// [START bigquery_view_dataset_access_policy]
func viewDatasetAccessPolicies(w io.Writer, projectID, datasetID string) error {

	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// datasetID := "mydataset"

	ctx := context.Background()

	// Creates new client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Gets dataset's metadata
	metaData, err := client.Dataset(datasetID).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.Client.Dataset.Metadata: %w", err)
	}

	fmt.Fprintf(w, "Details for Access entries in dataset %v.\n", datasetID)
	// Iterate over access permissions
	for _, access := range metaData.Access {
		fmt.Fprintf(w, "Role %s : %s\n", access.Role, access.Entity)
	}

	return nil
}

// [END bigquery_view_dataset_access_policy]
