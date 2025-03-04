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

// [START bigquery_view_table_or_view_access_policy]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// viewTableOrViewAccessPolicies retrieves the ACL for the given resource
// For more information on the types of ACLs available see:
// https://cloud.google.com/storage/docs/access-control/lists
func viewTableOrViewAccessPolicies(w io.Writer, projectID, datasetID, resourceID string) error {
	// Resource can be a table or a view
	//
	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// datasetID := "my-dataset-id"
	// resourceID := "my-resource-id"

	ctx := context.Background()

	// Create new client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Get resource's policy access.
	policy, err := client.Dataset(datasetID).Table(resourceID).IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Table.IAM.Policy: %w", err)
	}

	fmt.Fprintf(w, "Details for Access entries in table or view %v.\n", resourceID)

	for _, role := range policy.Roles() {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Role: %s\n", role)
		fmt.Fprintf(w, "Entities: %v\n", policy.Members(role))
	}

	return nil
}

// [END bigquery_view_table_or_view_access_policy]
