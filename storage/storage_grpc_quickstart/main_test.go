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
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	testPrefix      = "golang-grpc-test"
	bucketExpiryAge = time.Hour * 24
)

func TestGRPCQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}
	os.Setenv("GOOGLE_CLOUD_DISABLE_DIRECT_PATH", "true")
	ctx := context.Background()
	client, err := storage.NewGRPCClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewGRPCClient: %v", err)
	}
	t.Cleanup(func() {
		// Clean up any old buckets that may remain.
		if err := testutil.DeleteExpiredBuckets(client, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
			log.Printf("DeleteExpiredBuckets: %v", err)
		}
	})

	bucketName := testutil.UniqueBucketName(testPrefix)

	stdOut, stdErr, err := m.Run(nil, time.Minute, "--project", tc.ProjectID, "--bucket", bucketName)

	if err != nil {
		t.Errorf("stdout: %v", string(stdOut[:]))
		t.Errorf("stderr: %v", string(stdErr[:]))
		t.Errorf("execution failed: %v", err)
	}

	strStdErr := string(stdErr[:])
	if got, want := strStdErr, "Failed to enable client metrics"; strings.Contains(got, want) {
		t.Errorf("got output: %q, want to contain: %q", got, want)
	}

	if got, want := string(stdOut[:]), fmt.Sprintf("Bucket %v", bucketName); !strings.Contains(got, want) {
		t.Errorf("got output: %q, want to contain: %q", got, want)
	}
}
