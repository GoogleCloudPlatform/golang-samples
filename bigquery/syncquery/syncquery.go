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

// Command syncquery queries a Google BigQuery dataset.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ", os.Args[0], "'query text'")
		os.Exit(1)
	}
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
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
	defer client.Close()

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
