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

package functions_testing

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/functions/security"
	"github.com/GoogleCloudPlatform/golang-samples/internal/functest"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMakeGetRequest(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	fn := functest.NewCloudFunction("fixture", tc.ProjectID)
	fn.Entrypoint = "Fixture"
	fn.Dir = "../security"
	if err := fn.Deploy(); err != nil {
		t.Fatalf("CloudFunction.Deploy: %v", err)
	}
	defer fn.Teardown()

	// Test function-to-function requests.
	relayFn := functest.NewCloudFunction("relay", tc.ProjectID)
	relayFn.Entrypoint = "Fixture"
	relayFn.DeployCommand = fmt.Sprintf("gcloud functions deploy %s --trigger-http --update-env-vars TARGET_URL=%s", relayFn.DeployName(), fn.URL())
	relayFn.Dir = "../security"
	if err := relayFn.Deploy(); err != nil {
		t.Fatalf("CloudFunction.Deploy: %v", err)
	}
	defer relayFn.Teardown()

	// Prepare an authenticated request to the relay function.
	// The relay function will send it's own authenticated request to the fixture.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, relayFn.URL().String(), nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext: %v", err)
	}
	client, err := relayFn.HTTPClient()
	if err != nil {
		t.Fatalf("CloudFunction.HTTPClient: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Client.Do: %v", err)
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %v", err)
	}

	if want, got := "Success", string(out); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

}
