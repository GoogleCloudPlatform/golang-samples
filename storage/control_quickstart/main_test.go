// Copyright 2024 Google LLC
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
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	testPrefix      = "test-gcs-go-control"
	bucketExpiryAge = time.Hour * 24
)

func TestControlQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() {
		// Clean up any old buckets that may remain.
		if err := testutil.DeleteExpiredBuckets(client, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
			log.Printf("DeleteExpiredBuckets: %v", err)
		}
	})

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	stdOut, stdErr, err := m.Run(nil, time.Minute, "--bucket", bucketName)

	if err != nil {
		t.Errorf("stdout: %v", string(stdOut[:]))
		t.Errorf("stderr: %v", string(stdErr[:]))
		t.Errorf("execution failed: %v", err)
	}

	if got, want := string(stdOut[:]), "location type multi-region"; !strings.Contains(got, want) {
		t.Errorf("got output: %q, want to contain: %q", got, want)
	}
}
