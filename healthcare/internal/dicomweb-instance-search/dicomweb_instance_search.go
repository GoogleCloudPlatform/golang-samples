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

// Package dicomwebsearchinstances contains a sample for searching instances.
// It's in a separate package so each search sample can define queryParamOpt.
package dicomwebsearchinstances

// [START healthcare_dicomweb_search_instances]
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

// DicomWebSearchInstances searches instances.
func DicomWebSearchInstances(w io.Writer, projectID, location, datasetID, dicomStoreID, dicomWebPath string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"
	// dicomStoreID := "my-dicom-store"
	// dicomWebPath := "studies/1.3.6.1.4.1.11129.5.5.1113639985/series/1.3.6.1.4.1.11129.5.5.1953511724/instances/1.3.6.1.4.1.11129.5.5.9562821369"

	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	call := storesService.SearchForInstances(parent, dicomWebPath)
	// Refine your search by appending DICOM tags to the
	// request in the form of query parameters.
	includeAllFields := queryParamOpt{key: "includefield", value: "all"}
	resp, err := call.Do(includeAllFields)
	if err != nil {
		return fmt.Errorf("SearchForInstances: %v", err)
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("SearchForInstances: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	return nil
}

// [END healthcare_dicomweb_search_instances]
