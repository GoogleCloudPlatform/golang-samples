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
	"bytes"
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSimpleApp(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	rows, err := query(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := printResults(&b, rows); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(b.String(), "views:") {
		t.Errorf("got output: %q; want it to contain views:", b.String())
	}
}
