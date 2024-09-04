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

// [START spanner_delete_graph_data]

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func deleteGraphData(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Apply a series of mutations to tables underpinning edges and nodes in the
	// example graph. If there are referential integrity constraints defined
	// between edges and the nodes they connect, the edge must be deleted
	// before the nodes that the edge connects are deleted.
	m := []*spanner.Mutation{
		// spanner.Key can be used to delete a specific set of rows.
		// Delete the PersonOwnAccount rows with the key values (1,7) and (2,20).
		spanner.Delete("PersonOwnAccount", spanner.Key{1, 7}),
		spanner.Delete("PersonOwnAccount", spanner.Key{2, 20}),

		// spanner.KeyRange can be used to delete rows with a key in a specific range.
		// Delete a range of rows where the key prefix is >=1 and <8
		spanner.Delete("AccountTransferAccount",
			spanner.KeyRange{Start: spanner.Key{1}, End: spanner.Key{8}, Kind: spanner.ClosedOpen}),

		// spanner.AllKeys can be used to delete all the rows in a table.
		// Delete all Account rows, which will also delete the remaining
		// AccountTransferAccount rows since it was defined with ON DELETE CASCADE.
		spanner.Delete("Account", spanner.AllKeys()),

		// Delete remaining Person rows, which will also delete the remaining
		// PersonOwnAccount rows since it was defined with ON DELETE CASCADE.
		spanner.Delete("Person", spanner.AllKeys()),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_delete_graph_data]
