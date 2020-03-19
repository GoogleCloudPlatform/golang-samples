// Copyright 2019 Google LLC
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

// [START bigquery_grant_view_access]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// updateViewDelegated demonstrates the setup of an authorized view, which allows access to a view's results
// without the caller having direct access to the underlying source data.
func updateViewDelegated(projectID, srcDatasetID, viewDatasetID, viewID string) error {
	// projectID := "my-project-id"
	// srcDatasetID := "sourcedata"
	// viewDatasetID := "views"
	// viewID := "myview"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	srcDataset := client.Dataset(srcDatasetID)
	viewDataset := client.Dataset(viewDatasetID)
	view := viewDataset.Table(viewID)

	// First, we'll add a group to the ACL for the dataset containing the view.  This will allow users within
	// that group to query the view, but they must have direct access to any tables referenced by the view.
	vMeta, err := viewDataset.Metadata(ctx)
	if err != nil {
		return err
	}
	vUpdateMeta := bigquery.DatasetMetadataToUpdate{
		Access: append(vMeta.Access, &bigquery.AccessEntry{
			Role:       bigquery.ReaderRole,
			EntityType: bigquery.GroupEmailEntity,
			Entity:     "example-analyst-group@google.com",
		}),
	}
	if _, err := viewDataset.Update(ctx, vUpdateMeta, vMeta.ETag); err != nil {
		return err
	}

	// Now, we'll authorize a specific view against a source dataset, delegating access enforcement.
	// Once this has been completed, members of the group previously added to the view dataset's ACL
	// no longer require access to the source dataset to successfully query the view.
	srcMeta, err := srcDataset.Metadata(ctx)
	if err != nil {
		return err
	}
	srcUpdateMeta := bigquery.DatasetMetadataToUpdate{
		Access: append(srcMeta.Access, &bigquery.AccessEntry{
			EntityType: bigquery.ViewEntity,
			View:       view,
		}),
	}
	if _, err := srcDataset.Update(ctx, srcUpdateMeta, srcMeta.ETag); err != nil {
		return err
	}
	return nil
}

// [END bigquery_grant_view_access]
