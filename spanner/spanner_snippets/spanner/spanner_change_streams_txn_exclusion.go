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

// [START spanner_set_exclude_txn_from_change_streams]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/status"
)

// rwTxnExcludedFromChangeStreams executes the insert and update DMLs on Singers table excluded from allowed tracking change streams
func rwTxnExcludedFromChangeStreams(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("rwTxnExcludedFromChangeStreams.NewClient: %w", err)
	}
	defer client.Close()

	_, err = client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName)
					VALUES (111, 'Virginia', 'Watson')`,
		}
		_, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("rwTxnExcludedFromChangeStreams.Update: %w", err)
		}
		fmt.Fprintf(w, "New singer inserted.")
		stmt = spanner.Statement{
			SQL: `UPDATE Singers SET FirstName = 'Hi' WHERE SingerId = 111`,
		}
		_, err = txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("rwTxnExcludedFromChangeStreams.Update: %w", err)
		}
		fmt.Fprint(w, "Singer first name updated.")
		return nil
	}, spanner.TransactionOptions{ExcludeTxnFromChangeStreams: true})
	if err != nil {
		return err
	}
	return nil
}

// applyExcludedFromChangeStreams apply the insert mutations on Singers table excluded from allowed tracking change streams
func applyExcludedFromChangeStreams(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("applyExcludedFromChangeStreams.NewClient: %w", err)
	}
	defer client.Close()
	m := spanner.Insert("Singers",
		[]string{"SingerId", "FirstName", "LastName"},
		[]interface{}{999, "Foo", "Bar"})
	_, err = client.Apply(ctx, []*spanner.Mutation{m}, spanner.ExcludeTxnFromChangeStreams())

	if err != nil {
		return err
	}
	fmt.Fprint(w, "applyExcludedFromChangeStreams.Apply: New singer inserted.")
	return err
}

// applyExcludedFromChangeStreams apply the insert mutations on Singers table excluded from allowed tracking change streams
func applyAtLeastOnceExcludedFromChangeStreams(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("applyExcludedFromChangeStreams.NewClient: %w", err)
	}
	defer client.Close()
	m := spanner.Insert("Singers",
		[]string{"SingerId", "FirstName", "LastName"},
		[]interface{}{989, "Hellen", "Lee"})
	_, err = client.Apply(ctx, []*spanner.Mutation{m}, []spanner.ApplyOption{spanner.ExcludeTxnFromChangeStreams(), spanner.ApplyAtLeastOnce()}...)

	if err != nil {
		return err
	}
	fmt.Fprint(w, "applyExcludedFromChangeStreams.ApplyAtLeastOnce: New singer inserted.")
	return err
}

// batchWriteExcludedFromChangeStreams executes the insert mutation on Singers table excluded from allowed tracking change streams
func batchWriteExcludedFromChangeStreams(w io.Writer, db string) error {
	// db := "projects/my-project/instances/my-instance/databases/my-database"
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("batchWriteExcludedFromChangeStreams.NewClient: %w", err)
	}
	defer client.Close()

	singerColumns := []string{"SingerId", "FirstName", "LastName"}
	mutationGroups := make([]*spanner.MutationGroup, 1)

	mutationGroup1 := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{0127, "Scarlet", "Terry"}),
	}
	mutationGroups[0] = &spanner.MutationGroup{Mutations: mutationGroup1}

	iter := client.BatchWriteWithOptions(ctx, mutationGroups, spanner.BatchWriteOptions{ExcludeTxnFromChangeStreams: true})
	// See https://pkg.go.dev/cloud.google.com/go/spanner#BatchWriteResponseIterator.Do
	doFunc := func(response *sppb.BatchWriteResponse) error {
		if err = status.ErrorProto(response.GetStatus()); err == nil {
			fmt.Fprintf(w, "batchWriteExcludedFromChangeStreams.BatchWriteWithOptions: Mutation group indexes %v have been applied with commit timestamp %v",
				response.GetIndexes(), response.GetCommitTimestamp())
		} else {
			fmt.Fprintf(w, "batchWriteExcludedFromChangeStreams.BatchWriteWithOptions: Mutation group indexes %v could not be applied with error %v",
				response.GetIndexes(), err)
		}
		// Return an actual error as needed.
		return nil
	}
	return iter.Do(doFunc)
}

// pdmlExcludedFromChangeStreams executes the partitioned update DML on Singers table excluded from allowed tracking change streams
func pdmlExcludedFromChangeStreams(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("pdmlExcludedFromChangeStreams.NewClient: %w", err)
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: "UPDATE Singers SET FirstName = 'Hello' WHERE SingerId > 500"}
	rowCount, err := client.PartitionedUpdateWithOptions(ctx, stmt, spanner.QueryOptions{ExcludeTxnFromChangeStreams: true})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "pdmlExcludedFromChangeStreams.PartitionedUpdateWithOptions: %d record(s) updated.\n", rowCount)
	return nil
}

// [END spanner_set_exclude_txn_from_change_streams]
