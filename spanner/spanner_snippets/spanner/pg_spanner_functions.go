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

// [START spanner_postgresql_functions]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// pgFunctions shows how to call a server side function on a Spanner PostgreSQL database.
func pgFunctions(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Use the PostgreSQL `to_timestamp` function to convert a number of seconds since epoch to a
	// timestamp. 1284352323 seconds = Monday, September 13, 2010 4:32:03 AM.
	iter := client.Single().Query(ctx, spanner.Statement{
		SQL: "SELECT to_timestamp(1284352323) AS t",
	})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var ts time.Time
		if err := row.Columns(&ts); err != nil {
			return err
		}
		fmt.Fprintf(w, "1284352323 seconds after epoch is %s\n", ts)
	}
}

// [END spanner_postgresql_functions]
