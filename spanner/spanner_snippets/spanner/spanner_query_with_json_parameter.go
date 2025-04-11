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

// [START spanner_query_with_json_parameter]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// queryWithJsonParameter queries data on the JSON type column of the database
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

	type VenueDetails struct {
		Name   spanner.NullString   `json:"name"`
		Rating spanner.NullFloat64  `json:"rating"`
		Open   interface{}          `json:"open"`
		Tags   []spanner.NullString `json:"tags"`
	}

	stmt := spanner.Statement{
		SQL: `SELECT VenueId, VenueDetails FROM Venues WHERE JSON_VALUE(VenueDetails, '$.rating') = JSON_VALUE(@details, '$.rating')`,
		Params: map[string]interface{}{
			"details": spanner.NullJSON{Value: VenueDetails{
				Rating: spanner.NullFloat64{Float64: 9, Valid: true},
			}, Valid: true},
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
		var venueDetails spanner.NullJSON
		if err := row.Columns(&venueID, &venueDetails); err != nil {
			return err
		}
		fmt.Fprintf(w, "The venue details for venue id %v is %v\n", venueID, venueDetails)
	}
}

// [END spanner_query_with_json_parameter]
