// Copyright 2024 Google LLC
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

package amlai

// [START antimoneylaunderingai_list_locations]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
)

// listLocations lists all AML AI API locations for a given project.
func listLocations(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()

	// The AML AI API endpoint, API version, and request path.
	name := fmt.Sprintf("https://financialservices.googleapis.com/v1/projects/%s/locations", projectID)

	// DefaultClient returns an HTTP Client that uses the
	// Application Default Credentials to obtain authentication credentials.
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatal(err)
	}

	// NewRequest takes an io.Reader as its third argument,
	// but the GET request does not pass any data in its body.
	req, err := http.NewRequest(http.MethodGet, name, nil)
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}

	// Sets required header on the request.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Do: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

// [END antimoneylaunderingai_list_locations]
