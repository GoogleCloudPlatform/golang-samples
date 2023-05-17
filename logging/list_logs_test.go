// Copyright 2023 Google LLC
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

package snippets

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const logID = "test-list-logs"

var projectID string

func TestMain(m *testing.M) {
	ctx := context.Background()

	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("test project not set up properly")
		return
	}
	projectID = tc.ProjectID

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("logging.NewClient(%q) failed: %v", projectID, err)
	}
	defer client.Close()

	// Create a log
	logger := client.Logger(logID)
	logger.Log(logging.Entry{Payload: "create a log"})
	if err := logger.Flush(); err != nil {
		log.Fatalf("logger.Flush() failed: %v", err)
	}

	m.Run()
}

func TestListLogs(t *testing.T) {
	testutil.Retry(t, 6, 10*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := listLogs(buf, projectID); err != nil {
			r.Errorf("listLogs(%q) failed: %v", projectID, err)
			return
		}
		if !strings.Contains(buf.String(), logID) {
			r.Errorf("listLogs got %q, want to contain %q", buf.String(), logID)
		}
	})
}
