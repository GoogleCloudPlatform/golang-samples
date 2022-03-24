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

package spanner

// [START spanner_postgresql_batch_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

// pgBatchDml shows how to execute a batch of DML statements on a Spanner PostgreSQL database.
func pgBatchDml(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Spanner PostgreSQL supports BatchDML statements. This will batch multiple DML statements
	// into one request, which reduces the number of round trips that is needed for multiple DML
	// statements.
	insertSQL := `INSERT INTO Singers (SingerId, FirstName, LastName) 
            VALUES ($1, $2, $3)`
	stmts := []spanner.Statement{
		{
			SQL: insertSQL,
			Params: map[string]interface{}{
				"p1": 1,
				"p2": "Alice",
				"p3": "Henderson",
			},
		},
		{
			SQL: insertSQL,
			Params: map[string]interface{}{
				"p1": 2,
				"p2": "Bruce",
				"p3": "Allison",
			},
		},
	}
	var updateCounts []int64
	_, err = client.ReadWriteTransaction(context.Background(), func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		updateCounts, err = transaction.BatchUpdate(ctx, stmts)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted %v singers\n", updateCounts)
	return nil
}

// [END spanner_postgresql_batch_dml]
