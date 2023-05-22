// Copyright 2019 Google LLC
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

	"cloud.google.com/go/logging"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSimplelog(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := logging.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("logging.NewClient: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	defer func() {
		testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
			if err := deleteLog(tc.ProjectID); err != nil {
				r.Errorf("deleteLog: %v", err)
			}
		})
	}()

	client.OnError = func(err error) {
		t.Errorf("OnError: %v", err)
	}

	structuredWrite(tc.ProjectID)

	testutil.Retry(t, 20, 10*time.Second, func(r *testutil.R) {
		entries, err := getEntries(tc.ProjectID)
		if err != nil {
			r.Errorf("getEntries: %v", err)
			return
		}

		if len(entries) < 1 {
			r.Errorf("len(entries) = %d; want at least 1 entry", len(entries))
			return
		}

		wantContain := map[string]*logging.Entry{
			"Anything":                     entries[0],
			"The payload can be any type!": entries[0],
		}

		for want, entry := range wantContain {
			msg := fmt.Sprintf("%s", entry.Payload)
			if !strings.Contains(msg, want) {
				r.Errorf("want %q to contain %q", msg, want)
			}
		}
	})
}
