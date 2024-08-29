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
	"io"
	"net/http"
	"testing"

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
	defer func(service *cloudrunci.Service) {
		err := service.Clean()
		if err != nil {
			t.Fatalf("service.Clean %q: %v", service.Name, err)
		}
	}(service)

	resp, err := service.Request("GET", "/")
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}

	if got, want := string(out), "Hello Override!\n"; got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}
}
