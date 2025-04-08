// Copyright 2021 Google LLC
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

package routine

// [START bigquery_update_routine]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// updateRoutine demonstrates updating an existing BigQuery UDF using the routine API.
func updateRoutine(projectID, datasetID, routineID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// routineID := "myroutineid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	routineRef := client.Dataset(datasetID).Routine(routineID)

	// fetch existing metadata
	meta, err := routineRef.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("couldn't retrieve metadata: %w", err)
	}

	// Due to a limitation in the backend, supply all the properties for update.
	update := &bigquery.RoutineMetadataToUpdate{
		Type:        meta.Type,
		Language:    meta.Language,
		Arguments:   meta.Arguments,
		Description: meta.Description,
		ReturnType:  meta.ReturnType,
		Body:        "x * 4",
	}

	if _, err := routineRef.Update(ctx, update, meta.ETag); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

// [END bigquery_update_routine]
