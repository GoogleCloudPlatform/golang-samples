// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

var inspectTopicName = "dlp-inspect-test-topic"
var inspectSubscriptionName = "dlp-inspect-test-sub"

func TestInspectString(t *testing.T) {
	testutil.SystemTest(t)
	tests := []struct {
		s    string
		want bool
	}{
		{
			s:    "My SSN is 111222333",
			want: true,
		},
		{
			s: "Does not match",
		},
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		inspectString(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, []string{"US_SOCIAL_SECURITY_NUMBER"}, nil, nil, test.s)
		if got := buf.String(); test.want != strings.Contains(got, "US_SOCIAL_SECURITY_NUMBER") {
			if test.want {
				t.Errorf("inspectString(%s) = %q, want 'US_SOCIAL_SECURITY_NUMBER' substring", test.s, got)
			} else {
				t.Errorf("inspectString(%s) = %q, want to not contain 'US_SOCIAL_SECURITY_NUMBER'", test.s, got)
			}
		}
		buf.Reset()
		inspectString(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, nil, []string{"SSN"}, []string{"\\d{9}"}, test.s)
		if got := buf.String(); test.want != strings.Contains(got, "CUSTOM_DICTIONARY_0") && strings.Contains(got, "CUSTOM_REGEX_0") {
			if test.want {
				t.Errorf("inspectString(%s) = %q, want 'CUSTOM_DICTIONARY_0' and 'CUSTOM_REGEX_0' substring", test.s, got)
			} else {
				t.Errorf("inspectString(%s) = %q, want to not contain 'CUSTOM_DICTIONARY_0' and 'CUSTOM_REGEX_0'", test.s, got)
			}
		}
	}
}

func TestInspectFile(t *testing.T) {
	testutil.SystemTest(t)
	tests := []struct {
		s    string
		want bool
	}{
		{
			s:    "My SSN is 111222333",
			want: true,
		},
		{
			s: "Does not match",
		},
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		inspectFile(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, []string{"US_SOCIAL_SECURITY_NUMBER"}, nil, nil, dlppb.ByteContentItem_TEXT_UTF8, strings.NewReader(test.s))
		if got := buf.String(); test.want != strings.Contains(got, "US_SOCIAL_SECURITY_NUMBER") {
			if test.want {
				t.Errorf("inspectString(%s) = %q, want 'US_SOCIAL_SECURITY_NUMBER' substring", test.s, got)
			} else {
				t.Errorf("inspectString(%s) = %q, want to not contain 'US_SOCIAL_SECURITY_NUMBER'", test.s, got)
			}
		}
		buf.Reset()
		inspectFile(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, nil, []string{"SSN"}, []string{"\\d{9}"}, dlppb.ByteContentItem_TEXT_UTF8, strings.NewReader(test.s))
		if got := buf.String(); test.want != strings.Contains(got, "CUSTOM_DICTIONARY_0") && strings.Contains(got, "CUSTOM_REGEX_0") {
			if test.want {
				t.Errorf("inspectString(%s) = %q, want 'CUSTOM_DICTIONARY_0' and 'CUSTOM_REGEX_0' substring", test.s, got)
			} else {
				t.Errorf("inspectString(%s) = %q, want to not contain 'CUSTOM_DICTIONARY_0' and 'CUSTOM_REGEX_0'", test.s, got)
			}
		}
	}
}

const (
	ssnFileName             = "fake_ssn.txt"
	nothingEventfulFileName = "nothing_eventful.txt"
	bucketName              = "golang-samples-dlp-test"
)

func writeTestGCSFiles(t *testing.T, projectID string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		switch err {
		case storage.ErrObjectNotExist:
			if err := bucket.Create(ctx, projectID, nil); err != nil {
				t.Fatalf("Failed to create bucket: %v", err)
			}
		default:
			t.Fatalf("error getting bucket attrs: %v", err)
		}
	}
	if err := writeObject(ctx, bucket, ssnFileName, "My SSN is 111222333"); err != nil {
		t.Fatalf("error writing object: %v", err)
	}
	if err := writeObject(ctx, bucket, nothingEventfulFileName, "Nothing eventful"); err != nil {
		t.Fatalf("error writing object: %v", err)
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

func TestInspectGCS(t *testing.T) {
	testutil.SystemTest(t)
	writeTestGCSFiles(t, projectID)
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
		buf := new(bytes.Buffer)
		inspectGCSFile(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, inspectTopicName, inspectSubscriptionName, bucketName, test.fileName)
		if got := buf.String(); !strings.Contains(got, test.want) {
			t.Errorf("inspectString(%s) = %q, want %q substring", test.fileName, got, test.want)
		}
	}
}

type SSNTask struct {
	Description string
}

type BoringTask struct {
	Description string
}

func writeTestDatastoreFiles(t *testing.T, projectID string) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
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

func TestInspectDatastore(t *testing.T) {
	testutil.SystemTest(t)
	writeTestDatastoreFiles(t, projectID)
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
		buf := new(bytes.Buffer)
		inspectDatastore(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, inspectTopicName, inspectSubscriptionName, projectID, "", test.kind)
		if got := buf.String(); !strings.Contains(got, test.want) {
			t.Errorf("inspectDatastore(%s) = %q, want %q substring", test.kind, got, test.want)
		}
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

func createBigqueryTestFiles(projectID, datasetID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()
	d := client.Dataset(datasetID)
	if _, err := d.Metadata(ctx); err != nil {
		if err := d.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
			return err
		}
	}
	schema, err := bigquery.InferSchema(Item{})
	if err != nil {
		return err
	}
	if err := uploadBigQuery(ctx, d, schema, harmlessTable, "Nothing meaningful"); err != nil {
		return err
	}
	return uploadBigQuery(ctx, d, schema, harmfulTable, "My SSN is 111222333")
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
	testutil.SystemTest(t)
	if err := createBigqueryTestFiles(projectID, bqDatasetID); err != nil {
		t.Fatalf("error creating test BigQuery files: %v", err)
	}
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
		buf := new(bytes.Buffer)
		inspectBigquery(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, true, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, inspectTopicName, inspectSubscriptionName, projectID, bqDatasetID, test.table)
		if got := buf.String(); !strings.Contains(got, test.want) {
			t.Errorf("inspectBigquery(%s) = %q, want %q substring", test.table, got, test.want)
		}
	}
}
