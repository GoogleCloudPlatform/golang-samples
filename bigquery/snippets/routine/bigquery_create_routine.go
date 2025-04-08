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

// [START bigquery_create_routine]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createRoutine demonstrates creating a new BigQuery UDF using the routine API.
func createRoutine(projectID, datasetID, routineID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// routineID := "myroutineid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	metaData := &bigquery.RoutineMetadata{
		Type:     "SCALAR_FUNCTION",
		Language: "SQL",
		Body:     "x * 3",
		Arguments: []*bigquery.RoutineArgument{
			{Name: "x", DataType: &bigquery.StandardSQLDataType{TypeKind: "INT64"}},
		},
	}

	routineRef := client.Dataset(datasetID).Routine(routineID)
	if err := routineRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_routine]
