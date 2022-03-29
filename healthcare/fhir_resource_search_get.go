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

package snippets

// [START healthcare_search_resources_get]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
)

// searchFhirResourcesGet uses a GET request to search for FHIR resources in a given FHIR store.
func searchFHIRResourcesGet(w io.Writer, projectID, location, datasetID, fhirStoreID, resourceType string) error {
	ctx := context.Background()

	// The Healthcare API endpoint, API version, and request path.
	name := fmt.Sprintf("https://healthcare.googleapis.com/v1/projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s", projectID, location, datasetID, fhirStoreID, resourceType)

	// DefaultClient returns an HTTP Client that uses the
	// DefaultTokenSource (Application Default Credentials)
	// to obtain authentication credentials.
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatal(err)
	}

	// NewRequest takes an io.Reader as its third argument,
	// but the GET request to search for FHIR resources does
	// not pass any data in its body.
	req, err := http.NewRequest(http.MethodGet, name, nil)
	if err != nil {
		return fmt.Errorf("NewRequest: %v", err)
	}

	// To set additional parameters for search filtering, append the
	// search terms as query parameters, then assign the encoded
	// query string to the request.
	// For example, to search for a Patient with the family name "Smith",
	// specify a Patient resourceType and then set the following:
	// q := req.URL.Query()
	// q.Add("family:exact", "Smith")
	// req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Do: %v", err)
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("search: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

// [END healthcare_search_resources_get]
