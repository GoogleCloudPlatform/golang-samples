// Copyright 2024 Google LLC
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

// [START spanner_query_graph_data_with_parameter]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryGraphDataWithParameter(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Execute a Spanner query statement comprising a graph query. Graph queries
	// are characterized by 'MATCH' statements describing node and edge
	// patterns.
	//
	// This statement finds entities ('Account's) owned by all 'Person's 'b' to
	// which transfers greater than a specified 'min' have been made by entities
	// ('Account's) owned by any 'Person' 'a' in the graph called 'FinGraph'.
	// It then returns the names of all such 'Person's 'a' and 'b', and the
	// amount and time of the transfer.
	// The parameter 'min' is distinguished by a leading '@', and the value
	// of the parameter is specified using a map.
	stmt := spanner.Statement{
		SQL: `Graph FinGraph 
			  MATCH (a:Person)-[o:Owns]->()-[t:Transfers]->()<-[p:Owns]-(b:Person)
			  WHERE t.amount >= @min
			  RETURN a.name AS sender, b.name AS receiver, t.amount, t.create_time AS transfer_at`,
		Params: map[string]interface{}{
			"min": 500,
		},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	// The results are returned in tabular form. Iterate over the
	// result rows and print them.
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var sender, receiver string
		var amount float64
		var transfer_at time.Time
		if err := row.Columns(&sender, &receiver, &amount, &transfer_at); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s %s %f %s\n", sender, receiver, amount, transfer_at.Format(time.RFC3339))
	}
}

// [END spanner_query_graph_data_with_parameter]
