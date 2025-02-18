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
	"cloud.google.com/go/iam"
)

// [START bigquery_grant_access_to_table_or_view]
func grantAccessToResource(w io.Writer, projectID, datasetID, resourceID string) error {

	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// resourceID := "myresource"

	ctx := context.Background()

	// Creates new client
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Gets resource policy.
	policy, err := client.Dataset(datasetID).Table(resourceID).IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Table.IAM.Policy: %w", err)
	}

	// Find more details about BigQuery Entity Types here:
	// https://pkg.go.dev/cloud.google.com/go/bigquery#EntityType
	//
	// Find more details about IAM Roles here:
	// https://pkg.go.dev/cloud.google.com/go/iam#RoleName

	entityType := bigquery.GroupEmailEntity
	entityID := "example-analyst-group@google.com"
	roleType := iam.Viewer

	// Adds new policy.
	policy.Add(fmt.Sprintf("%v:%s", entityType, entityID), roleType)

	// Updates resource's policy.
	err = client.Dataset(datasetID).Table(resourceID).IAM().SetPolicy(ctx, policy)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Table.IAM.Policy: %w", err)
	}

	// Gets resource policy again expecting the update.
	updatedPolicy, err := client.Dataset(datasetID).Table(resourceID).IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Table.IAM.Policy: %w", err)
	}

	fmt.Fprintf(w, "Details for Access entries in table or view %v.\n", resourceID)

	for _, role := range updatedPolicy.Roles() {
		fmt.Fprintf(w, "Role %s : %s\n", role, policy.Members(role))
	}

	return nil
}

// [END bigquery_grant_access_to_table_or_view]
