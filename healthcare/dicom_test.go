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

// TestDICOMStore runs all DICOM store tests to avoid having to
// create/delete DICOM stores for every sample function that needs to be
// tested.
func TestDICOMStore(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	location := "us-central1"
	datasetID := "dicom-dataset"
	dicomStoreID := "my-dicom-store"

	// Delete test dataset if it already exists.
	if err := getDataset(buf, tc.ProjectID, location, datasetID); err == nil {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
				r.Errorf("deleteDataset got err: %v", err)
			}
		})
	}

	if err := createDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
		t.Skipf("Unable to create test dataset: %v", err)
		return
	}

	if err := createDICOMStore(buf, tc.ProjectID, location, datasetID, dicomStoreID); err != nil {
		t.Errorf("createDICOMStore got err: %v", err)
	}

	dicomStoreName := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", tc.ProjectID, location, datasetID, dicomStoreID)
	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := listDICOMStores(buf, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("listDICOMStores got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, dicomStoreName) {
			r.Errorf("listDICOMStores got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, dicomStoreName)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := getDICOMStore(buf, tc.ProjectID, location, datasetID, dicomStoreID); err != nil {
			r.Errorf("getDICOMStore got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, dicomStoreName) {
			r.Errorf("listDICOMStores got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, dicomStoreName)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebStoreInstance(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, "studies/1.3.6.1.4.1.11129.5.5.111396399361969898205364400549799252857604", "./testdata/dicom_00000001_000.dcm"); err != nil {
			r.Errorf("dicomStoreInstance got err: %v", err)
		}
	})

	// testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
	// 	if err := dicomWebSearchInstances(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, "studies/1.3.6.1.4.1.11129.5.5.111396399361969898205364400549799252857604"); err != nil {
	// 		r.Errorf("dicomWebSearchInstances got err: %v", err)
	// 	}
	// })

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebRetrieveStudy(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, "studies/1.3.6.1.4.1.11129.5.5.111396399361969898205364400549799252857604"); err != nil {
			r.Errorf("dicomWebSearchInstances got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebDeleteStudy(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, "studies/1.3.6.1.4.1.11129.5.5.111396399361969898205364400549799252857604"); err != nil {
			r.Errorf("dicomWebSearchInstances got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteDICOMStore(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID); err != nil {
			r.Errorf("deleteDICOMStore got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("deleteDataset got err: %v", err)
		}
	})
}
