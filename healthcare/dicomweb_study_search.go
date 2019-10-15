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

// [START healthcare_dicomweb_search_studies]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// queryParamOpt is a googleapi.Option (https://godoc.org/google.golang.org/api/googleapi#CallOption)
// that adds query parameters to an API call.
type queryParamOpt struct {
	key, value string
}

func (qp queryParamOpt) Get() (string, string) { return qp.key, qp.value }

// dicomWebSearchStudies refines a DICOMweb studies search by appending DICOM tags to the request.
func dicomWebSearchStudies(w io.Writer, projectID, location, datasetID, dicomStoreID, dicomWebPath string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"
	// dicomStoreID := "my-dicom-store"
	// dicomWebPath := "studies"
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	call := storesService.SearchForStudies(name, dicomWebPath)
	// Refine your search by appending DICOM tags to the
	// request in the form of query parameters. This sample
	// searches for studies containing a patient's name.
	patientName := queryParamOpt{key: "PatientName", value: "Sally Zhang"}
	resp, err := call.Do(patientName)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("SearchForStudies: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	respString := string(respBytes)
	fmt.Fprintf(w, "Found studies: %s\n", respString)

	return nil
}

// [END healthcare_dicomweb_search_studies]
