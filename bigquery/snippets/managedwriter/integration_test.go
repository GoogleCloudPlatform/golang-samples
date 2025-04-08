// Copyright 2022 Google LLC
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

package managedwriter

import (
	"context"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAppends(t *testing.T) {
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
	testDatasetID, err := bqtestutil.UniqueBQName("snippet_managedwriter_tests")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}
	if err := client.Dataset(testDatasetID).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID).DeleteWithContents(ctx)

	testTableID, err := bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}

	pendingSchema := bigquery.Schema{
		{Name: "bool_col", Type: bigquery.BooleanFieldType},
		{Name: "bytes_col", Type: bigquery.BytesFieldType},
		{Name: "float64_col", Type: bigquery.FloatFieldType},
		{Name: "int64_col", Type: bigquery.IntegerFieldType},
		{Name: "string_col", Type: bigquery.StringFieldType},
		{Name: "date_col", Type: bigquery.DateFieldType},
		{Name: "datetime_col", Type: bigquery.DateTimeFieldType},
		{Name: "geography_col", Type: bigquery.GeographyFieldType},
		{Name: "numeric_col", Type: bigquery.NumericFieldType},
		{Name: "bignumeric_col", Type: bigquery.BigNumericFieldType},
		{Name: "time_col", Type: bigquery.TimeFieldType},
		{Name: "timestamp_col", Type: bigquery.TimestampFieldType},

		{Name: "int64_list", Type: bigquery.IntegerFieldType, Repeated: true},
		{Name: "struct_col", Type: bigquery.RecordFieldType,
			Schema: bigquery.Schema{
				{Name: "sub_int_col", Type: bigquery.IntegerFieldType},
			}},
		{Name: "struct_list", Type: bigquery.RecordFieldType, Repeated: true,
			Schema: bigquery.Schema{
				{Name: "sub_int_col", Type: bigquery.IntegerFieldType},
			}},
		{Name: "row_num", Type: bigquery.IntegerFieldType, Required: true},
	}

	if err := client.Dataset(testDatasetID).Table(testTableID).Create(ctx, &bigquery.TableMetadata{
		Schema: pendingSchema,
	}); err != nil {
		t.Fatalf("failed to create destination table(%q %q): %v", testDatasetID, testTableID, err)
	}

	t.Run("PendingStream", func(t *testing.T) {
		if err := appendToPendingStream(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
			t.Errorf("appendToPendingStream(%q %q): %v", testDatasetID, testTableID, err)
		}
	})

	t.Run("DefaultStream", func(t *testing.T) {
		if err := appendToDefaultStream(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
			t.Errorf("appendToDefaultStream(%q %q): %v", testDatasetID, testTableID, err)
		}
	})

}
