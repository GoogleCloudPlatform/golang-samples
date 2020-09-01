// Copyright 2020 Google LLC
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

/*
Package functest facilitates end-to-end testing of Cloud Functions.

Example Usage

	package main_test

	import (
		"io/ioutil"
		"log"
		"os"
		"strings"
		"testing"

		"github.com/GoogleCloudPlatform/golang-samples/internal/functest"
	)

	func TestCloudFunction(t *testing.T) {
		// Define the Cloud Function to send requests.
		fn := functest.NewCloudFunction("hello", os.Getenv("GOOGLE_CLOUD_PROJECT"))
		// Override the Entrypoint if different from the function name.
		fn.Entrypoint = "HelloWorld"

		// This is not required: you could test an existing function.
		if err := fn.Deploy(); err != nil {
			log.Fatalf("CloudFunction.Deploy: %v", err)
		}
		defer fn.Teardown()

		// Create an HTTP Request.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fn.URL(), nil)
		if err != nil {
			log.Fatalf("http.NewRequestWithContext: %v", err)
		}

		// Get the HTTP client and execute the request.
		// Automatically adds an identity token and light request/response logging.
		client, err := fn.HTTPClient()
		if err != nil {
			log.Fatalf("CloudFunction.HTTPClient: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Client.Do: %v", err)
		}
		defer resp.Body.Close()

		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("ioutil.ReadAll: %q", err)
		}

		want := "Hello World!"
		if got := string(got); !strings.Contains(got, want) {
			t.Errorf("got\n----\n%s\n----\nWant to contain:\n----\n%s\n", got, shouldContain)
		}
	}

*/
package functest
