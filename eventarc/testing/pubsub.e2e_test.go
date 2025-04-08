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
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPubSubSinkService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("pubsub", tc.ProjectID)
	service.Dir = "../pubsub"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

	requestPath := "/"
	req, err := service.NewRequest("POST", requestPath)
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}

	client := http.Client{Timeout: 10 * time.Second}

	testutil.Retry(t, 3, 5, func(r *testutil.R) {
		fmt.Printf("Attempt #%d: client.Do: %s %s\n", r.Attempt, req.Method, req.URL)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusInternalServerError {
			r.Errorf("Cloud Run service not ready")
		} else if got := resp.StatusCode; got != http.StatusBadRequest {
			t.Errorf("response status: got %d, want %d", got, http.StatusBadRequest)
		}
	})
}
