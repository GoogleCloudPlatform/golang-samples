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

// Package risk contains example snippets using the DLP API to create risk jobs.
package risk

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

const (
	riskTopicName        = "dlp-risk-test-topic-"
	riskSubscriptionName = "dlp-risk-test-sub-"
)

func cleanupPubsub(t *testing.T, client *pubsub.Client, topicName, subName string) {
	ctx := context.Background()
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		t.Logf("Exists: %v", err)
		return
	} else if exists {
		if err := topic.Delete(ctx); err != nil {
			t.Logf("Delete: %v", err)
		}
	}

	s := client.Subscription(subName)
	if exists, err := s.Exists(ctx); err != nil {
		t.Logf("Exists: %v", err)
		return
	} else if exists {
		if err := s.Delete(ctx); err != nil {
			t.Logf("Delete: %v", err)
		}
	}
}

var (
	u         = uuid.New().String()[:8]
	projectID string
	tableID   = fmt.Sprint("dlp_test_risk_table" + u)
	dataSetID = fmt.Sprint("dlp_test_risk_dataset" + u)
)

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	projectID = tc.ProjectID
	if !ok {
		log.Fatal("couldn't initialize test")
		return
	}

	createBigQueryDataSetId(tc.ProjectID, dataSetID)
	createTableInsideDataset(tc.ProjectID, dataSetID, tableID)
	m.Run()
	deleteBigQueryAssets(tc.ProjectID, dataSetID)

}

func createBigQueryDataSetId(projectID, dataSetID string) error {

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}

	if err := client.Dataset(dataSetID).Create(ctx, meta); err != nil {
		return err
	}

	return nil
}

func createTableInsideDataset(projectID, dataSetID, tableID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "user_id", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
		{Name: "title", Type: bigquery.StringFieldType},
		{Name: "score", Type: bigquery.StringFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0), // Table will be automatically deleted in 1 year.
	}

	tableRef := client.Dataset(dataSetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		log.Printf("[INFO] createBigQueryDataSetId Error while table creation: %v", err)
		return err
	}

	duration := time.Duration(90) * time.Second
	time.Sleep(duration)

	inserter := client.Dataset(dataSetID).Table(tableID).Inserter()
	items := []*BigQueryTableItem{
		// Item implements the ValueSaver interface.
		{UserId: "602-61-8588", Age: 32, Title: "Biostatistician III", Score: "A"},
		{UserId: "618-96-2322", Age: 69, Title: "Programmer I", Score: "C"},
		{UserId: "618-96-2322", Age: 69, Title: "Executive Secretary", Score: "C"},
	}
	if err := inserter.Put(ctx, items); err != nil {
		return err
	}

	return nil
}

type BigQueryTableItem struct {
	UserId string
	Age    int
	Title  string
	Score  string
}

func (i *BigQueryTableItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id": i.UserId,
		"age":     i.Age,
		"title":   i.Title,
		"score":   i.Score,
	}, bigquery.NoDedupeID, nil
}

func deleteBigQueryAssets(projectID, dataSetID string) error {
	log.Printf("[START] deleteBigQueryAssets: projectID %v", projectID)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Printf("[INFO] deleteBigQueryAssets: delete dataset err %v", err)

	if err := client.Dataset(dataSetID).DeleteWithContents(ctx); err != nil {
		log.Printf("[INFO] deleteBigQueryAssets: delete dataset err %v", err)
		return err
	}

	duration := time.Duration(30) * time.Second
	time.Sleep(duration)

	log.Printf("[END] deleteBigQueryAssets:")
	return nil
}
