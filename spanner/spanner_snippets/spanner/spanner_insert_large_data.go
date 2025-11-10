// Copyright 2025 Google LLC
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

// [START spanner_insert_large_data]

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"cloud.google.com/go/spanner"
)

func writeLargeData(w io.Writer, db string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	singerColumns := []string{"SingerId", "FirstName", "LastName", "SingerInfo"}
	token := make([]byte, 10000000)
	if _, err := rand.Read(token); err != nil {
		return err
	}
	// Mutation is under the 100MB limit
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{1, "Marc", "Richards", token}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{2, "Catalina", "Smith", token}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{3, "Alice", "Trentor", token}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{4, "Lea", "Martin", token}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{5, "David", "Lomond", token}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{6, "Marc", "Richards", token}),
		spanner.InsertOrUpdate("Singers", singerColumns, []interface{}{7, "Catalina", "Smith", token}),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_insert_large_data]
