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

func queryInformationSchemaDatabaseOptions(w io.Writer, db string) error {
	matches := regexp.MustCompile("^(.+)/databases/(.+)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("createDatabaseWithCustomerManagedEncryptionKey: invalid database id %q", db)
	}
	databaseId := matches[2]

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, MarketingBudget FROM Albums`}
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
		var OPTION_NAME, OPTION_VALUE string
		if err := row.Columns(&OPTION_NAME, &OPTION_VALUE); err != nil {
			return err
		}
		fmt.Fprintf(w, "The result of the query to get %s for %s is %s", OPTION_NAME, databaseId, OPTION_VALUE)
	}
}

// [END spanner_dml_batch_update]
