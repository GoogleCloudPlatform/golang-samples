// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDownloadProfiles(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	projectID := tc.ProjectID
	ctx := context.Background()
	serviceVersion := fmt.Sprintf("v%v", time.Now().Unix())

	// We profile this test run to generate profiles which can later be fetched
	err := profiler.Start(profiler.Config{
		Service:              "aaa-test-download-profiles",
		ServiceVersion:       serviceVersion,
		NoHeapProfiling:      true,
		NoAllocProfiling:     true,
		NoGoroutineProfiling: true,
		DebugLogging:         true,
		ProjectID:            projectID,
	})
	if err != nil {
		t.Fatalf("failed to start the profiler: %v", err)
	}
	// do some work that generates profiles
	busyloop(t)
	var b bytes.Buffer
	// Download a single profile
	err = downloadProfiles(ctx, &b, projectID, "", 1, 1)
	if err != nil {
		t.Fatalf("download profiles failed: %v", err)
	}
	if output := b.String(); !strings.Contains(output, ProfilesDownloadedSuccessfully) || !strings.Contains(output, serviceVersion) {
		t.Errorf("downloadProfiles: expected output: %q to contain %q and %q", output, ProfilesDownloadedSuccessfully, serviceVersion)
	}
}

// Helper functions to keep CPU busy
// See profiler_quickstart for details on how this works.
func busyloop(t *testing.T) {
	t.Helper()
	for start := time.Now(); time.Since(start) < time.Minute*2; {
		load()
		runtime.Gosched()
	}
}

func load() {
	for i := 0; i < (1 << 20); i++ {
	}
}
