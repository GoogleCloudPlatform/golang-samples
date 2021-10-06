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

// [START spanner_set_transaction_tag]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

// readWriteTransactionWithTag executes the update and insert queries on venues table with appropriate transaction and requests tag
func readWriteTransactionWithTag(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `UPDATE Venues SET Capacity = CAST(Capacity/4 AS INT64) WHERE OutdoorVenue = false`,
		}
		_, err := txn.UpdateWithOptions(ctx, stmt, spanner.QueryOptions{RequestTag: "app=concert,env=dev,action=update"})
		if err != nil {
			return err
		}
		fmt.Fprint(w, "Venue capacities updated.")
		stmt = spanner.Statement{
			SQL: `INSERT INTO Venues (VenueId, VenueName, Capacity, OutdoorVenue, LastUpdateTime) 
                   VALUES (@venueId, @venueName, @capacity, @outdoorVenue, PENDING_COMMIT_TIMESTAMP())`,
			Params: map[string]interface{}{
				"venueId":      81,
				"venueName":    "Venue 81",
				"capacity":     1440,
				"outdoorVenue": true,
			},
		}
		_, err = txn.UpdateWithOptions(ctx, stmt, spanner.QueryOptions{RequestTag: "app=concert,env=dev,action=insert"})
		if err != nil {
			return err
		}
		fmt.Fprint(w, "New venue inserted.")
		return nil
	}, spanner.TransactionOptions{TransactionTag: "app=concert,env=dev"})
	return err
}

// [END spanner_set_transaction_tag]
