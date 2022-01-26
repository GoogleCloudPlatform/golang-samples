// Copyright 2022 Google LLC
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

package taillogs

import (
	"testing"
	"time"
	"context"
	"os"
	"fmt"

	"github.com/google/uuid"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
)

func TestTailLogs(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("Skipping tail logs sample test. Set GOLANG_SAMPLES_PROJECT_ID.")
	}

	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create logging client: %v", err)
	}
	defer client.Close()

	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()

	uuid, err := uuid.NewRandom()
	logID := fmt.Sprintf("tail-sample-log-%s", uuid.String()[:8])

	go func() {
		// 10 seconds is a recommended time to wait till streaming channel is established
		// when entries.tail reach PRD the sleep should be removed
		time.Sleep(10 * time.Second)

		logger := client.Logger(logID)
		defer logger.Flush()

		logger.Log(logging.Entry{
			Payload: "test tail logs entry 1",
			Severity: logging.Debug,
		})
		logger.Log(logging.Entry{
			Payload: "test tail logs entry 2",
			Severity: logging.Debug,
		})
	}()
	t.Cleanup(func() {
		adminClient.DeleteLog(ctx, logID)
	})

	// ingest a couple of logs to finish the test
	tailLogs(projectID)
}
