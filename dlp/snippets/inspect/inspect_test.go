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

package inspect

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

const (
	topicName        = "dlp-inspect-test-topic-"
	subscriptionName = "dlp-inspect-test-sub-"

	ssnFileName = "fake_ssn.txt"
	bucketName  = "golang-samples-dlp-test2"

	inspectsGCSTestFileName                 = "test.txt"
	filePathToUpload                        = "./testdata/test.txt"
	dirPathForInspectGCSSendToScc           = "dlp-go-lang-test-for-inspect-gcs-send-to-scc/"
	bucketnameForInspectGCSFileWithSampling = "dlp-job-go-lang-test-inspect-gcs-file-with-sampling"
)

type SSNTask struct {
	Description string
}

func writeTestDatastoreFiles(t *testing.T, projectID string) {
	t.Helper()
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("datastore.NewClient: %v", err)
	}
	kind := "SSNTask"
	name := "ssntask1"
	ssnKey := datastore.NameKey(kind, name, nil)
	task := SSNTask{
		Description: "My SSN is 111222333",
	}
	if _, err := client.Put(ctx, ssnKey, &task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
}

func writeTestGCSFiles(t *testing.T, projectID string) {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		switch err {
		case storage.ErrObjectNotExist:
			if err := bucket.Create(ctx, projectID, nil); err != nil {
				t.Fatalf("bucket.Create: %v", err)
			}
		default:
			t.Fatalf("error getting bucket attrs: %v", err)
		}
	}
	if err := writeObject(ctx, bucket, ssnFileName, "My SSN is 111222333"); err != nil {
		t.Fatalf("writeObject: %v", err)
	}
}

func writeObject(ctx context.Context, bucket *storage.BucketHandle, fileName, content string) error {
	obj := bucket.Object(fileName)
	_, err := obj.Attrs(ctx)
	if err != nil {
		switch err {
		case storage.ErrObjectNotExist:
			w := obj.NewWriter(ctx)
			w.Write([]byte(content))
			if err := w.Close(); err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil
}

type Item struct {
	Description string
}

const (
	harmfulTable = "harmful"
	bqDatasetID  = "golang_samples_dlp"
)

func mustCreateBigqueryTestFiles(t *testing.T, projectID, datasetID string) {
	t.Helper()

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()
	d := client.Dataset(datasetID)
	if _, err := d.Metadata(ctx); err != nil {
		if err := d.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}
	schema, err := bigquery.InferSchema(Item{})
	if err != nil {
		t.Fatalf("InferSchema: %v", err)
	}
	if err := uploadBigQuery(ctx, d, schema, harmfulTable, "My SSN is 111222333"); err != nil {
		t.Fatalf("uploadBigQuery: %v", err)
	}
}

func uploadBigQuery(ctx context.Context, d *bigquery.Dataset, schema bigquery.Schema, table, content string) error {
	t := d.Table(table)
	if _, err := t.Metadata(ctx); err == nil {
		return nil
	}
	if err := t.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		return err
	}
	source := bigquery.NewReaderSource(strings.NewReader(content))
	l := t.LoaderFrom(source)
	job, err := l.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	return status.Err()
}

// filePathtoGCS uploads a file test.txt in given path from the testdata directory.
func filePathtoGCS(t *testing.T, projectID, bucketNameForInspectGCSSendToScc, dirPathForInspectGCSSendToScc string) error {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Check if the bucket already exists.
	bucketExists := false
	_, err = client.Bucket(bucketNameForInspectGCSSendToScc).Attrs(ctx)
	if err == nil {
		bucketExists = true
	}

	// If the bucket doesn't exist, create it.
	if !bucketExists {
		if err := client.Bucket(bucketNameForInspectGCSSendToScc).Create(ctx, projectID, &storage.BucketAttrs{
			StorageClass: "STANDARD",
			Location:     "us-central1",
		}); err != nil {
			return err
		}
		fmt.Printf("Bucket '%s' created successfully.\n", bucketNameForInspectGCSSendToScc)
	} else {
		fmt.Printf("Bucket '%s' already exists.\n", bucketNameForInspectGCSSendToScc)
	}

	// Check if the directory already exists in the bucket.
	dirExists := false
	query := &storage.Query{Prefix: dirPathForInspectGCSSendToScc}
	it := client.Bucket(bucketNameForInspectGCSSendToScc).Objects(ctx, query)
	_, err = it.Next()
	if err == nil {
		dirExists = true
	}

	// If the directory doesn't exist, create it.
	if !dirExists {
		obj := client.Bucket(bucketNameForInspectGCSSendToScc).Object(dirPathForInspectGCSSendToScc)
		if _, err := obj.NewWriter(ctx).Write([]byte("")); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		fmt.Printf("Directory '%s' created successfully in bucket '%s'.\n", dirPathForInspectGCSSendToScc, bucketNameForInspectGCSSendToScc)
	} else {
		fmt.Printf("Directory '%s' already exists in bucket '%s'.\n", dirPathForInspectGCSSendToScc, bucketNameForInspectGCSSendToScc)
	}

	// file upload code

	// Open local file.
	file, err := ioutil.ReadFile(filePathToUpload)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return err
	}

	// Get a reference to the bucket
	bucket := client.Bucket(bucketNameForInspectGCSSendToScc)

	// Upload the file
	object := bucket.Object(inspectsGCSTestFileName)
	writer := object.NewWriter(ctx)
	_, err = writer.Write(file)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
		return err
	}
	fmt.Printf("File uploaded successfully: %v\n", inspectsGCSTestFileName)

	// Check if the file exists in the bucket
	_, err = bucket.Object(inspectsGCSTestFileName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			fmt.Printf("File %v does not exist in bucket %v\n", inspectsGCSTestFileName, bucketNameForInspectGCSSendToScc)
		} else {
			log.Fatalf("Failed to check file existence: %v", err)
		}
	} else {
		fmt.Printf("File %v exists in bucket %v\n", inspectsGCSTestFileName, bucketNameForInspectGCSSendToScc)
	}

	fmt.Println("filePathtoGCS function is executed-------")
	return nil
}
