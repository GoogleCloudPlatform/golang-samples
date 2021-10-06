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

// [START spanner_query_data_with_index]

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryUsingIndex(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT AlbumId, AlbumTitle, MarketingBudget
			FROM Albums@{FORCE_INDEX=AlbumsByAlbumTitle}
			WHERE AlbumTitle >= @start_title AND AlbumTitle < @end_title`,
		Params: map[string]interface{}{
			"start_title": "Aardvark",
			"end_title":   "Goo",
		},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var albumID int64
		var marketingBudget spanner.NullInt64
		var albumTitle string
		if err := row.ColumnByName("AlbumId", &albumID); err != nil {
			return err
		}
		if err := row.ColumnByName("AlbumTitle", &albumTitle); err != nil {
			return err
		}
		if err := row.ColumnByName("MarketingBudget", &marketingBudget); err != nil {
			return err
		}
		budget := "NULL"
		if marketingBudget.Valid {
			budget = strconv.FormatInt(marketingBudget.Int64, 10)
		}
		fmt.Fprintf(w, "%d %s %s\n", albumID, albumTitle, budget)
	}
	return nil
}

// [END spanner_query_data_with_index]
