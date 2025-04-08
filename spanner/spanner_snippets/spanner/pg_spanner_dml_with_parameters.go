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

// [START spanner_postgresql_dml_with_parameters]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

// pgDmlWithParameters shows how to execute a DML statement with parameters
// on a Spanner PostgreSQL database. The PostgreSQL dialect uses positional
// parameters, as opposed to the named parameters that are used by Cloud Spanner.
func pgDmlWithParameters(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `INSERT INTO Singers (SingerId, FirstName, LastName) 
              VALUES ($1, $2, $3), ($4, $5, $6)`,
		// Use 'p1' to bind to the parameter with index 1.
		Params: map[string]interface{}{
			"p1": 1, "p2": "Alice", "p3": "Henderson",
			"p4": 2, "p5": "Bruce", "p6": "Allison",
		},
	}
	var updateCount int64
	if _, err := client.ReadWriteTransaction(context.Background(), func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		updateCount, err = transaction.Update(ctx, stmt)
		return err
	}); err != nil {
		return err
	}
	fmt.Fprintf(w, "Inserted %d singers\n", updateCount)
	return nil
}

// [END spanner_postgresql_dml_with_parameters]
