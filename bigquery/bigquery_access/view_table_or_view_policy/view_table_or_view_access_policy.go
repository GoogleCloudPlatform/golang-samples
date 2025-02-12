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

package viewtableorviewpolicy

// [START bigquery_view_table_access_policies]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

func viewTableAccessPolicies(w io.Writer, projectID, datasetName, resourceName string) error {
	ctx := context.Background()

	// Creates new client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Creates handle for managing dataset's table.
	table := client.Dataset(datasetName).Table(resourceName)

	// Gets table's policy access.
	policy, err := table.IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("table.IAM.Policy %w", err)
	}

	fmt.Fprintf(w, "Details for Access entries in table or view %v.\n", resourceName)

	for _, role := range policy.Roles() {
		fmt.Fprintf(w, "Role %s : %s\n", role, policy.Members(role))
	}

	return nil
}

// [START bigquery_view_table_access_policies]
