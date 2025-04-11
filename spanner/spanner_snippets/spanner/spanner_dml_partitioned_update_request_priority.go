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

// [START spanner_dml_partitioned_update_request_priority]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
)

// updateUsingPartitionedDMLRequestPriority updates a database using partitioned DML with request priority set
func updateUsingPartitionedDMLRequestPriority(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: "UPDATE Albums SET MarketingBudget = 100000 WHERE SingerId > 1"}
	rowCount, err := client.PartitionedUpdateWithOptions(ctx, stmt, spanner.QueryOptions{Priority: sppb.RequestOptions_PRIORITY_LOW})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%d record(s) updated.\n", rowCount)
	return nil
}

// [END spanner_dml_partitioned_update_request_priority]
