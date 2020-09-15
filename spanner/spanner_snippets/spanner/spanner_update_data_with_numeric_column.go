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

// [START spanner_update_data_with_numeric_column]
import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func updateDataWithNumericColumn(w io.Writer, db string) error {
	ctx := context.Background()

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	cols := []string{"VenueId", "Revenue"}
	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Venues", cols, []interface{}{4, "35000"}),
		spanner.Update("Venues", cols, []interface{}{19, "104500"}),
		spanner.Update("Venues", cols, []interface{}{42, "99999999999999999999999999999.99"}),
	})
	return err
}

// [END spanner_update_data_with_numeric_column]
