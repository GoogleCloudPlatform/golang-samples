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

// [START spanner_query_with_timestamp_parameter]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithTimestampParameter(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	var exampleTimestamp = time.Now()
	stmt := spanner.Statement{
		SQL: `SELECT VenueId, VenueName, LastUpdateTime FROM Venues
		WHERE LastUpdateTime <= @lastUpdateTime`,
		Params: map[string]interface{}{
			"lastUpdateTime": exampleTimestamp,
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
		var lastUpdateTime time.Time
		if err := row.Columns(&venueID, &venueName, &lastUpdateTime); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s %s\n", venueID, venueName, lastUpdateTime)
	}
}

// [END spanner_query_with_timestamp_parameter]
