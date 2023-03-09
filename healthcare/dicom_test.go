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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// The studyUID, seriesUID, and instanceUID are hard-coded into
// the metadata of the DICOM file and make up the DicomWebPath
// in the requests to retrieve studies/instances/frames/rendered.
const (
	studyPath          = "studies"
	studyUID           = "studies/1.3.6.1.4.1.11129.5.5.111396399361969898205364400549799252857604/"
	seriesUID          = "series/1.3.6.1.4.1.11129.5.5.195628213694300498946760767481291263511724/"
	instanceUID        = "instances/1.3.6.1.4.1.11129.5.5.153751009835107614666834563294684339746480/"
	studyOutputFile    = "study.multipart"
	instanceOutputFile = "instance.dcm"
	renderedOutputFile = "rendered_image.png"
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
		t.Fatalf("Unable to create test dataset: %v", err)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := createDICOMStore(buf, tc.ProjectID, location, datasetID, dicomStoreID); err != nil {
			r.Errorf("createDICOMStore got err: %v", err)
		}
	})

	if t.Failed() {
		return
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
			r.Errorf("getDICOMStores got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, dicomStoreName)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebStoreInstance(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, studyUID, "./testdata/dicom_00000001_000.dcm"); err != nil {
			r.Errorf("dicomStoreInstance got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		// Remove the output file if it already exists.
		os.Remove(studyOutputFile)
		buf.Reset()

		err := dicomWebRetrieveStudy(buf, tc.ProjectID, location, datasetID, dicomStoreID, studyUID, studyOutputFile)
		if err != nil {
			r.Errorf("dicomWebRetrieveStudy: %v", err)
		}

		got := buf.String()

		if want := "Study retrieved and downloaded to file: study.multipart\n"; !strings.Contains(got, want) {
			r.Errorf("got %q, want %q", got, want)
		}

		stat, err := os.Stat(studyOutputFile)
		if err != nil {
			r.Errorf("os.Stat: %v", err)
		}

		if stat.Size() == 0 {
			t.Error("Empty output DICOM study file")
		}

		os.Remove(studyOutputFile)
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		// Remove the output file if it already exists.
		os.Remove(instanceOutputFile)
		buf.Reset()

		err := dicomWebRetrieveInstance(buf, tc.ProjectID, location, datasetID, dicomStoreID, studyUID+seriesUID+instanceUID, instanceOutputFile)
		if err != nil {
			r.Errorf("dicomWebRetrieveInstance: %v", err)
		}

		got := buf.String()

		if want := "DICOM instance retrieved and downloaded to file: instance.dcm\n"; !strings.Contains(got, want) {
			r.Errorf("got %q, want %q", got, want)
		}

		stat, err := os.Stat(instanceOutputFile)
		if err != nil {
			r.Errorf("os.Stat: %v", err)
			return
		}

		if stat.Size() == 0 {
			r.Errorf("Empty output DICOM instance file")
		}

		os.Remove(instanceOutputFile)
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebSearchStudies(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, studyPath); err != nil {
			r.Errorf("dicomWebSearchStudies got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebSearchInstances(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID); err != nil {
			r.Errorf("dicomWebSearchInstances got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		// Remove the output file if it already exists.
		os.Remove(renderedOutputFile)
		buf.Reset()

		err := dicomWebRetrieveRendered(buf, tc.ProjectID, location, datasetID, dicomStoreID, studyUID+seriesUID+instanceUID+"rendered", renderedOutputFile)
		if err != nil {
			r.Errorf("dicomWebRetrieveRendered: %v", err)
		}

		got := buf.String()

		if want := "Rendered PNG image retrieved and downloaded to file: rendered_image.png\n"; !strings.Contains(got, want) {
			r.Errorf("got %q, want %q", got, want)
		}

		stat, err := os.Stat(renderedOutputFile)
		if err != nil {
			r.Errorf("os.Stat: %v", err)
		}

		if stat.Size() == 0 {
			t.Error("Empty output rendered image file")
		}

		os.Remove(renderedOutputFile)
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := dicomWebDeleteStudy(ioutil.Discard, tc.ProjectID, location, datasetID, dicomStoreID, studyUID); err != nil {
			r.Errorf("dicomWebDeleteStudy got err: %v", err)
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
