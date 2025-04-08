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
	tc := testutil.SystemTest(t)
	// Run the example to export to the samples project
	testStart := time.Now()
	sout, serr, err := m.Run(map[string]string{"GOOGLE_CLOUD_PROJECT": tc.ProjectID}, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to run the trace example binary: %v - \n%s\n%s\n", err, sout, serr)

	}
	testEnd := time.Now()

	// Check that the project contains our trace
	ctx := context.Background()
	client, err := trace.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create trace client: %v", err)
	}
	defer client.Close()
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		// Count the number of traces returned, and ensure it is non-zero
		req := &tracepb.ListTracesRequest{
			ProjectId: tc.ProjectID,
			StartTime: timestamppb.New(testStart),
			EndTime:   timestamppb.New(testEnd),
			Filter:    "root:foo",
		}
		var numMatchingTraces int
		it := client.ListTraces(ctx, req)
		for {
			_, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				r.Errorf("Failed to get next item from ListTraces: %v", err)
				r.Fail()
			}
			numMatchingTraces++
		}
		if numMatchingTraces == 0 {
			r.Errorf("No traces found matching the filter")
		}
	})
}
