// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
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
	empty := uniqueBQName("golang_example_table_emptyschema")

	if err := createTableInferredSchema(client, datasetID, inferred); err != nil {
		t.Errorf("createTableInferredSchema(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := createTableExplicitSchema(client, datasetID, explicit); err != nil {
		t.Errorf("createTableExplicitSchema(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	if err := createTableEmptySchema(client, datasetID, empty); err != nil {
		t.Errorf("createTableEmptySchema(dataset:%q table:%q): %v", datasetID, empty, err)
	}

	if err := updateTableDescription(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableDescription(dataset:%q table:%q): %v", datasetID, explicit, err)
	}
	if err := updateTableExpiration(client, datasetID, explicit); err != nil {
		t.Errorf("updateTableExpiration(dataset:%q table:%q): %v", datasetID, explicit, err)
	}

	buf := &bytes.Buffer{}
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
	if got := buf.String(); !strings.Contains(got, empty) {
		t.Errorf("want table list %q to contain table %q", got, empty)
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
