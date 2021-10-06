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

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestApp(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	// Setup an output bucket.
	bucket, cleanup, err := setupStorage(tc.ProjectID)
	if err != nil {
		t.Fatalf("error setting up storage: %v", err)
	}
	defer cleanup()

	stdOut, stdErr, err := m.Run(nil, 30*time.Second, fmt.Sprintf("--project_id=%s", tc.ProjectID), fmt.Sprintf("--output=%s", bucket))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}

	// Look for a known substring in the output
	if !strings.Contains(string(stdOut), " ended in state COMPLETED") {
		t.Errorf("Did not find expected output.  Stdout: %s", string(stdOut))
	}

	if strings.Contains(string(stdOut), " with processing error") {
		t.Errorf("Workflow indicated it had processing errors.  Stdout: %s", string(stdOut))
	}

	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
}

// setupStorage is responsible for setting up a temporary bucket to hold artifacts from the quickstart.
func setupStorage(projectID string) (string, func(), error) {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", nil, err
	}
	bucket, err := bqtestutil.UniqueBucketName("golang-migration", "")
	if err != nil {
		storageClient.Close()
		return "", nil, fmt.Errorf("couldn't construct unique bucket name: %v", err)
	}
	if err := storageClient.Bucket(bucket).Create(ctx, projectID, nil); err != nil {
		storageClient.Close()
		return "", nil, fmt.Errorf("error creating output bucket: %v", err)
	}
	return bucket, func() {
		storageClient.Bucket(bucket).Delete(ctx)
		storageClient.Close()
	}, nil
}
