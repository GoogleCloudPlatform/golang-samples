// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const metricType = "custom.googleapis.com/golang-samples-tests/get"

func TestGetMetricDescriptor(t *testing.T) {
	tc := testutil.SystemTest(t)

	m, err := createCustomMetric(ioutil.Discard, tc.ProjectID, metricType)
	if err != nil {
		t.Fatalf("createMetric: %v", err)
	}
	defer deleteMetric(ioutil.Discard, m.GetName())

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getMetricDescriptor(buf, tc.ProjectID, metricType); err != nil {
			r.Errorf("getMetricDescriptor: %v", err)
		}
		want := "Name:"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("getMetricDescriptor got %q, want to contain %q", got, want)
		}
	})
}
