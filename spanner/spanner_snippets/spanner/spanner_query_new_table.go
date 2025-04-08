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
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryNewTable(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT SingerId, VenueId, EventDate, Revenue, LastUpdateTime FROM Performances
				ORDER BY LastUpdateTime DESC`}
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
		var singerID, venueID int64
		var revenue spanner.NullInt64
		var eventDate, lastUpdateTime time.Time
		if err := row.ColumnByName("SingerId", &singerID); err != nil {
			return err
		}
		if err := row.ColumnByName("VenueId", &venueID); err != nil {
			return err
		}
		if err := row.ColumnByName("EventDate", &eventDate); err != nil {
			return err
		}
		if err := row.ColumnByName("Revenue", &revenue); err != nil {
			return err
		}
		currentRevenue := "NULL"
		if revenue.Valid {
			currentRevenue = strconv.FormatInt(revenue.Int64, 10)
		}
		if err := row.ColumnByName("LastUpdateTime", &lastUpdateTime); err != nil {
			return err
		}

		fmt.Fprintf(w, "%d %d %s %s %s\n", singerID, venueID, eventDate, currentRevenue, lastUpdateTime)
	}
}
