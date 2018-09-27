// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

// uniqueBQName returns a more unique name for a BigQuery resource.
func uniqueBQName(prefix string) string {
	t := time.Now()
	return fmt.Sprintf("%s_%d", sanitize(prefix, '_'), t.Unix())
}

// uniqueBucketName returns a more unique name cloud storage bucket.
func uniqueBucketName(prefix, projectID string) string {
	t := time.Now()
	f := fmt.Sprintf("%s-%s-%d", sanitize(prefix, '-'), sanitize(projectID, '-'), t.Unix())
	// bucket max name length is 63 chars, so we truncate.
	if len(f) > 63 {
		return f[:63]
	}
	return f
}

func sanitize(s string, allowedSeparator rune) string {
	pattern := fmt.Sprintf("[^a-zA-Z0-9%s]", string(allowedSeparator))
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return reg.ReplaceAllString(s, "")
}
func TestAll(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := uniqueBQName("golang_example_dataset")
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(datasetID).DeleteWithContents(ctx)

	if err := updateDatasetAccessControl(client, datasetID); err != nil {
		t.Errorf("updateDataSetAccessControl(%q): %v", datasetID, err)
	}
	if err := addDatasetLabel(client, datasetID); err != nil {
		t.Errorf("updateDatasetAddLabel: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := datasetLabels(client, buf, datasetID); err != nil {
		t.Errorf("getDatasetLabels(%q): %v", datasetID, err)
	}
	want := "color:green"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getDatasetLabel(%q) expected %q to contain %q", datasetID, got, want)
	}

	if err := addDatasetLabel(client, datasetID); err != nil {
		t.Errorf("updateDatasetAddLabel: %v", err)
	}
	buf.Reset()
	if err := listDatasetsByLabel(client, buf); err != nil {
		t.Errorf("listDatasetsByLabel: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, datasetID) {
		t.Errorf("listDatasetsByLabel expected %q to contain %q", got, want)
	}
	if err := deleteDatasetLabel(client, datasetID); err != nil {
		t.Errorf("updateDatasetDeleteLabel: %v", err)
	}

	// Test empty dataset creation/ttl/delete.
	deletionDatasetID := uniqueBQName("golang_example_quickdelete")
	if err := createDataset(client, deletionDatasetID); err != nil {
		t.Errorf("createDataset(%q): %v", deletionDatasetID, err)
	}
	if err = updateDatasetDefaultExpiration(client, deletionDatasetID); err != nil {
		t.Errorf("updateDatasetDefaultExpiration(%q): %v", deletionDatasetID, err)
	}
	if err := deleteEmptyDataset(client, deletionDatasetID); err != nil {
		t.Errorf("deleteEmptyDataset(%q): %v", deletionDatasetID, err)
	}

	if err := updateDatasetDescription(client, datasetID); err != nil {
		t.Errorf("updateDatasetDescription(%q): %v", datasetID, err)
	}
	if err := listDatasets(client); err != nil {
		t.Errorf("listDatasets: %v", err)
	}

	inferred := uniqueBQName("golang_example_table_inferred")
	explicit := uniqueBQName("golang_example_table_explicit")

	if err := createTableInferredSchema(client, datasetID, inferred); err != nil {
		t.Errorf("createTableInferredSchema(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := createTableExplicitSchema(client, datasetID, explicit); err != nil {
		t.Errorf("createTableExplicitSchema(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	complex := uniqueBQName("golang_example_table_complex")
	if err := createTableComplexSchema(client, datasetID, complex); err != nil {
		t.Errorf("createTableComplexSchema(dataset:%q table:%q): %v", datasetID, complex, err)
	}

	tableCMEK := uniqueBQName("golang_example_table_cmek")
	if err := createTableWithCMEK(client, datasetID, tableCMEK); err != nil {
		t.Errorf("createTableWithCMEK(dataset:%q table:%q): %v", datasetID, tableCMEK, err)
	}

	required := uniqueBQName("golang_example_table_required")
	if err := relaxTableAPI(client, datasetID, required); err != nil {
		t.Errorf("relaxTableApi(dataset:%q table:%q): %v", datasetID, required, err)
	}

	widenLoad := uniqueBQName("golang_example_table_widen_load")
	filenameWiden := "testdata/people.csv"
	if err := createTableAndWidenLoad(client, datasetID, widenLoad, filenameWiden); err != nil {
		t.Errorf("createTableAndWidenLoad(dataset:%q table:%q): %v", datasetID, widenLoad, err)
	}

	widenQuery := uniqueBQName("golang_example_table_widen_query")
	if err := createTableAndWidenQuery(client, datasetID, widenQuery); err != nil {
		t.Errorf("createTableAndWidenQuery(dataset:%q table:%q): %v", datasetID, widenQuery, err)
	}

	if err := updateTableDescription(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableDescription(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	if err := updateTableExpiration(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableExpiration(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	if err := updateTableAddColumn(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableAddColumn(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	if err := addTableLabel(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableAddLabel(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	if err := deleteTableLabel(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableAddLabel(dataset:%q table:%q): %v", datasetID, explicit, err)
	}

	buf.Reset()
	if err := listTables(client, buf, datasetID); err != nil {
		t.Errorf("listTables(%q): %v", datasetID, err)
	}
	// Ensure all three tables are in the list.
	if got := buf.String(); !strings.Contains(got, inferred) {
		t.Errorf("want table list %q to contain table %q", got, inferred)
	}
	if got := buf.String(); !strings.Contains(got, explicit) {
		t.Errorf("want table list %q to contain table %q", got, explicit)
	}

	if err := printDatasetInfo(client, datasetID); err != nil {
		t.Errorf("printDatasetInfo: %v", err)
	}

	// Stream data, read, query the inferred schema table.
	if err := insertRows(client, datasetID, inferred); err != nil {
		t.Errorf("insertRows(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := browseTable(client, datasetID, inferred); err != nil {
		t.Errorf("browseTable(dataset:%q table:%q): %v", datasetID, inferred, err)
	}

	if err := queryBasic(client); err != nil {
		t.Errorf("queryBasic: %v", err)
	}
	batchTable := uniqueBQName("golang_example_batchresults")
	if err := queryBatch(client, datasetID, batchTable); err != nil {
		t.Errorf("queryBatch(dataset:%q table:%q): %v", datasetID, batchTable, err)
	}
	if err := queryDisableCache(client); err != nil {
		t.Errorf("queryBasicDisableCache: %v", err)
	}
	if err := queryDryRun(client); err != nil {
		t.Errorf("queryDryRun: %v", err)
	}
	sql := "SELECT 17 as foo"
	if err := queryLegacy(client, sql); err != nil {
		t.Errorf("queryLegacy: %v", err)
	}
	largeResults := uniqueBQName("golang_example_legacy_largeresults")
	if err := queryLegacyLargeResults(client, datasetID, largeResults); err != nil {
		t.Errorf("queryLegacyLargeResults(dataset:%q table:%q): %v", datasetID, largeResults, err)
	}
	if err := queryWithArrayParams(client); err != nil {
		t.Errorf("queryWithArrayParams: %v", err)
	}
	if err := queryWithNamedParams(client); err != nil {
		t.Errorf("queryWithNamedParams: %v", err)
	}
	if err := queryWithPositionalParams(client); err != nil {
		t.Errorf("queryWithPositionalParams: %v", err)
	}
	if err := queryWithTimestampParam(client); err != nil {
		t.Errorf("queryWithTimestampParam: %v", err)
	}
	if err := queryWithStructParam(client); err != nil {
		t.Errorf("queryWithStructParam: %v", err)
	}

	// Run query variations
	persisted := uniqueBQName("golang_example_table_queryresult")
	if err := queryWithDestination(client, datasetID, persisted); err != nil {
		t.Errorf("queryWithDestination(dataset:%q table:%q): %v", datasetID, persisted, err)
	}
	persistedCMEK := uniqueBQName("golang_example_table_queryresult_cmek")
	if err := queryWithDestinationCMEK(client, datasetID, persistedCMEK); err != nil {
		t.Errorf("queryWithDestinationCMEK(dataset:%q table:%q): %v", datasetID, persistedCMEK, err)
	}

	// Control a job lifecycle explicitly: create, report status, cancel.
	exampleJobID := uniqueBQName("golang_example_job")
	q := client.Query("Select 17 as foo")
	q.JobID = exampleJobID
	q.Priority = bigquery.BatchPriority
	q.Run(ctx)
	if err := getJobInfo(client, exampleJobID); err != nil {
		t.Errorf("getJobInfo(%s): %v", exampleJobID, err)
	}
	if err := cancelJob(client, exampleJobID); err != nil {
		t.Errorf("cancelJobInfo(%s): %v", exampleJobID, err)
	}

	// Print information about tables (extended and simple).
	if err := printTableInfo(client, datasetID, inferred); err != nil {
		t.Errorf("printTableInfo(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := printTableInfo(client, datasetID, explicit); err != nil {
		t.Errorf("printTableInfo(dataset:%q table:%q): %v", datasetID, explicit, err)
	}

	dstTableID := uniqueBQName("golang_example_tabledst")
	if err := copyTable(client, datasetID, inferred, dstTableID); err != nil {
		t.Errorf("copyTable(dataset:%q src:%q dst:%q): %v", datasetID, inferred, dstTableID, err)
	}
	if err := deleteTable(client, datasetID, inferred); err != nil {
		t.Errorf("deleteTable(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := deleteAndUndeleteTable(client, datasetID, dstTableID); err != nil {
		t.Errorf("undeleteTable(dataset:%q table:%q): %v", datasetID, dstTableID, err)
	}

	dstTableID = uniqueBQName("golang_multicopydest")
	if err := copyMultiTable(client, datasetID, dstTableID); err != nil {
		t.Errorf("copyMultiTable(dataset:%q table:%q): %v", datasetID, dstTableID, err)
	}
	dstTableID = uniqueBQName("golang_example_copycmek")
	if err := copyTableWithCMEK(client, datasetID, dstTableID); err != nil {
		t.Errorf("copyTableWithCMEK(dataset:%q table:%q): %v", datasetID, dstTableID, err)
	}

	if err := listJobs(client); err != nil {
		t.Errorf("listJobs: %v", err)
	}

}

// Exercise BigQuery logical views.
func TestViews(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	srcDatasetID := uniqueBQName("golang_example_view_source")
	if err := createDataset(client, srcDatasetID); err != nil {
		t.Errorf("createDataset(%q): %v", srcDatasetID, err)
	}
	defer client.Dataset(srcDatasetID).DeleteWithContents(ctx)
	viewDatasetID := uniqueBQName("golang_example_view_container")
	if err := createDataset(client, viewDatasetID); err != nil {
		t.Errorf("createDataset(%q): %v", viewDatasetID, err)
	}
	defer client.Dataset(viewDatasetID).DeleteWithContents(ctx)

	viewID := uniqueBQName("golang_example_view")

	if err := createView(client, viewDatasetID, viewID); err != nil {
		t.Fatalf("createView(dataset:%q view:%q): %v", viewDatasetID, viewID, err)
	}

	if err := getView(client, viewDatasetID, viewID); err != nil {
		t.Fatalf("getView(dataset:%q view:%q): %v", viewDatasetID, viewID, err)
	}

	if err := updateView(client, viewDatasetID, viewID); err != nil {
		t.Fatalf("updateView(dataset:%q view:%q): %v", viewDatasetID, viewID, err)
	}

	if err := updateViewDelegated(client, srcDatasetID, viewDatasetID, viewID); err != nil {
		t.Fatalf("updateViewDelegated(srcdataset:%q viewdataset:%q view:%q): %v", srcDatasetID, viewDatasetID, viewID, err)
	}

}

func TestImportExport(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := uniqueBQName("golang_example_dataset_importexport")
	tableID := uniqueBQName("golang_example_dataset_importexport")
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}
	defer client.Dataset(datasetID).DeleteWithContents(ctx)

	filename := "testdata/people.csv"
	if err := importCSVFromFile(client, datasetID, tableID, filename); err != nil {
		t.Fatalf("importCSVFromFile(dataset:%q table:%q filename:%q): %v", datasetID, tableID, filename, err)
	}

	autoCSV := uniqueBQName("golang_example_csv_autodetect")
	if err := importCSVAutodetectSchema(client, datasetID, autoCSV); err != nil {
		t.Fatalf("importCSVAutodetectSchema(dataset:%q table:%q): %v", datasetID, autoCSV, err)
	}

	autoCSVTruncate := uniqueBQName("golang_example_csv_truncate")
	if err := importCSVTruncate(client, datasetID, autoCSVTruncate); err != nil {
		t.Fatalf("importCSVTruncate(dataset:%q table:%q): %v", datasetID, autoCSVTruncate, err)
	}

	explicitCSV := uniqueBQName("golang_example_dataset_importcsv_explicit")
	if err := importCSVExplicitSchema(client, datasetID, explicitCSV); err != nil {
		t.Fatalf("importCSVExplicitSchema(dataset:%q table:%q): %v", datasetID, explicitCSV, err)
	}

	explicitJSON := uniqueBQName("golang_example_dataset_importjson_explicit")
	if err := importJSONExplicitSchema(client, datasetID, explicitJSON); err != nil {
		t.Fatalf("importJSONExplicitSchema(dataset:%q table:%q): %v", datasetID, explicitJSON, err)
	}

	autodetectJSON := uniqueBQName("golang_example_dataset_importjson_autodetect")
	if err := importJSONAutodetectSchema(client, datasetID, autodetectJSON); err != nil {
		t.Fatalf("importJSONAutodetectSchema(dataset:%q table:%q): %v", datasetID, autodetectJSON, err)
	}

	autoJSONwithCMEK := uniqueBQName("golang_example_importjson_cmek")
	if err := importJSONWithCMEK(client, datasetID, autoJSONwithCMEK); err != nil {
		t.Fatalf("importJSONWithCMEK(dataset:%q table:%q): %v", datasetID, autoJSONwithCMEK, err)
	}

	autoJSONTruncate := uniqueBQName("golang_example_importjson_truncate")
	if err := importJSONTruncate(client, datasetID, autoJSONTruncate); err != nil {
		t.Fatalf("importJSONTruncate(dataset:%q table:%q): %v", datasetID, autoJSONTruncate, err)
	}

	orc := uniqueBQName("golang_example_importorc")
	if err := importORC(client, datasetID, orc); err != nil {
		t.Errorf("importOrc(dataset:%q table: %q): %v", datasetID, orc, err)
	}
	if err := importORCTruncate(client, datasetID, orc); err != nil {
		t.Errorf("importOrcTruncate(dataset:%q table: %q): %v", datasetID, orc, err)
	}

	parquet := uniqueBQName("golang_example_importparquet")
	if err := importParquet(client, datasetID, parquet); err != nil {
		t.Errorf("importParquet(dataset:%q table: %q): %v", datasetID, parquet, err)
	}
	if err := importParquetTruncate(client, datasetID, parquet); err != nil {
		t.Errorf("importParquetTruncate(dataset:%q table: %q): %v", datasetID, parquet, err)
	}

	requiredImport := uniqueBQName("golang_example_table_required_import")
	filenameRelax := "testdata/people.csv"
	if err := relaxTableImport(client, datasetID, requiredImport, filenameRelax); err != nil {
		t.Errorf("relaxTableImport(dataset:%q table:%q): %v", datasetID, requiredImport, err)
	}

	requiredQuery := uniqueBQName("golang_example_table_required_query")
	if err := relaxTableQuery(client, datasetID, requiredQuery); err != nil {
		t.Errorf("relaxTableQuery(dataset:%q table:%q): %v", datasetID, requiredImport, err)
	}

	bucket := uniqueBucketName("golang-example-bucket", tc.ProjectID)
	const object = "values.csv"

	if err := storageClient.Bucket(bucket).Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatalf("cannot create bucket: %v", err)
	}

	// extract shakespeare sample as CSV
	gcsURI := fmt.Sprintf("gs://%s/%s", bucket, "shakespeare.csv")
	if err := exportSampleTableAsCSV(client, gcsURI); err != nil {
		t.Errorf("exportSampleTableAsCSV(%q): %v", gcsURI, err)
	}

	// extract shakespeare sample as GZIP-compressed CSV
	gcsURI = fmt.Sprintf("gs://%s/%s", bucket, "shakespeare.csv.gz")
	if err := exportSampleTableAsCompressedCSV(client, gcsURI); err != nil {
		t.Errorf("exportSampleTableAsCompressedCSV(%q): %v", gcsURI, err)
	}

	// extract shakespeare sample as newline-delimited JSON
	gcsURI = fmt.Sprintf("gs://%s/%s", bucket, "shakespeare.json")
	if err := exportSampleTableAsJSON(client, gcsURI); err != nil {
		t.Errorf("exportSampleTableAsJSON(%q): %v", gcsURI, err)
	}

	// Walk the bucket and delete objects
	it := storageClient.Bucket(bucket).Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err := storageClient.Bucket(bucket).Object(objAttrs.Name).Delete(ctx); err != nil {
			t.Errorf("failed to cleanup the GCS object: %v", err)
		}
	}

	time.Sleep(time.Second) // Give it a second, due to eventual consistency.
	if err := storageClient.Bucket(bucket).Delete(ctx); err != nil {
		t.Errorf("failed to cleanup the GCS bucket: %v", err)
	}
}

func TestPartitioningAndClustering(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := uniqueBQName("golang_example_dataset_partition_cluster")
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}
	defer client.Dataset(datasetID).DeleteWithContents(ctx)

	partitionedEmpty := uniqueBQName("golang_example_partitioned")
	if err := createTablePartitioned(client, datasetID, partitionedEmpty); err != nil {
		t.Errorf("createTablePartitioned(dataset:%q table:%q): %v", datasetID, partitionedEmpty, err)
	}

	partitionedLoad := uniqueBQName("golang_example_partitioned_load")
	if err := importPartitionedSampleTable(client, datasetID, partitionedLoad); err != nil {
		t.Errorf("importPartitionedStatesByDate(dataset:%q table:%q): %v", datasetID, partitionedLoad, err)
	}

	if err := queryPartitionedTable(client, datasetID, partitionedLoad); err != nil {
		t.Errorf("queryPartitionedTable(dataset:%q table:%q): %v", datasetID, partitionedLoad, err)
	}

	clusteredEmpty := uniqueBQName("golang_example_clustered")
	if err := createTableClustered(client, datasetID, clusteredEmpty); err != nil {
		t.Errorf("createTableClustered(dataset:%q table:%q): %v", datasetID, clusteredEmpty, err)
	}

	clusteredLoad := uniqueBQName("golang_example_clustered_transactions")
	if err := importClusteredSampleTable(client, datasetID, clusteredLoad); err != nil {
		t.Errorf("importClusteredSampleTable(dataset:%q table:%q): %v", datasetID, clusteredLoad, err)
	}

	if err := queryClusteredTable(client, datasetID, clusteredLoad); err != nil {
		t.Errorf("queryClusteredTable(dataset:%q table:%q): %v", datasetID, clusteredLoad, err)
	}
}
