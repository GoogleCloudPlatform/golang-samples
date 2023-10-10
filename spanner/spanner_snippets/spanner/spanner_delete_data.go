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

// [START spanner_delete_data]

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func delete(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	m := []*spanner.Mutation{
		// spanner.Key can be used to delete a specific set of rows.
		// Delete the Albums with the key values (2,1) and (2,3).
		spanner.Delete("Albums", spanner.Key{2, 1}),
		spanner.Delete("Albums", spanner.Key{2, 3}),
		// spanner.KeyRange can be used to delete rows with a key in a specific range.
		// Delete a range of rows where the column key is >=3 and <5
		spanner.Delete("Singers", spanner.KeyRange{Start: spanner.Key{3}, End: spanner.Key{5}, Kind: spanner.ClosedOpen}),
		// spanner.AllKeys can be used to delete all the rows in a table.
		// Delete remaining Singers rows, which will also delete the remaining Albums rows since it was
		// defined with ON DELETE CASCADE.
		spanner.Delete("Singers", spanner.AllKeys()),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_delete_data]
