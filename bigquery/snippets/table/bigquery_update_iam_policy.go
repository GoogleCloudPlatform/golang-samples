// Copyright 2022 Google LLC
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

package table

// [START bigquery_update_iam_policy]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
)

// updateIAMPolicy demonstrates updating the ACL on a BigQuery table by adding a new entry to
// an existing policy.
func updateIAMPolicy(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	table := client.Dataset(datasetID).Table(tableID)
	policy, err := table.IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("failed to get table policy: %w", err)
	}
	newRole := iam.RoleName("roles/bigquery.dataViewer")
	newMember := "allAuthenticatedUsers"
	policy.Add(newMember, newRole)

	if err := table.IAM().SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("failed to set new table policy: %w", err)
	}

	return nil
}

// [END bigquery_update_iam_policy]
