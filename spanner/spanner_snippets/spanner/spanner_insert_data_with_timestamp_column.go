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

// [START spanner_insert_data_with_timestamp_column]

import (
	"context"

	"cloud.google.com/go/spanner"
)

func writeWithTimestamp(db string) error {
	ctx := context.Background()

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	performanceColumns := []string{"SingerId", "VenueId", "EventDate", "Revenue", "LastUpdateTime"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Performances", performanceColumns, []interface{}{1, 4, "2017-10-05", 11000, spanner.CommitTimestamp}),
		spanner.InsertOrUpdate("Performances", performanceColumns, []interface{}{1, 19, "2017-11-02", 15000, spanner.CommitTimestamp}),
		spanner.InsertOrUpdate("Performances", performanceColumns, []interface{}{2, 42, "2017-12-23", 7000, spanner.CommitTimestamp}),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_insert_data_with_timestamp_column]
