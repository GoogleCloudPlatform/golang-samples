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

// Package job demonstrates interactions with BigQuery's job resources, which
// allow for execution of multiple kinds of workloads including queries, loads,
// copies, and data extracts.
package job

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

func TestJobs(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Control a job lifecycle explicitly: create, report status, cancel.
	exampleJobID, err := bqtestutil.UniqueBQName("golang_example_job")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}
	q := client.Query("Select 17 as foo")
	q.JobID = exampleJobID
	q.Priority = bigquery.BatchPriority
	q.Run(ctx)
	if err := getJobInfo(ioutil.Discard, tc.ProjectID, exampleJobID); err != nil {
		t.Errorf("getJobInfo(%s): %v", exampleJobID, err)
	}
	if err := cancelJob(tc.ProjectID, exampleJobID); err != nil {
		t.Errorf("cancelJobInfo(%s): %v", exampleJobID, err)
	}
	if err := listJobs(ioutil.Discard, tc.ProjectID); err != nil {
		t.Errorf("listJobs: %v", err)
	}
}

func TestCopiesAndExtracts(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	testDatasetID, err := bqtestutil.UniqueBQName("snippet_table_tests")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}
	if err := client.Dataset(testDatasetID).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID).DeleteWithContents(ctx)

	// Generate some dummy tables via a quick CTAS.
	if err := generateTableCTAS(client, testDatasetID, "table1"); err != nil {
		t.Fatalf("failed to generate example table1: %v", err)
	}
	if err := generateTableCTAS(client, testDatasetID, "table2"); err != nil {
		t.Fatalf("failed to generate example table2: %v", err)
	}

	if err := createJob(tc.ProjectID, "SELECT 17 as foo"); err != nil {
		t.Errorf("createJob: %v", err)
	}

	// Run copy job tests in parallel.
	t.Run("copy", func(t *testing.T) {
		t.Run("copyTable", func(t *testing.T) {
			t.Parallel()
			if err := copyTable(tc.ProjectID, testDatasetID, "table1", "copy1"); err != nil {
				t.Errorf("copyTable(%s): %v", testDatasetID, err)
			}
		})

		t.Run("copyTableWithCMEK", func(t *testing.T) {
			if bqtestutil.SkipCMEKTests() {
				t.Skip("Skipping CMEK tests")
			}
			t.Parallel()
			if err := copyTableWithCMEK(tc.ProjectID, testDatasetID, "copycmek"); err != nil {
				t.Errorf("copyTableWithCMEK(%s): %v", testDatasetID, err)
			}
		})

		t.Run("copyMultiTable", func(t *testing.T) {
			t.Parallel()
			if err := copyMultiTable(tc.ProjectID, testDatasetID, []string{"table1", "table2"}, testDatasetID, "copymulti"); err != nil {
				t.Errorf("copyMultiTable(%s): %v", testDatasetID, err)
			}
		})
	})

	// Extract tests - setup bucket
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer storageClient.Close()

	bucket, err := bqtestutil.UniqueBucketName("golang-example-bucket", tc.ProjectID)
	if err != nil {
		t.Fatalf("cannot generate unique bucket name: %v", err)
	}
	if err := storageClient.Bucket(bucket).Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatalf("cannot create bucket: %v", err)
	}

	model := client.DatasetInProject(tc.ProjectID, testDatasetID).Model("model")
	if err := generateModel(client, tc.ProjectID, testDatasetID, model.ModelID); err != nil {
		t.Fatalf("cannot create BQ ML model: %v", err)
	}
	defer model.Delete(ctx)

	// Run extract job tests in parallel.
	t.Run("extract", func(t *testing.T) {
		t.Run("exportTableAsCSV", func(t *testing.T) {
			t.Parallel()
			gcsURI := fmt.Sprintf("gs://%s/%s", bucket, "shakespeare.csv")
			if err := exportTableAsCSV(tc.ProjectID, gcsURI); err != nil {
				t.Errorf("exportTableAsCSV(%s): %v", gcsURI, err)
			}
		})

		t.Run("exportTableAsCompressedCSV", func(t *testing.T) {
			t.Parallel()
			gcsURI := fmt.Sprintf("gs://%s/%s", bucket, "shakespeare.csv.gz")
			if err := exportTableAsCompressedCSV(tc.ProjectID, gcsURI); err != nil {
				t.Errorf("exportTableAsCompressedCSV(%s): %v", gcsURI, err)
			}
		})
		t.Run("exportTableAsJSON", func(t *testing.T) {
			t.Parallel()
			gcsURI := fmt.Sprintf("gs://%s/%s", bucket, "shakespeare.json")
			if err := exportTableAsJSON(tc.ProjectID, gcsURI); err != nil {
				t.Errorf("exportTableAsJSON(%s): %v", gcsURI, err)
			}
		})
		t.Run("exportModel", func(t *testing.T) {
			t.Parallel()
			gcsURI := fmt.Sprintf("gs://%s/%s", bucket, "model")
			if err := exportModel(tc.ProjectID, testDatasetID, model.ModelID, gcsURI); err != nil {
				t.Errorf("exportModel(%s): %v", gcsURI, err)
			}
		})
	})

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

// generateModel creates an example BigQuery ML model.
func generateModel(client *bigquery.Client, projectID, datasetID, modelID string) error {
	ctx := context.Background()
	modelRef := fmt.Sprintf("%s.%s.%s", projectID, datasetID, modelID)

	// Create a ML model via a query.
	sql := fmt.Sprintf(`
	CREATE MODEL `+"`%s`"+`
	OPTIONS (
		model_type='linear_reg',
		max_iteration=1,
		learn_rate=0.4,
		learn_rate_strategy='constant'
	) AS (
		SELECT 'a' AS f1, 2.0 AS label
		UNION ALL
		SELECT 'b' AS f1, 3.8 AS label
	)`, modelRef)
	q := client.Query(sql)
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
