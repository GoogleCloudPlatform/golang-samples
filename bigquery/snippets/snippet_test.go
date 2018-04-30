// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	rawbq "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/iterator"
)

func init() {
	// Workaround for Travis:
	// https://docs.travis-ci.com/user/common-build-problems/#Build-times-out-because-no-output-was-received
	if os.Getenv("TRAVIS") == "true" {
		go func() {
			for {
				time.Sleep(5 * time.Minute)
				log.Print("Still testing. Don't kill me!")
			}
		}()
	}
}

func TestAll(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := fmt.Sprintf("golang_example_dataset_%d", time.Now().Unix())
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}

	if err := updateDatasetAccessControl(client, datasetID); err != nil {
		t.Errorf("updateDataSetAccessControl(%q): %v", datasetID, err)
	}

	// Test empty dataset creation/ttl/delete.
	deletionDatasetID := fmt.Sprintf("%s_quickdelete", datasetID)
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

	inferred := fmt.Sprintf("golang_example_table_inferred_%d", time.Now().Unix())
	explicit := fmt.Sprintf("golang_example_table_explicit_%d", time.Now().Unix())
	empty := fmt.Sprintf("golang_example_table_emptyschema_%d", time.Now().Unix())

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

	// Stream data, read, query the inferred schema table.
	if err := insertRows(client, datasetID, inferred); err != nil {
		t.Errorf("insertRows(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := listRows(client, datasetID, inferred); err != nil {
		t.Errorf("listRows(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := browseTable(client, datasetID, inferred); err != nil {
		t.Errorf("browseTable(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := basicQuery(client, datasetID, inferred); err != nil {
		t.Errorf("basicQuery(dataset:%q table:%q): %v", datasetID, inferred, err)
	}

	// Print information about tables (extended and simple).
	if err := printTableMetadataSimple(client, datasetID, inferred); err != nil {
		t.Errorf("printTableMetadata(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := printTableMetadataSimple(client, datasetID, explicit); err != nil {
		t.Errorf("printTableMetadata(dataset:%q table:%q): %v", datasetID, explicit, err)
	}

	dstTableID := fmt.Sprintf("golang_example_tabledst_%d", time.Now().Unix())
	if err := copyTable(client, datasetID, inferred, dstTableID); err != nil {
		t.Errorf("copyTable(dataset:%q src:%q dst:%q): %v", datasetID, inferred, dstTableID, err)
	}
	if err := deleteTable(client, datasetID, inferred); err != nil {
		t.Errorf("deleteTable(dataset:%q table:%q): %v", datasetID, inferred, err)
	}
	if err := deleteTable(client, datasetID, dstTableID); err != nil {
		t.Errorf("deleteTable(dataset:%q table:%q): %v", datasetID, dstTableID, err)
	}

	deleteDataset(t, ctx, datasetID)
}

func deleteDataset(t *testing.T, ctx context.Context, datasetID string) {
	tc := testutil.SystemTest(t)
	hc, err := google.DefaultClient(ctx, rawbq.CloudPlatformScope)
	if err != nil {
		t.Errorf("DefaultClient: %v", err)
	}
	s, err := rawbq.New(hc)
	if err != nil {
		t.Errorf("bigquery.New: %v", err)
	}
	call := s.Datasets.Delete(tc.ProjectID, datasetID)
	call.DeleteContents(true)
	call.Context(ctx)
	if err := call.Do(); err != nil {
		t.Errorf("deleteDataset(%q): %v", datasetID, err)
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

	datasetID := fmt.Sprintf("golang_example_dataset_importexport_%d", time.Now().Unix())
	tableID := fmt.Sprintf("golang_example_dataset_importexport_%d", time.Now().Unix())
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}
	defer deleteDataset(t, ctx, datasetID)

	filename := "testdata/people.csv"
	if err := importCSVFromFile(client, datasetID, tableID, filename); err != nil {
		t.Fatalf("importCSVFromFile(dataset:%q table:%q filename:%q): %v", datasetID, tableID, filename, err)
	}

	explicitCSV := fmt.Sprintf("golang_example_dataset_importcsv_explicit_%d", time.Now().Unix())
	if err := importCSVExplicitSchema(client, datasetID, explicitCSV); err != nil {
		t.Fatalf("importCSVExplicitSchema(dataset:%q table:%q): %v", datasetID, explicitCSV, err)
	}

	explicitJSON := fmt.Sprintf("golang_example_dataset_importjson_explicit_%d", time.Now().Unix())
	if err := importJSONExplicitSchema(client, datasetID, explicitJSON); err != nil {
		t.Fatalf("importJSONExplicitSchema(dataset:%q table:%q): %v", datasetID, explicitJSON, err)
	}

	autodetectJSON := fmt.Sprintf("golang_example_dataset_importjson_autodetect_%d", time.Now().Unix())
	if err := importJSONAutodetectSchema(client, datasetID, autodetectJSON); err != nil {
		t.Fatalf("importJSONAutodetectSchema(dataset:%q table:%q): %v", datasetID, autodetectJSON, err)
	}
	bucket := fmt.Sprintf("golang-example-bigquery-importexport-bucket-%d", time.Now().Unix())
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
