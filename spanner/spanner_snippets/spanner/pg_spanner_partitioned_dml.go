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

// [START spanner_postgresql_partitioned_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

// pgPartitionedDml shows how to execute Partitioned DML on a Spanner PostgreSQL database.
// See also https://cloud.google.com/spanner/docs/dml-partitioned.
func pgPartitionedDml(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Spanner PostgreSQL has the same transaction limits as normal Spanner. This includes a
	// maximum of 20,000 mutations in a single read/write transaction. Large update operations can
	// be executed using Partitioned DML. This is also supported on Spanner PostgreSQL.
	// See https://cloud.google.com/spanner/docs/dml-partitioned for more information.
	deletedCount, err := client.PartitionedUpdate(ctx, spanner.Statement{
		SQL: "DELETE FROM users WHERE active=false",
	})
	if err != nil {
		return err
	}
	// The returned update count is the lower bound of the number of records that was deleted.
	fmt.Fprintf(w, "Deleted at least %d inactive users\n", deletedCount)

	return nil
}

// [END spanner_postgresql_partitioned_dml]
