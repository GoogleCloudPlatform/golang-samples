// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command syncquery queries a Google BigQuery dataset.
package main

import (
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"golang.org/x/net/context"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ", os.Args[0], "'query text'")
		os.Exit(1)
	}
	proj := os.Getenv("GCLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GCLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	rows, err := Query(proj, os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	for _, row := range rows {
		fmt.Println(row)
	}
}

// Query returns a slice of the results of a query.
func Query(proj, q string) ([][]bigquery.Value, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		return nil, err
	}

	query := client.Query(q)
	iter, err := query.Read(ctx)
	if err != nil {
		return nil, err
	}

	var rows [][]bigquery.Value

	for {
		var row []bigquery.Value
		err := iter.Next(&row)
		if err == iterator.Done {
			return rows, nil
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
}
