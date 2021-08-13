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

package spanner

// [START spanner_query_information_schema_database_options]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// queryInformationSchemaDatabaseOptions queries the database options from the
// information schema table.
func queryInformationSchemaDatabaseOptions(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.+)/databases/(.+)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("queryInformationSchemaDatabaseOptions: invalid database id %q", db)
	}
	databaseID := matches[2]

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT OPTION_NAME, OPTION_VALUE
	                                FROM INFORMATION_SCHEMA.DATABASE_OPTIONS 
                                    WHERE OPTION_NAME = 'default_leader'`}
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
		var option_name, option_value string
		if err := row.Columns(&option_name, &option_value); err != nil {
			return err
		}
		fmt.Fprintf(w, "The result of the query to get %s for %s is %s", option_name, databaseID, option_value)
	}
}

// [END spanner_query_information_schema_database_options]
