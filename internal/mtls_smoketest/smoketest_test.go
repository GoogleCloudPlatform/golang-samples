// Copyright 2020 Google LLC
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

package mtls_smoketest

import (
	"context"
	"os"
	"testing"
	"time"

	bqstorage "cloud.google.com/go/bigquery/storage/apiv1"
	gaming "cloud.google.com/go/gaming/apiv1beta"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	bqstoragepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1beta"
)

var shouldFail = os.Getenv("GOOGLE_API_USE_MTLS") == "always"

// checkErr expects an error under mtls_smoketest, and no error otherwise.
func checkErr(err error, t *testing.T) {
	t.Helper()
	if shouldFail && err == nil {
		t.Fatalf("got no err when wanted one - this means you should delete this test and un-skip the tests it's referring to.")
	}
	if !shouldFail && err != nil {
		t.Fatalf("got err when wanted no error: %v", err)
	}
}

// When this test starts failing, delete it and the corresponding lines in system_tests.bash
//
// vision/detect
// vision/label
// vision/product_search
// run/image-processing/imagemagick
func TestVision(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	ctx := context.Background()
	// NOTE(cbro): Observed successful and unsuccessful calls take under 1s.
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithQuotaProject(tc.ProjectID))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	f, err := os.Open("../../vision/testdata/cat.jpg")
	if err != nil {
		t.Fatal(err)
	}
	image, err := vision.NewImageFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.DetectLabels(ctx, image, nil, 10)
	checkErr(err, t)
}

// When this test starts failing, delete it and the corresponding lines in system_tests.bash
//
// bigquery/bigquery_storage_quickstart
func TestBigquerystorage(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	ctx := context.Background()
	// NOTE(cbro): Observed successful calls take around 1s. Unsuccessful calls hang indefinitely.
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	client, err := bqstorage.NewBigQueryReadClient(ctx)
	if err != nil {
		t.Fatalf("NewBigQueryStorageClient: %v", err)
	}
	defer client.Close()

	createReadSessionRequest := &bqstoragepb.CreateReadSessionRequest{
		Parent: "projects/" + tc.ProjectID,
		ReadSession: &bqstoragepb.ReadSession{
			Table:      "projects/bigquery-public-data/datasets/usa_names/tables/usa_1910_current",
			DataFormat: bqstoragepb.DataFormat_AVRO,
		},
		MaxStreamCount: 1,
	}

	_, err = client.CreateReadSession(ctx, createReadSessionRequest)
	checkErr(err, t)
}

// When this test starts failing, delete it and the corresponding lines in system_tests.bash
//
// gaming/servers
func TestGameservices(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	ctx := context.Background()
	// NOTE(cbro): Observed successful and unsuccessful calls take under 1s.
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	client, err := gaming.NewRealmsClient(ctx)
	if err != nil {
		t.Fatalf("NewRealmsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.ListRealmsRequest{
		Parent: "projects/" + tc.ProjectID + "/locations/global",
	}

	it := client.ListRealms(ctx, req)
	_, err = it.Next()
	if err == iterator.Done {
		err = nil
	}
	checkErr(err, t)
}
