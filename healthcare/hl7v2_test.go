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

// TestHL7V2Store runs all HL7V2 store tests to avoid having to
// create/delete HL7V2 stores for every sample function that needs to be
// tested.
func TestHL7V2Store(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	location := "us-central1"
	datasetID := "hl7v2-dataset"
	hl7V2StoreID := "my-hl7v2-store"

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

	if err := createHL7V2Store(ioutil.Discard, tc.ProjectID, location, datasetID, hl7V2StoreID); err != nil {
		t.Fatalf("createHL7V2Store got err: %v", err)
	}

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		hl7V2StoreName := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7V2Stores/%s", tc.ProjectID, location, datasetID, hl7V2StoreID)
		if err := listHL7V2Stores(buf, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("listHL7V2Stores got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, hl7V2StoreName) {
			r.Errorf("listHL7V2Stores got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, hl7V2StoreName)
		}
	})

	messageID := "2yqbdhYHlk_ucSmWkcKOVm_N0p0OpBXgIlVG18rB-cw=" // TODO(cbro): use return value from create. seems to be stable though.

	dataFile := "testdata/hl7v2message.dat" // size = 167 bytes

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := createHL7V2Message(buf, tc.ProjectID, location, datasetID, hl7V2StoreID, dataFile); err != nil {
			r.Errorf("createHL7V2Message got err: %v", err)
		}
		if got, wantContain := buf.String(), messageID; !strings.Contains(got, wantContain) {
			r.Errorf("createHL7V2Message got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, wantContain)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := ingestHL7V2Message(buf, tc.ProjectID, location, datasetID, hl7V2StoreID, dataFile); err != nil {
			r.Errorf("ingestHL7V2Message got err: %v", err)
		}
		if got, wantContain := buf.String(), messageID; !strings.Contains(got, wantContain) {
			r.Errorf("ingestHL7V2Message got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, wantContain)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := getHL7V2Message(buf, tc.ProjectID, location, datasetID, hl7V2StoreID, messageID); err != nil {
			r.Errorf("getHL7V2Message got err: %v", err)
		}
		if got, wantContain := buf.String(), "Raw length: 167"; !strings.Contains(got, wantContain) {
			r.Errorf("getHL7V2Message got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, wantContain)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := patchHL7V2Message(ioutil.Discard, tc.ProjectID, location, datasetID, hl7V2StoreID, messageID, dataFile); err != nil {
			r.Errorf("patchHL7V2Message got err: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := deleteHL7V2Message(ioutil.Discard, tc.ProjectID, location, datasetID, hl7V2StoreID, messageID); err != nil {
			r.Errorf("deleteHL7V2Message got err: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := deleteHL7V2Store(ioutil.Discard, tc.ProjectID, location, datasetID, hl7V2StoreID); err != nil {
			r.Errorf("deleteHL7V2Store got err: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := deleteDataset(ioutil.Discard, tc.ProjectID, location, datasetID); err != nil {
			r.Errorf("deleteDataset got err: %v", err)
		}
	})
}
