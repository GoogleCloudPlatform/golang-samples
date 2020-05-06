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
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetMetricDescriptor(t *testing.T) {
	tc := testutil.SystemTest(t)

	metricType := "custom.googleapis.com/golang-samples-tests/get"
	m, err := createCustomMetric(ioutil.Discard, tc.ProjectID, metricType)
	if err != nil {
		t.Fatalf("createMetric: %v", err)
	}
	defer deleteMetric(ioutil.Discard, m.GetName())

	testutil.Retry(t, 20, 10*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getMetricDescriptor(buf, tc.ProjectID, metricType); err != nil {
			r.Errorf("getMetricDescriptor: %v", err)
			return
		}
		want := "Name:"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("getMetricDescriptor got %q, want to contain %q", got, want)
		}
	})
}
