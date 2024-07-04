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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	conditionaldelete "github.com/GoogleCloudPlatform/golang-samples/healthcare/internal/fhir-resource-conditional-delete"
	conditionalpatch "github.com/GoogleCloudPlatform/golang-samples/healthcare/internal/fhir-resource-conditional-patch"
	conditionalupdate "github.com/GoogleCloudPlatform/golang-samples/healthcare/internal/fhir-resource-conditional-update"
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
	resourceType := "Patient"

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

	testutil.Retry(t, 10, 60*time.Second, func(r *testutil.R) {
		if err := createFHIRStore(buf, tc.ProjectID, location, datasetID, fhirStoreID); err != nil {
			t.Errorf("createFHIRStore got err: %v", err)
		}
	})

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
		buf.Reset()
		if err := createFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType); err != nil {
			r.Errorf("createFHIRResource got err: %v", err)
		}
	})

	type resource struct {
		Active bool
		ID     string
		Meta   struct {
			VersionID string // Used for getFHIRResourceHistory.
		}
	}
	res := resource{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Errorf("json.Unmarshal createFHIRResource output: %v", err)
	}

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := getFHIRMetadata(buf, tc.ProjectID, location, datasetID, fhirStoreID); err != nil {
			r.Errorf("getFHIRMetadata got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, fhirStoreID) {
			r.Errorf("getFHIRMetadata got \n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, fhirStoreID)
		}
	})

	buf.Reset()
	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := getFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID); err != nil {
			r.Errorf("getFHIRResource got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, res.ID) {
			r.Errorf("getFHIRResource got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, res.ID)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := searchFHIRResourcesGet(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType); err != nil {
			r.Errorf("searchFHIRResourcesGet got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, res.ID) {
			r.Errorf("searchFHIRResourcesGet got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, resourceType)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := searchFHIRResourcesPost(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType); err != nil {
			r.Errorf("searchFHIRResourcesPost got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, res.ID) {
			r.Errorf("searchFHIRResourcesPost got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, resourceType)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := updateFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID, false); err != nil {
			r.Errorf("updateFHIRResource got err: %v", err)
			return
		}
		updatedRes := resource{}
		if err := json.Unmarshal(buf.Bytes(), &updatedRes); err != nil {
			r.Errorf("json.Unmarshal updateFHIRResource output: %v", err)
			return
		}
		if updatedRes.Active {
			r.Errorf("updateFHIRResource got active=true, expected active=false")
		}

		buf.Reset()
		updatedRes = resource{}
		if err := updateFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID, true); err != nil {
			r.Errorf("updateFHIRResource got err: %v", err)
			return
		}
		if err := json.Unmarshal(buf.Bytes(), &updatedRes); err != nil {
			r.Errorf("json.Unmarshal updateFHIRResource output: %v", err)
			return
		}
		if !updatedRes.Active {
			r.Errorf("updateFHIRResource got active=false, expected active=true")
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := patchFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID, false); err != nil {
			r.Errorf("patchFHIRResource got err: %v", err)
			return
		}
		patchedRes := resource{}
		if err := json.Unmarshal(buf.Bytes(), &patchedRes); err != nil {
			r.Errorf("json.Unmarshal patchFHIRResource output: %v", err)
			return
		}
		if patchedRes.Active {
			r.Errorf("patchFHIRResource got active=true, expected active=false")
		}

		buf.Reset()
		patchedRes = resource{}
		if err := updateFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID, true); err != nil {
			r.Errorf("patchFHIRResource got err: %v", err)
			return
		}
		if err := json.Unmarshal(buf.Bytes(), &patchedRes); err != nil {
			r.Errorf("json.Unmarshal patchFHIRResource output: %v", err)
			return
		}
		if !patchedRes.Active {
			r.Errorf("patchFHIRResource got active=false, expected active=true")
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		patchedRes := resource{}
		if err := conditionalpatch.ConditionalPatchFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, false); err != nil {
			r.Errorf("ConditionalPatchFHIRResource got err: %v", err)
			return
		}
		if err := json.Unmarshal(buf.Bytes(), &patchedRes); err != nil {
			r.Errorf("json.Unmarshal ConditionalPatchFHIRResource output: %v", err)
			return
		}
		if patchedRes.ID != res.ID {
			r.Errorf("ConditionalPatchFHIRResource got ID=%v, want %v", patchedRes.ID, res.ID)
			return
		}
		if patchedRes.Active {
			r.Errorf("ConditionalPatchFHIRResource got active=true, expected active=false")
		}

		buf.Reset()
		patchedRes = resource{}
		if err := conditionalupdate.ConditionalUpdateFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, true); err != nil {
			r.Errorf("ConditionalUpdateFHIRResource got err: %v", err)
			return
		}
		if err := json.Unmarshal(buf.Bytes(), &patchedRes); err != nil {
			r.Errorf("json.Unmarshal ConditionalUpdateFHIRResource output: %v", err)
			return
		}
		if patchedRes.ID != res.ID {
			r.Errorf("ConditionalUpdateFHIRResource got ID=%v, want %v; wrong condition led to creating a new resource?", patchedRes.ID, res.ID)
			return
		}
		if !patchedRes.Active {
			r.Errorf("ConditionalUpdateFHIRResource got active=false, expected active=true")
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := fhirGetPatientEverything(buf, tc.ProjectID, location, datasetID, fhirStoreID, res.ID); err != nil {
			r.Errorf("fhirGetPatientEverything got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, res.ID) {
			r.Errorf("fhirGetPatientEverything got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, res.ID)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := listFHIRResourceHistory(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID); err != nil {
			r.Errorf("listFHIRResourceHistory got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, res.Meta.VersionID) {
			r.Errorf("listFHIRResourceHistory got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, res.Meta.VersionID)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := getFHIRResourceHistory(buf, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID, res.Meta.VersionID); err != nil {
			r.Errorf("getFHIRResourceHistory got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, res.Meta.VersionID) {
			r.Errorf("getFHIRResourceHistory got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, res.Meta.VersionID)
		}
	})

	// Longer retry time to avoid bucket create/delete API quota issues.
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf.Reset()
		// Note: the Healthcare Agent needs access to the bucket.
		// golang-samples-tests have been given access, but other projects may
		// fail until they give the agent access to the bucket.
		bucketName := testutil.TestBucket(context.Background(), t, tc.ProjectID, "healthcare-test")
		gsURIPrefix := "gs://" + bucketName + "/fhir-export/"

		if err := exportFHIRResource(buf, tc.ProjectID, location, datasetID, fhirStoreID, gsURIPrefix); err != nil {
			r.Errorf("exportFHIRResource got err: %v", err)
		}

		if err := importFHIRResource(ioutil.Discard, tc.ProjectID, location, datasetID, fhirStoreID, gsURIPrefix+"**"); err != nil {
			r.Errorf("importFHIRResource got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := conditionaldelete.ConditionalDeleteFHIRResource(ioutil.Discard, tc.ProjectID, location, datasetID, fhirStoreID, resourceType); err != nil {
			r.Errorf("ConditionalDeleteFHIRResource got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		if err := deleteFHIRResource(ioutil.Discard, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID); err != nil {
			r.Errorf("deleteFHIRResource got err: %v", err)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := purgeFHIRResource(ioutil.Discard, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID); err != nil {
			r.Errorf("purgeFHIRResource got err: %v", err)
		}

		if err := getFHIRResource(ioutil.Discard, tc.ProjectID, location, datasetID, fhirStoreID, resourceType, res.ID); err == nil {
			r.Errorf("getFHIRResource got %q, want it to be not found after purgeFHIRResource", res.ID)
		}
	})

	testutil.Retry(t, 10, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := fhirExecuteBundle(buf, tc.ProjectID, location, datasetID, fhirStoreID); err != nil {
			r.Errorf("fhirExecuteBundle got err: %v", err)
		}
		if got, want := buf.String(), "201 Created"; !strings.Contains(got, want) {
			r.Errorf("fhirExecuteBundle got\n----\n%s\n----\nWant to contain:\n----\n%s\n----\n", got, want)
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
