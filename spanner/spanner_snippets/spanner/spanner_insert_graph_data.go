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

// [START spanner_insert_graph_data]

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/spanner"
)

func parseTime(rfc3339Time string) time.Time {
	t, _ := time.Parse(time.RFC3339, rfc3339Time)
	return t
}

func insertGraphData(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	// Values are inserted into the node and edge tables corresponding to
	// using Spanner 'Insert' mutations.
	// The tables and columns comply with the schema defined for the
	// property graph 'FinGraph', comprising 'Person' and 'Account' nodes,
	// and 'PersonOwnAccount' and 'AccountTransferAccount' edges.
	personColumns := []string{"id", "name", "birthday", "country", "city"}
	accountColumns := []string{"id", "create_time", "is_blocked", "nick_name"}
	ownColumns := []string{"id", "account_id", "create_time"}
	transferColumns := []string{"id", "to_id", "amount", "create_time", "order_number"}
	m := []*spanner.Mutation{
		spanner.Insert("Account", accountColumns,
			[]interface{}{7, parseTime("2020-01-10T06:22:20.12Z"), false, "Vacation Fund"}),
		spanner.Insert("Account", accountColumns,
			[]interface{}{16, parseTime("2020-01-27T17:55:09.12Z"), true, "Vacation Fund"}),
		spanner.Insert("Account", accountColumns,
			[]interface{}{20, parseTime("2020-02-18T05:44:20.12Z"), false, "Rainy Day Fund"}),
		spanner.Insert("Person", personColumns,
			[]interface{}{1, "Alex", parseTime("1991-12-21T00:00:00.12Z"), "Australia", " Adelaide"}),
		spanner.Insert("Person", personColumns,
			[]interface{}{2, "Dana", parseTime("1980-10-31T00:00:00.12Z"), "Czech_Republic", "Moravia"}),
		spanner.Insert("Person", personColumns,
			[]interface{}{3, "Lee", parseTime("1986-12-07T00:00:00.12Z"), "India", "Kollam"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{7, 16, 300.0, parseTime("2020-08-29T15:28:58.12Z"), "304330008004315"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{7, 16, 100.0, parseTime("2020-10-04T16:55:05.12Z"), "304120005529714"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{16, 20, 300.0, parseTime("2020-09-25T02:36:14.12Z"), "103650009791820"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{20, 7, 500.0, parseTime("2020-10-04T16:55:05.12Z"), "304120005529714"}),
		spanner.Insert("AccountTransferAccount", transferColumns,
			[]interface{}{20, 16, 200.0, parseTime("2020-10-17T03:59:40.12Z"), "302290001255747"}),
		spanner.Insert("PersonOwnAccount", ownColumns,
			[]interface{}{1, 7, parseTime("2020-01-10T06:22:20.12Z")}),
		spanner.Insert("PersonOwnAccount", ownColumns,
			[]interface{}{2, 20, parseTime("2020-01-27T17:55:09.12Z")}),
		spanner.Insert("PersonOwnAccount", ownColumns,
			[]interface{}{3, 16, parseTime("2020-02-18T05:44:20.12Z")}),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_insert_graph_data]
