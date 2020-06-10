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

// [START spanner_write_data_for_struct_queries]

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func writeStructData(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	singerColumns := []string{"SingerId", "FirstName", "LastName"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{6, "Elena", "Campbell"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{7, "Gabriel", "Wright"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{8, "Benjamin", "Martinez"}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{9, "Hannah", "Harris"}),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_write_data_for_struct_queries]
