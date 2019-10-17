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
	"bytes"
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	topicName        = "dlp-inspect-test-topic"
	subscriptionName = "dlp-inspect-test-sub"

	ssnFileName             = "fake_ssn.txt"
	nothingEventfulFileName = "nothing_eventful.txt"
	bucketName              = "golang-samples-dlp-test"
)

func TestInspectDatastore(t *testing.T) {
	t.Skip("https://github.com/GoogleCloudPlatform/golang-samples/issues/1039")

	tc := testutil.EndToEndTest(t)
	writeTestDatastoreFiles(t, tc.ProjectID)
	tests := []struct {
		kind string
		want string
	}{
		{
			kind: "SSNTask",
			want: "US_SOCIAL_SECURITY_NUMBER",
		},
		{
			kind: "BoringTask",
			want: "No results",
		},
	}
	for _, test := range tests {
		t.Run(test.kind, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			if err := inspectDatastore(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName, subscriptionName, tc.ProjectID, "", test.kind); err != nil {
				t.Errorf("inspectDatastore(%s) got err: %v", test.kind, err)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("inspectDatastore(%s) = %q, want %q substring", test.kind, got, test.want)
			}
		})
	}
}

type SSNTask struct {
	Description string
}

type BoringTask struct {
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

	kind = "BoringTask"
	name = "boringtask1"
	boringKey := datastore.NameKey(kind, name, nil)
	boringTask := BoringTask{
		Description: "Nothing meaningful",
	}
	if _, err := client.Put(ctx, boringKey, &boringTask); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
}

func TestInspectGCS(t *testing.T) {
	t.Skip("https://github.com/GoogleCloudPlatform/golang-samples/issues/1039")

	tc := testutil.SystemTest(t)
	writeTestGCSFiles(t, tc.ProjectID)
	tests := []struct {
		fileName string
		want     string
	}{
		{
			fileName: ssnFileName,
			want:     "US_SOCIAL_SECURITY_NUMBER",
		},
		{
			fileName: nothingEventfulFileName,
			want:     "No results",
		},
	}
	for _, test := range tests {
		t.Run(test.fileName, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			if err := inspectGCSFile(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName, subscriptionName, bucketName, test.fileName); err != nil {
				t.Errorf("inspectGCSFile(%s) got err: %v", test.fileName, err)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("inspectGCSFile(%s) = %q, want %q substring", test.fileName, got, test.want)
			}
		})
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
	if err := writeObject(ctx, bucket, nothingEventfulFileName, "Nothing eventful"); err != nil {
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

func TestInspectString(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	if err := inspectString(buf, tc.ProjectID, "I'm Gary and my email is gary@example.com"); err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectString got %q, want %q", got, want)
	}
}

func TestInspectTextFile(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	if err := inspectTextFile(buf, tc.ProjectID, "testdata/test.txt"); err != nil {
		t.Errorf("TestInspectTextFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectTextFile got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectTextFile got %q, want %q", got, want)
	}
}

type Item struct {
	Description string
}

const (
	harmlessTable = "harmless"
	harmfulTable  = "harmful"
	bqDatasetID   = "golang_samples_dlp"
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
	if err := uploadBigQuery(ctx, d, schema, harmlessTable, "Nothing meaningful"); err != nil {
		t.Fatalf("uploadBigQuery: %v", err)
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

func TestInspectBigquery(t *testing.T) {
	t.Skip("https://github.com/GoogleCloudPlatform/golang-samples/issues/1039")

	tc := testutil.EndToEndTest(t)

	mustCreateBigqueryTestFiles(t, tc.ProjectID, bqDatasetID)

	tests := []struct {
		table string
		want  string
	}{
		{
			table: harmfulTable,
			want:  "US_SOCIAL_SECURITY_NUMBER",
		},
		{
			table: harmlessTable,
			want:  "No results",
		},
	}
	for _, test := range tests {
		t.Run(test.table, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			if err := inspectBigquery(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName, subscriptionName, tc.ProjectID, bqDatasetID, test.table); err != nil {
				t.Errorf("inspectBigquery(%s) got err: %v", test.table, err)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("inspectBigquery(%s) = %q, want %q substring", test.table, got, test.want)
			}
		})
	}
}
