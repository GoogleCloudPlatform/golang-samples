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

package cloudruntests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestHelloworldService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("helloworld", tc.ProjectID)
	service.Env = cloudrunci.EnvVars{"NAME": "Override"}
	service.Dir = "../helloworld"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

	req, err := service.NewRequest("GET", "/")
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}

	testutil.Retry(t, 10, 20*time.Second, func(r *testutil.R) {
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			r.Errorf("client.Do: %v", err)
		}
		fmt.Printf("client.Do: %s %s\n", req.Method, req.URL)

		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			r.Errorf("ioutil.ReadAll: %v", err)
		}

		if got, want := string(out), "Hello Override!\n"; got != want {
			r.Errorf("body: got %q, want %q", got, want)
		}

		if got := resp.StatusCode; got != http.StatusOK {
			r.Errorf("response status: got %d, want %d", got, http.StatusOK)
		}
	})
}
