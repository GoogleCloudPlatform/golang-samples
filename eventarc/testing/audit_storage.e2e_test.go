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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudevent "github.com/cloudevents/sdk-go/v2"
)

func TestAuditStorageSinkService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("audit-storage", tc.ProjectID)
	service.Dir = "../audit_storage"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

	event := cloudevent.NewEvent("1.0")
	event.SetID("1")
	event.SetSource("test")
	event.SetSubject("storage.googleapis.com/projects/_/buckets/my-bucket")
	event.SetType("test")

	service_url, err := service.URL("/")
	if err != nil {
		t.Fatal(err)
	}
	req, err := cloudevent.NewHTTPRequestFromEvent(context.Background(),
		service_url, event)
	if err != nil {
		t.Fatal(err)
	}

	// add a valid auth header to the cloudevent request.
	authreq, _ := service.NewRequest("POST", "/")
	req.Header.Set("Authorization", authreq.Header.Get("Authorization"))

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("client.Do: %s %s\n", req.Method, req.URL)

	if got := resp.StatusCode; got != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(b)
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}
}
