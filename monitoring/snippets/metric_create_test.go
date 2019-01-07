// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateCustomMetric(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	m, err := createCustomMetric(buf, tc.ProjectID, metricType)
	if err != nil {
		t.Fatalf("createCustomMetric: %v", err)
	}
	defer deleteMetric(ioutil.Discard, m.GetName())

	want := "Created"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("createCustomMetric got %q, want %q", got, want)
	}
}
