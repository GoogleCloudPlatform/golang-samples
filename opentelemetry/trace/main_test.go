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

package main

import (
	"context"
	"os"
	"testing"
	"time"

	trace "cloud.google.com/go/trace/apiv1"
	"cloud.google.com/go/trace/apiv1/tracepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestWriteTraces(t *testing.T) {
	// Tests build.
	m := testutil.BuildMain(t)
	if !m.Built() {
		t.Fatalf("failed to build app")
	}
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("Skipping tail logs sample test. Set GOLANG_SAMPLES_PROJECT_ID.")
	}
	// Run the example to export to the samples project
	testStart := time.Now()
	_, _, err := m.Run(map[string]string{"GOOGLE_CLOUD_PROJECT": projectID}, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to run the trace example binary: %v", err)
	}
	testEnd := time.Now()

	// Check that the project contains our trace
	ctx := context.Background()
	client, err := trace.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create trace client: %v", err)
	}
	defer client.Close()
	// Wait a few seconds to ensure our trace is returned by the trace API
	time.Sleep(5 * time.Second)

	req := &tracepb.ListTracesRequest{
		ProjectId: projectID,
		StartTime: timestamppb.New(testStart),
		EndTime:   timestamppb.New(testEnd),
		Filter:    "root:foo",
	}
	// Count the number of traces returned, and ensure it is non-zero
	var numMatchingTraces int
	it := client.ListTraces(ctx, req)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Failed to get next item from ListTraces: %v", err)
		}
		numMatchingTraces++
	}
	if numMatchingTraces == 0 {
		t.Errorf("No traces found matching the filter")
	}
}
