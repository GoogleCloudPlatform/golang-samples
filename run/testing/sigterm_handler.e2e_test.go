// Copyright 2021 Google LLC
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
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSigtermHandlerService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("sigterm-handler", tc.ProjectID)
	service.Dir = "../sigterm-handler"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer GetLogEntries(service, t)
	defer service.Clean()

	// Explicitly send SIGTERM
	req, err := service.NewRequest("GET", "")
	q := req.URL.Query()
	q.Add("terminate", "1")
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}

	resp, err := service.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	defer resp.Body.Close()
	log.Printf("client.Do: %s %s\n", req.Method, req.URL)

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}
}

func GetLogEntries(service *cloudrunci.Service, t *testing.T) {
	// Create timestamp filters
	minsAgo := time.Now().Add(-5 * time.Minute)
	timeFormat := minsAgo.Format(time.RFC3339)
	filter := fmt.Sprintf(`timestamp>="%s" severity="default" NOT protoPayload.serviceName="run.googleapis.com"`, timeFormat)

	find := "terminated signal caught"
	attempts := 6
	found, err := service.LogEntries(filter, find, attempts)
	if err != nil || !found {
		t.Errorf("%q log entry not found. (%v)", find, err)
	}
	return
}
