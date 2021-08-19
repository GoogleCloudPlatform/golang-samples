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

// [START spanner_query_with_json_parameter]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithJsonParameter(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("addJsonColumn: invalid database id %s", db)
	}
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT VenueId, VenueDetails FROM Venues WHERE JSON_VALUE(VenueDetails, '$.rating') = JSON_VALUE(@details, '$.rating')`,
		Params: map[string]interface{}{
			"details": map[string]interface{}{"rating": 9},
		},
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
		var venueID int64
		var venueDetails string
		if err := row.Columns(&venueID, &venueDetails); err != nil {
			return err
		}
		fmt.Fprintf(w, "The venue details for venue id %v is %v\n", venueID, venueDetails)
	}
}

// [END spanner_query_with_json_parameter]
