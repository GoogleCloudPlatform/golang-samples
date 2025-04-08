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

// [START bigquery_update_materialized_view]
import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

// updateMaterializedView updates a materialized view by manipulating MV properties.
func updateMaterializedView(projectID, datasetID, viewID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// viewID := "myview"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Retrieve current view metadata.
	viewRef := client.Dataset(datasetID).Table(viewID)
	meta, err := viewRef.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("couldn't retrieve view metadata: %w", err)
	}

	if meta.MaterializedView == nil {
		return fmt.Errorf("provided view %q is not a materialized view", viewID)
	}

	// construct an updated MV definition.
	newMV := &bigquery.MaterializedViewDefinition{
		Query:           meta.MaterializedView.Query,
		EnableRefresh:   true,
		RefreshInterval: meta.MaterializedView.RefreshInterval + time.Minute,
	}

	// Issue the update to alter the view.
	_, err = viewRef.Update(ctx, bigquery.TableMetadataToUpdate{
		MaterializedView: newMV,
	}, meta.ETag)
	if err != nil {
		return fmt.Errorf("Update(): %w", err)
	}
	return nil
}

// [END bigquery_update_materialized_view]
