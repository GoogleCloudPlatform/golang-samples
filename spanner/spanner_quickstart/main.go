// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample spanner_quickstart is a basic program that uses Cloud Spanner.
package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// This database must exist.
	databaseID := "projects/your-project-id/instances/your-instance-id/databases/your-database-id"

	client, err := spanner.NewClient(ctx, databaseID)
	if err != nil {
		log.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: "SELECT 1"}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		log.Fatalf("Query failed with %v", err)
	}

	var i int64
	if row.Columns(&i) != nil {
		log.Fatalf("Failed to parse row %v", err)
	}
	fmt.Printf("Got value %v\n", i)
}
