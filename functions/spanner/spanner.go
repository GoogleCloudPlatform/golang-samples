// Copyright 2019 Google LLC
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

// [START spanner_functions_quickstart]

// Package spanner contains an example of using Spanner from a Cloud Function.
package spanner

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// client is a global Spanner client, to avoid initializing a new client for
// every request.
var client *spanner.Client
var clientOnce sync.Once

// db is the name of the database to query.
var db = "projects/my-project/instances/my-instance/databases/example-db"

// HelloSpanner is an example of querying Spanner from a Cloud Function.
func HelloSpanner(w http.ResponseWriter, r *http.Request) {
	clientOnce.Do(func() {
		// Declare a separate err variable to avoid shadowing client.
		var err error
		client, err = spanner.NewClient(context.Background(), db)
		if err != nil {
			http.Error(w, "Error initializing database", http.StatusInternalServerError)
			log.Printf("spanner.NewClient: %v", err)
			return
		}
	})

	fmt.Fprintln(w, "Albums:")
	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := client.Single().Query(r.Context(), stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return
		}
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			log.Printf("iter.Next: %v", err)
			return
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			http.Error(w, "Error parsing database response", http.StatusInternalServerError)
			log.Printf("row.Columns: %v", err)
			return
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_functions_quickstart]
