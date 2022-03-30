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

// [START spanner_postgresql_cast_data_type]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// pgCastDataType shows how to cast values from one data type to another in a
// Spanner PostgreSQL SQL statement.
func pgCastDataType(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// The `::` cast operator can be used to cast from one data type to another.
	query := `select 1::varchar as str, '2'::int as int, 3::decimal as dec,
				'4'::bytea as bytes, 5::float as float, 'true'::bool as bool,
				'2021-11-03T09:35:01UTC'::timestamptz as timestamp`
	stmt := spanner.Statement{SQL: query}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	// Row represents a row returned from above cast query.
	type Row struct {
		Str       string
		Int       int64
		Dec       spanner.PGNumeric
		Bytes     []byte
		Float     float64
		Bool      bool
		Timestamp time.Time
	}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var val Row
		if err := row.ToStruct(&val); err != nil {
			return err
		}
		fmt.Fprintf(w, "String: %s\n", val.Str)
		fmt.Fprintf(w, "Int: %d\n", val.Int)
		fmt.Fprintf(w, "Decimal: %s\n", val.Dec)
		fmt.Fprintf(w, "Bytes: %s\n", val.Bytes)
		fmt.Fprintf(w, "Float: %f\n", val.Float)
		fmt.Fprintf(w, "Bool: %v\n", val.Bool)
		fmt.Fprintf(w, "Timestamp: %s\n", val.Timestamp)
	}
}

// [END spanner_postgresql_cast_data_type]
