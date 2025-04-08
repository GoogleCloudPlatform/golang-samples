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

// [START spanner_update_data_with_timestamp_column]

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func updateWithTimestamp(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	cols := []string{"SingerId", "AlbumId", "MarketingBudget", "LastUpdateTime"}
	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Albums", cols, []interface{}{1, 1, 1000000, spanner.CommitTimestamp}),
		spanner.Update("Albums", cols, []interface{}{2, 2, 750000, spanner.CommitTimestamp}),
	})
	return err
}

// [END spanner_update_data_with_timestamp_column]
