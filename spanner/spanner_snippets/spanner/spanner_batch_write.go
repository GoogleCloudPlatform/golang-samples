// Copyright 2023 Google LLC
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

// [START spanner_batch_write]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

func batchWrite(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	singerColumns := []string{"SingerId", "FirstName", "LastName"}
	albumColumns := []string{"SingerId", "AlbumId", "AlbumTitle"}
	mutationGroups := make([]*spanner.MutationGroup, 2)

	group1 := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{16, "Scarlet", "Terry"}),
	}
	mutationGroups[0] = &spanner.MutationGroup{Mutations: group1}

	group2 := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{17, "Marc", ""}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{18, "Catalina", "Smith"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{17, 1, "Total Junk"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{18, 2, "Go, Go, Go"}),
	}
	mutationGroups[1] = &spanner.MutationGroup{Mutations: group2}

	iter := client.BatchWrite(ctx, mutationGroups)
	doFunc := func(response *sppb.BatchWriteResponse) error {
		if ts := response.GetCommitTimestamp(); ts == nil {
			return fmt.Errorf("invalid commit timesamp")
		}
		return nil
	}
	if err = iter.Do(doFunc); err != nil {
		return err
	}
	fmt.Fprintf(w, "BatchWrite successful")
	return nil
}

// [END spanner_batch_write]
