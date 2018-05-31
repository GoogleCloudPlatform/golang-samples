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
		Location: "US", // Create the dataset in the US.
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
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.DatasetMetadataToUpdate{
		Description: "Updated Description.",
	}
	if _, err = ds.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_update_dataset_description]
	return nil
}

func updateDatasetDefaultExpiration(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_update_dataset_expiration]
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.DatasetMetadataToUpdate{
		DefaultTableExpiration: 24 * time.Hour,
	}
	if _, err := client.Dataset(datasetID).Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_update_dataset_expiration]
	return nil
}

func updateDatasetAccessControl(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_update_dataset_access]
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	// Append a new access control entry to the existing access list.
	update := bigquery.DatasetMetadataToUpdate{
		Access: append(meta.Access, &bigquery.AccessEntry{
			Role:       bigquery.ReaderRole,
			EntityType: bigquery.UserEmailEntity,
			Entity:     "sample.bigquery.dev@gmail.com"},
		),
	}

	// Leverage the ETag for the update to assert there's been no modifications to the
	// dataset since the metadata was originally read.
	if _, err := ds.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_update_dataset_access]
	return nil
}

func datasetLabels(client *bigquery.Client, w io.Writer, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_get_dataset_labels]
	meta, err := client.Dataset(datasetID).Metadata(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Dataset %s labels:\n", datasetID)
	if len(meta.Labels) == 0 {
		fmt.Fprintln(w, "Dataset has no labels defined.")
		return nil
	}
	for k, v := range meta.Labels {
		fmt.Fprintf(w, "\t%s:%s\n", k, v)
	}
	// [END bigquery_get_dataset_labels]
	return nil
}

func addDatasetLabel(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_label_dataset]
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}

	update := bigquery.DatasetMetadataToUpdate{}
	update.SetLabel("color", "green")
	if _, err := ds.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_label_dataset]
	return nil
}

func deleteDatasetLabel(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_delete_label_dataset]
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.DatasetMetadataToUpdate{}
	update.DeleteLabel("color")
	if _, err := ds.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_delete_label_dataset]
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

func listDatasetsByLabel(client *bigquery.Client, w io.Writer) error {
	ctx := context.Background()
	// [START bigquery_list_datasets_by_label]
	it := client.Datasets(ctx)
	it.Filter = "labels.color:green"
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "dataset: %s\n", dataset.DatasetID)
	}
	// [END bigquery_list_datasets_by_label]
	return nil
}

func printDatasetInfo(client *bigquery.Client, datasetID string) error {
	ctx := context.Background()
	// [START bigquery_get_dataset]
	meta, err := client.Dataset(datasetID).Metadata(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Dataset ID: %s\n", datasetID)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Println("Labels:")
	for k, v := range meta.Labels {
		fmt.Printf("\t%s: %s", k, v)
	}
	fmt.Println("Tables:")
	it := client.Dataset(datasetID).Tables(ctx)

	cnt := 0
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		cnt++
		fmt.Printf("\t%s\n", t.TableID)
	}
	if cnt == 0 {
		fmt.Println("\tThis dataset does not contain any tables.")
	}
	// [END bigquery_get_dataset]
	return nil
}

func listJobs(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_list_jobs]
	it := client.Jobs(ctx)
	for i := 0; i < 10; i++ {
		j, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		state := "Unknown"
		switch j.LastStatus().State {
		case bigquery.Pending:
			state = "Pending"
		case bigquery.Running:
			state = "Running"
		case bigquery.Done:
			state = "Done"
		}
		fmt.Printf("Job %s in state %s\n", j.ID(), state)
	}
	// [END bigquery_list_jobs]
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
	// bigquery.InferSchema infers BQ schema from native Go types.
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
	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0), // Table will be automatically deleted in 1 year.
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
	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, nil); err != nil {
		return err
	}
	// [END bigquery_create_table_without_schema]
	return nil
}

func updateTableDescription(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_update_table_description]
	tableRef := client.Dataset(datasetID).Table(tableID)
	meta, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.TableMetadataToUpdate{
		Description: "Updated description.",
	}
	if _, err = tableRef.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_update_table_description]
	return nil

}

func updateTableExpiration(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_update_table_expiration]
	tableRef := client.Dataset(datasetID).Table(tableID)
	meta, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.TableMetadataToUpdate{
		ExpirationTime: time.Now().Add(time.Duration(5*24) * time.Hour), // table expiration in 5 days.
	}
	if _, err = tableRef.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_update_table_expiration]
	return nil

}

func tableLabels(client *bigquery.Client, w io.Writer, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_get_table_labels]
	meta, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Table %s labels:\n", datasetID)
	if len(meta.Labels) == 0 {
		fmt.Println("Table has no labels defined.")
		return nil
	}
	for k, v := range meta.Labels {
		fmt.Fprintf(w, "\t%s:%s\n", k, v)
	}
	// [END bigquery_get_table_labels]
	return nil
}

func addTableLabel(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_label_table]
	tbl := client.Dataset(datasetID).Table(tableID)
	meta, err := tbl.Metadata(ctx)
	if err != nil {
		return err
	}

	update := bigquery.TableMetadataToUpdate{}
	update.SetLabel("color", "green")
	if _, err := tbl.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_label_table]
	return nil
}

func deleteTableLabel(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_delete_label_table]
	tbl := client.Dataset(datasetID).Table(tableID)
	meta, err := tbl.Metadata(ctx)
	if err != nil {
		return err
	}
	update := bigquery.TableMetadataToUpdate{}
	update.DeleteLabel("color")
	if _, err := tbl.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	// [END bigquery_delete_label_table]
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

func queryBasic(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query]

	q := client.Query(
		"SELECT name FROM `bigquery-public-data.usa_names.usa_1910_2013` " +
			"WHERE state = \"TX\" " +
			"LIMIT 100")
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// [END bigquery_query]
	return runAndRead(ctx, client, q)
}

func queryDisableCache(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_no_cache]

	q := client.Query(
		"SELECT corpus FROM `bigquery-public-data.samples.shakespeare` GROUP BY corpus;")
	q.DisableQueryCache = true
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// [END bigquery_query_no_cache]

	return runAndRead(ctx, client, q)
}

func queryBatch(client *bigquery.Client, dstDatasetID, dstTableID string) error {
	ctx := context.Background()
	// [START bigquery_query_batch]
	// Build an aggregate table.
	q := client.Query(`
		SELECT
  			corpus,
  			SUM(word_count) as total_words,
  			COUNT(1) as unique_words
		FROM ` + "`bigquery-public-data.samples.shakespeare`" + `
		GROUP BY corpus;`)
	q.Priority = bigquery.BatchPriority
	q.QueryConfig.Dst = client.Dataset(dstDatasetID).Table(dstTableID)

	// Start the job.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	// Job is started and will progress without interaction.
	// To simulate other work being done, sleep a few seconds.
	time.Sleep(5 * time.Second)
	status, err := job.Status(ctx)
	if err != nil {
		return err
	}

	state := "Unknown"
	switch status.State {
	case bigquery.Pending:
		state = "Pending"
	case bigquery.Running:
		state = "Running"
	case bigquery.Done:
		state = "Done"
	}
	// You can continue to monitor job progress until it reaches
	// the Done state by polling periodically.  In this example,
	// we print the latest status.
	fmt.Printf("Job %s in Location %s currently in state: %s\n", job.ID(), job.Location(), state)

	// [END bigquery_query_batch]
	job.Cancel(ctx)
	return nil
}

func queryDryRun(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_dry_run]
	q := client.Query(`
	SELECT
		name,
		COUNT(*) as name_count
	FROM ` + "`bigquery-public-data.usa_names.usa_1910_2013`" + `
	WHERE state = 'WA'
	GROUP BY name`)
	q.DryRun = true
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	// Dry run is not asynchronous, so get the latest status and statistics.
	status := job.LastStatus()
	if err != nil {
		return err
	}
	fmt.Printf("This query will process %d bytes\n", status.Statistics.TotalBytesProcessed)
	// [END bigquery_query_dry_run]
	return nil
}

func queryWithDestination(client *bigquery.Client, destDatasetID, destTableID string) error {
	ctx := context.Background()
	// [START bigquery_query_destination_table]

	q := client.Query("SELECT 17 as my_col")
	q.Location = "US" // Location must match the dataset(s) referenced in query.
	q.QueryConfig.Dst = client.Dataset(destDatasetID).Table(destTableID)
	// [END bigquery_query_destination_table]
	return runAndRead(ctx, client, q)
}

func queryLegacy(client *bigquery.Client, sqlString string) error {
	ctx := context.Background()
	// [START bigquery_query_legacy]
	q := client.Query(sqlString)
	q.UseLegacySQL = true

	// [END bigquery_query_legacy]
	return runAndRead(ctx, client, q)
}

func queryLegacyLargeResults(client *bigquery.Client, dstDatasetID, dstTableID string) error {
	ctx := context.Background()
	// [START bigquery_query_legacy_large_results]
	q := client.Query(
		"SELECT corpus FROM [bigquery-public-data:samples.shakespeare] GROUP BY corpus;")
	q.UseLegacySQL = true
	q.AllowLargeResults = true
	q.QueryConfig.Dst = client.Dataset(dstDatasetID).Table(dstTableID)
	// [END bigquery_query_legacy_large_results]
	return runAndRead(ctx, client, q)
}

func queryWithArrayParams(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_params_arrays]
	q := client.Query(
		`SELECT
			name,
			sum(number) as count 
        FROM ` + "`bigquery-public-data.usa_names.usa_1910_2013`" + `
		WHERE
			gender = @gender
        	AND state IN UNNEST(@states)
		GROUP BY
			name
		ORDER BY
			count DESC
		LIMIT 10;`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "gender",
			Value: "M",
		},
		{
			Name:  "states",
			Value: []string{"WA", "WI", "WV", "WY"},
		},
	}
	// [END bigquery_query_params_arrays]
	return runAndRead(ctx, client, q)
}

func queryWithNamedParams(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_params_named]
	q := client.Query(
		`SELECT word, word_count
        FROM ` + "`bigquery-public-data.samples.shakespeare`" + `
        WHERE corpus = @corpus
        AND word_count >= @min_word_count
        ORDER BY word_count DESC;`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "corpus",
			Value: "romeoandjuliet",
		},
		{
			Name:  "min_word_count",
			Value: 250,
		},
	}
	// [END bigquery_query_params_named]
	return runAndRead(ctx, client, q)
}

func queryWithPositionalParams(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_params_positional]
	q := client.Query(
		`SELECT word, word_count
        FROM ` + "`bigquery-public-data.samples.shakespeare`" + `
        WHERE corpus = ?
        AND word_count >= ?
        ORDER BY word_count DESC;`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Value: "romeoandjuliet",
		},
		{
			Value: 250,
		},
	}
	// [END bigquery_query_params_positional]
	return runAndRead(ctx, client, q)
}

func queryWithTimestampParam(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_params_timestamps]
	q := client.Query(
		`SELECT TIMESTAMP_ADD(@ts_value, INTERVAL 1 HOUR);`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "ts_value",
			Value: time.Date(2016, 12, 7, 8, 0, 0, 0, time.UTC),
		},
	}
	// [END bigquery_query_params_timestamps]
	return runAndRead(ctx, client, q)
}

func queryWithStructParam(client *bigquery.Client) error {
	ctx := context.Background()
	// [START bigquery_query_params_structs]
	type MyStruct struct {
		X int64
		Y string
	}
	q := client.Query(
		`SELECT @struct_value as s;`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "struct_value",
			Value: MyStruct{X: 1, Y: "foo"},
		},
	}
	// [END bigquery_query_params_structs]
	return runAndRead(ctx, client, q)
}

func printTableInfo(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_get_table]
	meta, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return err
	}
	// Print basic information about the table.
	fmt.Printf("Schema has %d top-level fields\n", len(meta.Schema))
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Printf("Rows in managed storage: %d\n", meta.NumRows)
	// [END bigquery_get_table]
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

// generateTableCTAS creates a quick table by issuing a CREATE TABLE AS SELECT
// query.
func generateTableCTAS(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	q := client.Query(
		fmt.Sprintf(
			`CREATE TABLE %s.%s 
		AS
		SELECT
		  2000 + CAST(18 * RAND() as INT64) as year,
		  IF(RAND() > 0.5,"foo","bar") as token
		FROM
		  UNNEST(GENERATE_ARRAY(0,5,1)) as r`, datasetID, tableID))
	job, err := q.Run(ctx)
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
	return nil
}

func copyMultiTable(client *bigquery.Client, datasetID, dstTableID string) error {
	ctx := context.Background()
	// Generate some dummy tables via a quick CTAS.
	if err := generateTableCTAS(client, datasetID, "table1"); err != nil {
		return err
	}
	if err := generateTableCTAS(client, datasetID, "table2"); err != nil {
		return err
	}
	// [START bigquery_copy_table_multiple_source]
	dataset := client.Dataset(datasetID)

	srcTableIDs := []string{"table1", "table2"}
	var tableRefs []*bigquery.Table
	for _, v := range srcTableIDs {
		tableRefs = append(tableRefs, dataset.Table(v))
	}
	copier := dataset.Table(dstTableID).CopierFrom(tableRefs...)
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
	// [END bigquery_copy_table_multiple_source]
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

func deleteAndUndeleteTable(client *bigquery.Client, datasetID, tableID string) error {
	ctx := context.Background()
	// [START bigquery_undelete_table]

	ds := client.Dataset(datasetID)
	if _, err := ds.Table(tableID).Metadata(ctx); err != nil {
		return err
	}
	// Record the current time.  We'll use this as the snapshot time
	// for recovering the table.
	snapTime := time.Now()

	// "Accidentally" delete the table.
	if err := client.Dataset(datasetID).Table(tableID).Delete(ctx); err != nil {
		return err
	}

	// Construct the restore-from tableID using a snapshot decorator.
	snapshotTableID := fmt.Sprintf("%s@%d", tableID, snapTime.UnixNano()/1e6)
	// Choose a new table ID for the recovered table data.
	recoverTableID := fmt.Sprintf("%s_recovered", tableID)

	// Construct and run a copy job.
	copier := ds.Table(recoverTableID).CopierFrom(ds.Table(snapshotTableID))
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

	// [END bigquery_undelete_table]
	ds.Table(recoverTableID).Delete(ctx)
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
	source.AutoDetect = true   // Allow BigQuery to determine schema.
	source.SkipLeadingRows = 1 // CSV has a single header line.

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
	// You can choose to run the job in a specific location for more complex data locality scenarios.
	// Ex: In this example, source dataset and GCS bucket are in the US.
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
	// You can choose to run the job in a specific location for more complex data locality scenarios.
	// Ex: In this example, source dataset and GCS bucket are in the US.
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
	// You can choose to run the job in a specific location for more complex data locality scenarios.
	// Ex: In this example, source dataset and GCS bucket are in the US.
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

// runAndRead executes a query then prints results.
func runAndRead(ctx context.Context, client *bigquery.Client, q *bigquery.Query) error {
	// [START bigquery_query]
	// [START bigquery_query_destination_table]
	// [START bigquery_query_legacy]
	// [START bigquery_query_legacy_large_results]
	// [START bigquery_query_no_cache]
	// [START bigquery_query_params_arrays]
	// [START bigquery_query_params_named]
	// [START bigquery_query_params_positional]
	// [START bigquery_query_params_timestamps]
	// [START bigquery_query_params_structs]
	job, err := q.Run(ctx)
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
	// [END bigquery_query_destination_table]
	// [END bigquery_query_legacy]
	// [END bigquery_query_legacy_large_results]
	// [END bigquery_query_no_cache]
	// [END bigquery_query_params_arrays]
	// [END bigquery_query_params_named]
	// [END bigquery_query_params_positional]
	// [END bigquery_query_params_timestamps]
	// [END bigquery_query_params_structs]
	return nil
}
