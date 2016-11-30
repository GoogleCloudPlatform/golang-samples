// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains snippets for the Google BigQuery Go package.
package snippets

import (
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func createDataset(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_create_dataset]
	if err := client.Dataset(datasetID).Create(ctx); err != nil {
		return err
	}
	// [END bigquery_create_dataset]
	return nil
}

func listDatasets(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_list_datasets]
	it := client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		fmt.Println(dataset.DatasetID)
	}
	// [END bigquery_list_datasets]
	return nil
}

// [START bigquery_create_table]

// Item represents a row item.
type Item struct {
	Name  string
	Count int
}

// Save implements the ValueSaver interface.
func (i *Item) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Name":  i.Name,
		"Count": i.Count,
	}, "", nil
}

// [END bigquery_create_table]

func createTable(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_create_table]
	schema, err := bigquery.InferSchema(Item{})
	if err != nil {
		return err
	}
	table := client.Dataset(datasetID).Table(tableID)
	if err := table.Create(ctx, schema); err != nil {
		return err
	}
	// [END bigquery_create_table]
	return nil
}

func listTables(client *bigquery.Client, w io.Writer, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_list_tables]
	ts := client.Dataset(datasetID).Tables(ctx)
	for {
		t, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Table: %q\n", t.TableID)
	}
	// [END bigquery_list_tables]
	return nil
}

func insertRows(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_insert_stream]
	u := client.Dataset(datasetID).Table(tableID).Uploader()
	items := []*Item{
		// Item implements the ValueSaver interface.
		{Name: "n1", Count: 7},
		{Name: "n2", Count: 2},
		{Name: "n3", Count: 1},
	}
	if err := u.Put(ctx, items); err != nil {
		return err
	}
	// [END bigquery_insert_stream]
	return nil
}

func listRows(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_list_rows]
	q := client.Query(fmt.Sprintf(`
		SELECT name, count
		FROM [%s.%s]
		WHERE count >= 5
	`, datasetID, tableID))
	it, err := q.Read(ctx)
	if err != nil {
		return err
	}

	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(row)
	}
	// [END bigquery_list_rows]
	return nil
}

func asyncQuery(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_async_query]
	q := client.Query(fmt.Sprintf(`
		SELECT name, count
		FROM [%s.%s]
	`, datasetID, tableID))
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}

	// Wait until async querying is done.
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}

	it, err := job.Read(ctx)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(row)
	}
	// [END bigquery_async_query]
	return nil
}

func browseTable(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_browse_table]
	table := client.Dataset(datasetID).Table(tableID)
	it := table.Read(ctx)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(row)
	}
	// [END bigquery_browse_table]
	return nil
}

func copyTable(client *bigquery.Client, datasetID, srcID, dstID string) error {
	ctx := context.Background()
	// [START bigquery_copy_table]
	dataset := client.Dataset(datasetID)
	copier := dataset.Table(dstID).CopierFrom(dataset.Table(srcID))
	copier.WriteDisposition = bigquery.WriteTruncate
	job, err := copier.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	// [END bigquery_copy_table]
	return nil
}

func deleteTable(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_delete_table]
	table := client.Dataset(datasetID).Table(tableID)
	if err := table.Delete(ctx); err != nil {
		return err
	}
	// [END bigquery_delete_table]
	return nil
}

func importFromGCS(client *bigquery.Client, datasetID, tableID, gcsURI string) error {
	ctx := context.Background()
	// [START bigquery_import_from_gcs]
	// For example, "gs://data-bucket/path/to/data.csv"
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.AllowJaggedRows = true
	// TODO: set other options on the GCSReference.

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.CreateDisposition = bigquery.CreateNever
	// TODO: set other options on the Loader.

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	// [END bigquery_import_from_gcs]
	return nil
}

func importFromFile(client *bigquery.Client, datasetID, tableID, filename string) error {
	ctx := context.Background()
	// [START bigquery_import_from_file]
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	source := bigquery.NewReaderSource(f)
	source.AllowJaggedRows = true
	// TODO: set other options on the GCSReference.

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(source)
	loader.CreateDisposition = bigquery.CreateNever
	// TODO: set other options on the Loader.

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	// [END bigquery_import_from_file]
	return nil
}

func exportToGCS(client *bigquery.Client, datasetID, tableID, gcsURI string) error {
	ctx := context.Background()
	// [START bigquery_export_gcs]
	// For example, "gs://data-bucket/path/to/data.csv"
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.FieldDelimiter = ","

	extractor := client.Dataset(datasetID).Table(tableID).ExtractorTo(gcsRef)
	extractor.DisableHeader = true
	job, err := extractor.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	// [END bigquery_export_gcs]
	return nil
}
