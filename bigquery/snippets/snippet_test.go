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

	// test empty dataset creation/ttl/delete
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

	tableID := fmt.Sprintf("golang_example_table_%d", time.Now().Unix())
	if err := createTable(client, datasetID, tableID); err != nil {
		t.Errorf("createTable(dataset:%q  table:%q): %v", datasetID, tableID, err)
	}
	buf := &bytes.Buffer{}
	if err := listTables(client, buf, datasetID); err != nil {
		t.Errorf("listTables(%q): %v", datasetID, err)
	}
	if got := buf.String(); !strings.Contains(got, tableID) {
		t.Errorf("want table list %q to contain table %q", got, tableID)
	}
	if err := insertRows(client, datasetID, tableID); err != nil {
		t.Errorf("insertRows(dataset:%q table:%q): %v", datasetID, tableID, err)
	}
	if err := listRows(client, datasetID, tableID); err != nil {
		t.Errorf("listRows(dataset:%q table:%q): %v", datasetID, tableID, err)
	}
	if err := browseTable(client, datasetID, tableID); err != nil {
		t.Errorf("browseTable(dataset:%q table:%q): %v", datasetID, tableID, err)
	}
	if err := asyncQuery(client, datasetID, tableID); err != nil {
		t.Errorf("failed to async query: %v", err)
	}

	dstTableID := fmt.Sprintf("golang_example_tabledst_%d", time.Now().Unix())
	if err := copyTable(client, datasetID, tableID, dstTableID); err != nil {
		t.Errorf("failed to copy table (dataset:%q src:%q dst:%q): %v", datasetID, tableID, dstTableID, err)
	}
	if err := deleteTable(client, datasetID, tableID); err != nil {
		t.Errorf("deleteTable(dataset:%q table:%q): %v", datasetID, tableID, err)
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
		t.Errorf("failed to create dataset: %v", err)
	}
	schema := bigquery.Schema{
		&bigquery.FieldSchema{Name: "Year", Type: bigquery.IntegerFieldType},
		&bigquery.FieldSchema{Name: "City", Type: bigquery.StringFieldType},
	}
	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, &bigquery.TableMetadata{
		Schema: schema,
	}); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	defer deleteDataset(t, ctx, datasetID)

	if err := importFromFile(client, datasetID, tableID, "testdata/olympics.csv"); err != nil {
		t.Fatalf("failed to import from file: %v", err)
	}

	jsonTableExplicit := fmt.Sprintf("golang_example_dataset_importjson_explicit_%d", time.Now().Unix())
	if err := importJSONExplicitSchema(client, datasetID, jsonTableExplicit); err != nil {
		t.Fatalf("Failed to ingest JSON sample with explicit schema: %v", err)
	}

	jsonTableAutodetect := fmt.Sprintf("golang_example_dataset_importjson_autodetect_%d", time.Now().Unix())
	if err := importJSONAutodetectSchema(client, datasetID, jsonTableAutodetect); err != nil {
		t.Fatalf("Failed to ingest JSON sample with explicit schema: %v", err)
	}

	bucket := fmt.Sprintf("golang-example-bigquery-importexport-bucket-%d", time.Now().Unix())
	const object = "values.csv"

	if err := storageClient.Bucket(bucket).Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatalf("cannot create bucket: %v", err)
	}

	gcsURI := fmt.Sprintf("gs://%s/%s", bucket, object)
	if err := exportToGCS(client, datasetID, tableID, gcsURI); err != nil {
		t.Errorf("failed to export to %v: %v", gcsURI, err)
	}

	// Cleanup the bucket and object.
	if err := storageClient.Bucket(bucket).Object(object).Delete(ctx); err != nil {
		t.Errorf("failed to cleanup the GCS object: %v", err)
	}
	time.Sleep(time.Second) // Give it a second, due to eventual consistency.
	if err := storageClient.Bucket(bucket).Delete(ctx); err != nil {
		t.Errorf("failed to cleanup the GCS bucket: %v", err)
	}
}
