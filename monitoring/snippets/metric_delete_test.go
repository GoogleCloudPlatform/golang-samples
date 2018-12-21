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

func TestDeleteMetric(t *testing.T) {
	tc := testutil.SystemTest(t)
	m, err := createCustomMetric(ioutil.Discard, tc.ProjectID, metricType)
	if err != nil {
		t.Fatalf("createCustomMetric: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := deleteMetric(buf, m.GetName()); err != nil {
		t.Fatalf("deleteMetric: %v", err)
	}
	want := "Deleted"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("deleteMetric got %q, want %q", got, want)
	}
}
