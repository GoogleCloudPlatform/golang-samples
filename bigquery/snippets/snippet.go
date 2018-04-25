// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains snippets for the Google BigQuery Go package.
package snippets

import (
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func createDataset(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_create_dataset]
	meta := &bigquery.DatasetMetadata{
		Location: "US", // Create the dataset in the US
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		return err
	}
	// [END bigquery_create_dataset]
	return nil
}

func updateDatasetDescription(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_update_dataset_description]
	ds := client.Dataset(datasetID)
	original, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	changes := bigquery.DatasetMetadataToUpdate{
		Description: "Updated Description.",
	}
	if _, err = ds.Update(ctx, changes, original.ETag); err != nil {
		return err
	}
	// [END bigquery_update_dataset_description]
	return nil
}

func updateDatasetDefaultExpiration(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_update_dataset_expiration]
	ds := client.Dataset(datasetID)
	original, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	changes := bigquery.DatasetMetadataToUpdate{
		DefaultTableExpiration: 24 * time.Hour,
	}
	if _, err := client.Dataset(datasetID).Update(ctx, changes, original.ETag); err != nil {
		return err
	}
	// [END bigquery_update_dataset_expiration]
	return nil
}

func updateDatasetAccessControl(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_update_dataset_access]
	ds := client.Dataset(datasetID)
	original, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	// Append a new access control entry to the existing access list
	changes := bigquery.DatasetMetadataToUpdate{
		Access: append(original.Access, &bigquery.AccessEntry{
			Role:       bigquery.ReaderRole,
			EntityType: bigquery.UserEmailEntity,
			Entity:     "sample.bigquery.dev@gmail.com"},
		),
	}

	// Leverage the ETag for the update to assert there's been no modifications to the
	// dataset since the metadata was originally read.
	if _, err := ds.Update(ctx, changes, original.ETag); err != nil {
		return err
	}
	// [END bigquery_update_dataset_access]
	return nil
}

func deleteEmptyDataset(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_delete_dataset]
	if err := client.Dataset(datasetID).Delete(ctx); err != nil {
		return fmt.Errorf("Failed to delete dataset: %v", err)
	}
	// [END bigquery_delete_dataset]
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

// Item represents a row item.
type Item struct {
	Name string
	Age  int
}

// Save implements the ValueSaver interface.
func (i *Item) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Name": i.Name,
		"Age":  i.Age,
	}, "", nil
}

func createTableInferredSchema(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// Demonstrates inferring a BigQuery schema from Go native types
	schema, err := bigquery.InferSchema(Item{})
	if err != nil {
		return err
	}
	table := client.Dataset(datasetID).Table(tableID)
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		return err
	}
	return nil
}

func createTableExplicitSchema(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_create_table]
	// Represent email_addrs as an array of strings, rather than a single address
	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0), // Table will be automatically deleted in 1 year
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	// [END bigquery_create_table]
	return nil
}

func createTableEmptySchema(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_create_table_without_schema]
	// Create the table as a partitioned table
	metaData := &bigquery.TableMetadata{
		TimePartitioning: &bigquery.TimePartitioning{
			Expiration: time.Duration(24*365) * time.Hour, // 365 day partition expiry
		},
	}
	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, metaData); err != nil {
		return err
	}
	// [END bigquery_create_table_without_schema]
	return nil
}

func updateTableDescription(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_update_table_description]
	tableRef := client.Dataset(datasetID).Table(tableID)
	original, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	newMeta := bigquery.TableMetadataToUpdate{
		Description: "Updated description.", // table expiration in 5 days
	}
	_, err = tableRef.Update(ctx, newMeta, original.ETag)
	if err != nil {
		return err
	}
	// [END bigquery_update_table_description]
	return nil

}

func updateTableExpiration(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_update_table_expiration]
	tableRef := client.Dataset(datasetID).Table(tableID)
	original, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	newMeta := bigquery.TableMetadataToUpdate{
		ExpirationTime: time.Now().Add(time.Duration(5*24) * time.Hour), // table expiration in 5 days
	}
	_, err = tableRef.Update(ctx, newMeta, original.ETag)
	if err != nil {
		return err
	}
	// [END bigquery_update_table_expiration]
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
	// [START bigquery_table_insert_rows]
	u := client.Dataset(datasetID).Table(tableID).Uploader()
	items := []*Item{
		// Item implements the ValueSaver interface.
		{Name: "Phred Phlyntstone", Age: 32},
		{Name: "Wylma Phlyntstone", Age: 29},
	}
	if err := u.Put(ctx, items); err != nil {
		return err
	}
	// [END bigquery_table_insert_rows]
	return nil
}

func listRows(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	q := client.Query(fmt.Sprintf(`
		SELECT name, age
		FROM %s.%s
		WHERE age >= 20
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
	return nil
}

func basicQuery(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_query]
	q := client.Query(
		"SELECT name FROM `bigquery-public-data.usa_names.usa_1910_2013` " +
			"WHERE state = \"TX\" " +
			"LIMIT 100")
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"

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
	// [END bigquery_query]
	return nil
}

func printTableMetadataSimple(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_get_table]
	meta, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return err
	}
	// Print information about the table
	fmt.Printf("Schema: %+v\n", meta.Schema)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Printf("Row in managed storage: %d\n", meta.NumRows)
	// [END bigquery_get_table]
	return nil
}
func printTableMetadataExtended(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// Demonstrates some of the information available in the metadata for a table.
	meta, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return err
	}
	// Print information about the table
	fmt.Printf("Table ID: %s\n", tableID)
	if len(meta.Name) > 0 {
		fmt.Printf("Has Friendly Name: %s", meta.Name)
	}

	if meta.Type == bigquery.ViewTable {
		// Table is a logical view, rather than a table with backing storage
		fmt.Println("Table is a logical view")
	} else {
		if meta.Type == bigquery.ExternalTable {
			// Table is federated against external data (GCS, BigTable, etc)
			fmt.Printf("Table is externally federated against %s\n", meta.ExternalDataConfig.SourceFormat)
		} else {
			// Table is a normal managed table.
			fmt.Printf("Table is managed by BigQuery, with %d rows and %d Bytes in managed storage.\n", meta.NumRows, meta.NumBytes)
			if meta.StreamingBuffer != nil {
				fmt.Printf("Table has active streaming buffer, with estimated %d bytes and %d rows in the buffer\n", meta.StreamingBuffer.EstimatedBytes, meta.StreamingBuffer.EstimatedRows)
			}
		}
		// Its a table, walk the top-level schema
		fmt.Printf("Defined schema has %d top-level fields\n", len(meta.Schema))
		topLevelArrays := 0
		topLevelRecords := 0
		for _, f := range meta.Schema {
			if f.Repeated {
				topLevelArrays++
			}
			if f.Type == bigquery.RecordFieldType {
				topLevelRecords++
			}
		}
		if topLevelArrays > 0 || topLevelRecords > 0 {
			fmt.Printf("Schema is complex.  %d fields are array-based, %d fields are record/structs.\n", topLevelArrays, topLevelRecords)
		}
	}
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

func importCSVFromFile(client *bigquery.Client, datasetID, tableID, filename string) error {
	ctx := context.Background()
	// [START bigquery_load_from_file]
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	source := bigquery.NewReaderSource(f)
	source.AutoDetect = true   // Allow BigQuery to determine schema
	source.SkipLeadingRows = 1 // CSV has a single header line

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(source)

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
	// [END bigquery_load_from_file]
	return nil
}

func importCSVExplicitSchema(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_load_table_gcs_csv]
	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/us-states/us-states.csv")
	gcsRef.SkipLeadingRows = 1
	gcsRef.Schema = bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "post_abbr", Type: bigquery.StringFieldType},
	}
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteEmpty

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("Job completed with error: %v", status.Err())
	}
	// [END bigquery_load_table_gcs_csv]
	return nil
}

func importJSONExplicitSchema(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_load_table_gcs_json]
	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/us-states/us-states.json")
	gcsRef.SourceFormat = bigquery.JSON
	gcsRef.Schema = bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "post_abbr", Type: bigquery.StringFieldType},
	}
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteEmpty

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("Job completed with error: %v", status.Err())
	}
	// [END bigquery_load_table_gcs_json]
	return nil
}

func importJSONAutodetectSchema(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_load_table_gcs_json_autodetect]
	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/us-states/us-states.json")
	gcsRef.SourceFormat = bigquery.JSON
	gcsRef.AutoDetect = true
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteEmpty

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("Job completed with error: %v", status.Err())
	}
	// [END bigquery_load_table_gcs_json_autodetect]
	return nil
}

func exportSampleTableAsCSV(client *bigquery.Client, gcsURI string) error {
	ctx := context.Background()
	// [START bigquery_extract_table]
	srcProject := "bigquery-public-data"
	srcDataset := "samples"
	srcTable := "shakespeare"

	// For example, gcsUri = "gs://mybucket/shakespeare.csv"
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.FieldDelimiter = ","

	extractor := client.DatasetInProject(srcProject, srcDataset).Table(srcTable).ExtractorTo(gcsRef)
	extractor.DisableHeader = true
	// You can choose to run the job in a specific location for more complex data locality scenarios
	// Ex: In this example, source dataset and GCS bucket are in the US
	extractor.Location = "US"

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
	// [END bigquery_extract_table]
	return nil
}

func exportSampleTableAsCompressedCSV(client *bigquery.Client, gcsURI string) error {
	ctx := context.Background()
	// [START bigquery_extract_table_compressed]
	srcProject := "bigquery-public-data"
	srcDataset := "samples"
	srcTable := "shakespeare"

	// For example, gcsUri = "gs://mybucket/shakespeare.csv"
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.Compression = bigquery.Gzip

	extractor := client.DatasetInProject(srcProject, srcDataset).Table(srcTable).ExtractorTo(gcsRef)
	extractor.DisableHeader = true
	// You can choose to run the job in a specific location for more complex data locality scenarios
	// Ex: In this example, source dataset and GCS bucket are in the US	extractor.Location = "US"
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
	// [END bigquery_extract_table_compressed]
	return nil
}

func exportSampleTableAsJSON(client *bigquery.Client, gcsURI string) error {
	ctx := context.Background()
	// [START bigquery_extract_table_json]
	srcProject := "bigquery-public-data"
	srcDataset := "samples"
	srcTable := "shakespeare"

	// For example, gcsUri = "gs://mybucket/shakespeare.json"
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.DestinationFormat = bigquery.JSON

	extractor := client.DatasetInProject(srcProject, srcDataset).Table(srcTable).ExtractorTo(gcsRef)
	// You can choose to run the job in a specific location for more complex data locality scenarios
	// Ex: In this example, source dataset and GCS bucket are in the US
	extractor.Location = "US"

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
	// [END bigquery_extract_table_json]
	return nil
}
