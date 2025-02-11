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

package viewdatasetaccesspolicy

// [START bigquery_view_dataset_access_policies]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

func viewDatasetAccessPolicies(w io.Writer, projectID, datasetName string) error {
	ctx := context.Background()

	// Creates new client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Creates handle for managing dataset
	dataset := client.Dataset(datasetName)

	// Gets dataset's metadata
	metaData, err := dataset.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("dataset.Metadata: %v", err)
	}

	// Iterate over access permissions
	for _, val := range metaData.Access {
		fmt.Fprintf(w, "Role %s : %s\n", val.Role, val.Entity)
	}

	return nil
}

// [START bigquery_view_dataset_access_policies]
