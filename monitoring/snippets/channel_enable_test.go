// Copyright 2018 Google Inc. All rights reserved.
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

func TestEnableChannel(t *testing.T) {
	tc := testutil.SystemTest(t)

	c, err := createChannel(tc.ProjectID)
	if err != nil {
		t.Fatalf("Error creating test channel: %v", err)
	}
	defer deleteChannel(ioutil.Discard, c.GetName())

	buf := &bytes.Buffer{}
	if err := enableChannel(buf, c.GetName()); err != nil {
		t.Fatalf("enableChannel: %v", err)
	}
	want := "Enabled channel"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("enableChannel got %q, want to contain %q", got, want)
	}
}
