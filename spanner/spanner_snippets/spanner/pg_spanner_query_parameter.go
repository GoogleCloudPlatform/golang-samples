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

// [START spanner_postgresql_query_parameter]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// pgQueryParameter shows how to execute a query with parameters on a Spanner
// PostgreSQL database. The PostgreSQL dialect uses positional parameters, as
// opposed to the named parameters of Cloud Spanner.
func pgQueryParameter(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT SingerId, FirstName, LastName FROM Singers
			WHERE LastName = $1`,
		Params: map[string]interface{}{
			"p1": "Garcia",
		},
	}
	type Singers struct {
		SingerID            int64
		FirstName, LastName string
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var val Singers
		if err := row.ToStruct(&val); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s %s\n", val.SingerID, val.FirstName, val.LastName)
	}
}

// [END spanner_postgresql_query_parameter]
