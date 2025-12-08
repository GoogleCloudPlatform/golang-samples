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

// [START bigquery_revoke_access_to_table_or_view]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
)

// revokeTableOrViewAccessPolicies creates a new ACL removing the VIEWER role to group "example-analyst-group@google.com"
// For more information on the types of ACLs available see:
// https://cloud.google.com/storage/docs/access-control/lists
func revokeTableOrViewAccessPolicies(w io.Writer, projectID, datasetID, resourceID string) error {
	// Resource can be a table or a view
	//
	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// resourceID := "myresource"

	ctx := context.Background()

	// Create new client
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Get resource policy.
	policy, err := client.Dataset(datasetID).Table(resourceID).IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Table.IAM.Policy: %w", err)
	}

	// Find more details about IAM Roles here:
	// https://pkg.go.dev/cloud.google.com/go/iam#RoleName
	entityID := "example-analyst-group@google.com"
	roleType := iam.Viewer

	// Revoke policy access.
	policy.Remove(fmt.Sprintf("group:%s", entityID), roleType)

	// Update resource's policy.
	err = client.Dataset(datasetID).Table(resourceID).IAM().SetPolicy(ctx, policy)
	if err != nil {
		return fmt.Errorf("bigquery.Dataset.Table.IAM.Policy: %w", err)
	}

	// Get resource policy again expecting the update.
	policy, err = client.Dataset(datasetID).Table(resourceID).IAM().Policy(ctx)
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

// [END bigquery_revoke_access_to_table_or_view]
