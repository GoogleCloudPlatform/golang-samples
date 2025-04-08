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

// [START bigquery_create_routine_ddl]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createRoutineDDL demonstrates creating a new BigQuery UDF using a DDL query.
func createRoutineDDL(projectID, datasetID, routineID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// routineID := "myroutineid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	routineName, err := client.Dataset(datasetID).Routine(routineID).Identifier(bigquery.StandardSQLID)
	if err != nil {
		return fmt.Errorf("couldn't generate identifier: %w", err)
	}

	sql := fmt.Sprintf(`CREATE FUNCTION %s(
        	arr ARRAY<STRUCT<name STRING, val INT64>>
    		) AS (
        	(SELECT SUM(IF(elem.name = "foo",elem.val,null)) FROM UNNEST(arr) AS elem)
    		)`, routineName)

	job, err := client.Query(sql).Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_routine_ddl]
