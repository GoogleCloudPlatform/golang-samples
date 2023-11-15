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

// [START spanner_batch_write_at_least_once]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/status"
)

// batchWrite demonstrates writing mutations to a Spanner database through
// BatchWrite API - https://pkg.go.dev/cloud.google.com/go/spanner#Client.BatchWrite
func batchWrite(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Database is assumed to exist - https://cloud.google.com/spanner/docs/getting-started/go#create_a_database
	singerColumns := []string{"SingerId", "FirstName", "LastName"}
	albumColumns := []string{"SingerId", "AlbumId", "AlbumTitle"}
	mutationGroups := make([]*spanner.MutationGroup, 2)

	mutationGroup1 := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{16, "Scarlet", "Terry"}),
	}
	mutationGroups[0] = &spanner.MutationGroup{Mutations: mutationGroup1}

	mutationGroup2 := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{17, "Marc", ""}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{18, "Catalina", "Smith"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{17, 1, "Total Junk"}),
		spanner.InsertOrUpdate("Albums", albumColumns, []interface{}{18, 2, "Go, Go, Go"}),
	}
	mutationGroups[1] = &spanner.MutationGroup{Mutations: mutationGroup2}
	iter := client.BatchWrite(ctx, mutationGroups)
	// See https://pkg.go.dev/cloud.google.com/go/spanner#BatchWriteResponseIterator.Do
	doFunc := func(response *sppb.BatchWriteResponse) error {
		if err = status.ErrorProto(response.GetStatus()); err == nil {
			fmt.Fprintf(w, "Mutation group indexes %v have been applied with commit timestamp %v",
				response.GetIndexes(), response.GetCommitTimestamp())
		} else {
			fmt.Fprintf(w, "Mutation group indexes %v could not be applied with error %v",
				response.GetIndexes(), err)
		}
		// Return an actual error as needed.
		return nil
	}
	return iter.Do(doFunc)
}

// [END spanner_batch_write_at_least_once]
