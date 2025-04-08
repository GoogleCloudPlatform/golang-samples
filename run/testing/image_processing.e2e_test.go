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

func TestImageProcessingService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("imageproc", tc.ProjectID)
	service.Dir = "../image-processing"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

	requestPath := "/"
	resp, err := service.Request("POST", requestPath,
		cloudrunci.WithAcceptFunc(cloudrunci.AcceptNonServerError))
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}

	defer resp.Body.Close()

	if got := resp.StatusCode; got != http.StatusBadRequest {
		t.Errorf("response status: got %d, want %d", got, http.StatusBadRequest)
	}
}
