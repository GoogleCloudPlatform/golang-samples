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

// [START spanner_quickstart]

// Sample spanner_quickstart is a basic program that uses Cloud Spanner.
package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	// This database must exist.
	databaseName := "projects/your-project-id/instances/your-instance-id/databases/your-database-id"

	client, err := spanner.NewClient(ctx, databaseName)
	if err != nil {
		log.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: "SELECT 1"}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			fmt.Println("Done")
			return
		}
		if err != nil {
			log.Fatalf("Query failed with %v", err)
		}

		var i int64
		if row.Columns(&i) != nil {
			log.Fatalf("Failed to parse row %v", err)
		}
		fmt.Printf("Got value %v\n", i)
	}
}

// [END spanner_quickstart]
