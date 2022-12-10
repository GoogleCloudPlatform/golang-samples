// Copyright 2022 Google LLC
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

// [START spanner_delete_data_with_proto_columns_using_dml]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
)

func deleteDataWithProtoMsgAndEnumUsingDML(w io.Writer, db string) error {
	ctx := context.Background()
	endpoint := "staging-wrenchworks.sandbox.googleapis.com:443"
	client, err := spanner.NewClient(ctx, db, option.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{SQL: `DELETE FROM Singers WHERE SingerId = 5`}

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) deleted.\n", rowCount)
		return nil
	})
	return err
}

// [END spanner_delete_data_with_proto_columns_using_dml]
