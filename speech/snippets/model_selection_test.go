// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestModelSelection(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	err := modelSelection(&buf, "../testdata/Google_Gnome.wav")
	if err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "the weather outside is sunny") {
		t.Fatalf(`modelSelection(../testdata/Google_Gnome.wav) = %q; want "the weather outside is sunny"`, got)
	}
}
