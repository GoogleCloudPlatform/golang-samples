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

// TestFHIRStore runs all FHIR store tests to avoid having to
// create/delete FHIR stores for every sample function that needs to be
// tested.
func TestFHIRStore(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	location := "us-central1"
	datasetID := "fhir-dataset"
	fhirStoreID := "my-fhir-store"
	if err := createDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
		t.Skipf("Unable to create test dataset: %v", err)
		return
	}

	if err := createFHIRStore(buf, tc.ProjectID, location, datasetID, fhirStoreID); err != nil {
		t.Errorf("createFHIRStore got err: %v", err)
	}

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		fhirStoreName := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", tc.ProjectID, location, datasetID, fhirStoreID)
		if err := listFHIRStores(buf, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("listFHIRStores got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, fhirStoreName) {
			r.Errorf("listFHIRStores got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, fhirStoreName)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteFHIRStore(ioutil.Discard, tc.ProjectID, location, datasetID, fhirStoreID); err != nil {
			r.Errorf("deleteFHIRStore got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("deleteDataset got err: %v", err)
		}
	})
}
