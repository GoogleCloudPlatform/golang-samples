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

package revokeiampolicytableorview

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
)

func revokeTableOrViewAccessPolicies(w io.Writer, projectID, datasetID, resourceID string) error {
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

	// Revokes policy access.
	analystEmail := "example-analyst-group@google.com"
	policy.Remove(fmt.Sprintf("group:%s", analystEmail), iam.Viewer)

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
