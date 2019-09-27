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

package snippets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestDataset runs all dataset tests to avoid having to create/delete
// datasets for every sample function that needs to be tested.
func TestDataset(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	location := "us-central1"
	datasetID := "my-dataset"
	deidentifiedDatasetID := "my-dataset-deidentified"

	// Delete test datasets if they already exist.
	if err := getDataset(buf, tc.ProjectID, location, datasetID); err == nil {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
				r.Errorf("deleteDataset got err: %v", err)
			}
		})
	}
	if err := getDataset(buf, tc.ProjectID, location, deidentifiedDatasetID); err == nil {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, deidentifiedDatasetID); err != nil {
				r.Errorf("deleteDataset got err: %v", err)
			}
		})
	}

	if err := createDataset(buf, tc.ProjectID, location, datasetID); err != nil {
		t.Fatalf("createDataset got err: %v", err)
	}
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", tc.ProjectID, location, datasetID)
	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := listDatasets(buf, tc.ProjectID, location); err != nil {
			r.Errorf("listDatasets got err: %v", err)
			return
		}
		if got := buf.String(); !strings.Contains(got, name) {
			r.Errorf("listDatasets got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, name)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := getDataset(buf, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("getDataset got err: %v", err)
			return
		}
		if got := buf.String(); !strings.Contains(got, name) {
			r.Errorf("getDataset got %q; want to contain %q", got, name)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deidentifyDataset(ioutil.Discard, tc.ProjectID, location, datasetID, deidentifiedDatasetID); err != nil {
			r.Errorf("deidentifyDataset got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := patchDataset(ioutil.Discard, tc.ProjectID, location, datasetID, "UTC"); err != nil {
			r.Errorf("patchDataset got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("deleteDataset got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, deidentifiedDatasetID); err != nil {
			r.Errorf("deleteDataset (deidentified) got err: %v", err)
		}
	})
}
