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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestReadTimeSeriesAlign(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := readTimeSeriesAlign(buf, tc.ProjectID); err != nil {
		t.Errorf("readTimeSeriesAlign: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("readTimeSeriesAlign() = %q, want to contain %q", got, want)
	}
}

func TestReadTimeSeriesReduce(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := readTimeSeriesReduce(buf, tc.ProjectID); err != nil {
		t.Errorf("readTimeSeriesReduce: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("readTimeSeriesReduce() = %q, want to contain %q", got, want)
	}
}

func TestReadTimeSeriesFields(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := readTimeSeriesFields(buf, tc.ProjectID); err != nil {
		t.Errorf("readTimeSeriesFields: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("readTimeSeriesFields() = %q, want to contain %q", got, want)
	}
}
