// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func TestAll(t *testing.T) {
	tc := testutil.SystemTest(t)

	client, err := bigquery.NewClient(context.Background(), tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	datasetID := fmt.Sprintf("golang_example_dataset_%d", time.Now().Unix())
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	if err := listDatasets(client); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}

	tableID := fmt.Sprintf("golang_example_table_%d", time.Now().Unix())
	if err := createTable(client, datasetID, tableID); err != nil {
		t.Errorf("failed to create table: %v", err)
	}
	if err := insertRows(client, datasetID, tableID); err != nil {
		t.Errorf("failed to insert rows: %v", err)
	}

	dstTableID := fmt.Sprintf("golang_example_tabledst_%d", time.Now().Unix())
	if err := copyTable(client, datasetID, tableID, dstTableID); err != nil {
		t.Errorf("failed to copy table: %v", err)
	}
	if err := deleteTable(client, datasetID, tableID); err != nil {
		t.Errorf("failed to delete table: %v", err)
	}
	if err := deleteTable(client, datasetID, dstTableID); err != nil {
		t.Errorf("failed to delete table: %v", err)
	}
}

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
	// [ENDSTART bigquery_insert_stream]
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

	for {
		status, err := job.Status(ctx)
		if err != nil {
			return err
		}
		if status.Done() {
			if status.Err() != nil {
				return status.Err()
			}
			break
		}
		time.Sleep(pollInterval)
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

const pollInterval = time.Second
