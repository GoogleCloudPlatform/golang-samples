// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command simpleapp queries the Shakespeare sample dataset in Google BigQuery.
// [START bigquery_simple_app_all]
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"golang.org/x/net/context"
)

func main() {
	proj := os.Getenv("GCLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GCLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	rows, err := Query(proj)
	if err != nil {
		log.Fatalln(err)
	}
	err = PrintResults(os.Stdout, rows)
	if err != nil {
		log.Fatalln(err)
	}
}

// Query returns a slice of the results of a query.
// [START bigquery_simple_app_query]
func Query(proj string) (*bigquery.RowIterator, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		return nil, err
	}

	query := client.Query(
		"SELECT " +
			"APPROX_TOP_COUNT(corpus, 10) as title, " +
			"COUNT(*) as unique_words " +
			"FROM `publicdata.samples.shakespeare`;")
	// Use standard SQL syntax for queries.
	// See: https://cloud.google.com/bigquery/sql-reference/
	query.QueryConfig.UseStandardSQL = true
	return query.Read(ctx)
}

// [END bigquery_simple_app_query]

// PrintResults prints results of a query to the Shakespeare dataset.
// [START bigquery_simple_app_print]
func PrintResults(w io.Writer, iter *bigquery.RowIterator) error {
	for {
		var row bigquery.ValueList
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}

		// TODO: use ValueLoader to load a struct instead of using the ValueList directly.
		// See: https://github.com/GoogleCloudPlatform/google-cloud-go/issues/399
		fmt.Fprintln(w, "titles:")
		ts := row[0].([]bigquery.Value)
		for _, t := range ts {
			record := t.([]bigquery.Value)
			title := record[0].(string)
			cnt := record[1].(int)
			fmt.Fprintf(w, "\t%s: %d\n", title, cnt)
		}

		words := row[1].(int)
		fmt.Fprintf(w, "total unique words: %d\n", words)
	}
}

// [END bigquery_simple_app_print]
// [END bigquery_simple_app_all]
