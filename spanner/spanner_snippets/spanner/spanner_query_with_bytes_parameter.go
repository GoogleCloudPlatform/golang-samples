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

// [START spanner_query_with_bytes_parameter]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithBytes(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	var exampleBytes = []byte("Hello World 1")
	stmt := spanner.Statement{
		SQL: `SELECT VenueId, VenueName FROM Venues
            	WHERE VenueInfo = @venueInfo`,
		Params: map[string]interface{}{
			"venueInfo": exampleBytes,
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
		var venueName string
		if err := row.Columns(&venueID, &venueName); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s\n", venueID, venueName)
	}
}

// [END spanner_query_with_bytes_parameter]
