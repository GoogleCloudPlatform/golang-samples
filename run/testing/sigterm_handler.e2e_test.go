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
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSigtermHandlerService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	runID := time.Now().Format("20060102-150405")
	service := cloudrunci.NewService("sigterm-handler", tc.ProjectID)
	service.Dir = "../sigterm-handler"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer GetLogEntries(service, runID, t)
	defer service.Clean()

	requestPath := "/"
	req, err := service.NewRequest("GET", requestPath)
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("client.Do: %s %s\n", req.Method, req.URL)

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}
}

func GetLogEntries(service *cloudrunci.Service, runID string, t *testing.T) {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create service and timestamp filters
	fiveMinAgo := time.Now().Add(-5 * time.Minute)
	timeFormat := fiveMinAgo.Format(time.RFC3339)
	filter := fmt.Sprintf(`resource.labels.service_name="%s" timestamp>="%s"`, fmt.Sprintf("%s-%s", service.Name, runID), timeFormat)
	preparedFilter := fmt.Sprintf(`resource.type="cloud_run_revision" severity="default" %s  NOT protoPayload.serviceName="run.googleapis.com"`, filter)

	fmt.Println("Waiting for logs...")
	time.Sleep(1 * time.Minute)
	MAX := 10
	for i := 0; i < MAX; i++ {
		fmt.Printf("Attempt #%d\n", i)
		it := client.Entries(ctx, logadmin.Filter(preparedFilter))
		for {
			entry, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Errorf("error fetching logs: %s", err)
			}
			if strings.Contains(fmt.Sprintf("%v", entry.Payload), "terminated signal caught") {
				fmt.Print("log entry: found.")
				return
			}
		}
		time.Sleep(15 * time.Second)
	}
	t.Errorf("log entry: not found.")
}
