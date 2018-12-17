// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListMonitoredResources(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if err := listMonitoredResources(buf, tc.ProjectID); err != nil {
		t.Fatalf("listMonitoredResources: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("listMonitoredResources got %q, want to contain %q", got, want)
	}
}
