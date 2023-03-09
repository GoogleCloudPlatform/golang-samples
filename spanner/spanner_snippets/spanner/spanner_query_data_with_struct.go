// Copyright 2020 Google LLC
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

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithStruct(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// [START spanner_create_struct_with_data]

	type name struct {
		FirstName string
		LastName  string
	}
	var singerInfo = name{"Elena", "Campbell"}

	// [END spanner_create_struct_with_data]

	// [START spanner_query_data_with_struct]

	stmt := spanner.Statement{
		SQL: `SELECT SingerId FROM SINGERS
				WHERE (FirstName, LastName) = @singerinfo`,
		Params: map[string]interface{}{"singerinfo": singerInfo},
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
		var singerID int64
		if err := row.Columns(&singerID); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d\n", singerID)
	}

	// [END spanner_query_data_with_struct]
}
