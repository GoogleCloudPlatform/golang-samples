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

package cloudruntests

import (
	"net/http"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPubSubService(t *testing.T) {
	// TODO: test failing due to:
	// pubsub.e2e_test.go:42: request: no acceptable response after 5 retries: %!w(<nil>)
	// See fusion log:
	// https://fusion2.corp.google.com/invocations/0286a27a-2a06-4c60-a4ed-84402ef2fd51/targets/cloud-devrel%2Fgo%2Fgolang-samples%2Fpresubmit%2Flatest-version;config=default/log
	t.Skip()

	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("pubsub", tc.ProjectID)
	service.Dir = "../pubsub"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer func(service *cloudrunci.Service) {
		err := service.Clean()
		if err != nil {
			t.Fatalf("service.Clean %q: %v", service.Name, err)
		}
	}(service)

	resp, err := service.Request("GET", "/",
		cloudrunci.WithAcceptFunc(func(resp *http.Response) bool {
			return resp.StatusCode != http.StatusBadRequest
		}),
	)

	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if got := resp.StatusCode; got != http.StatusBadRequest {
		t.Errorf("response status: got %d, want %d", got, http.StatusBadRequest)
	}
}
