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

func TestAutoPunctuation(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	err := autoPunctuation(&buf, "../testdata/commercial_mono.wav")
	if err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "Okay. Sure.") {
		t.Fatalf(`autoPunctuation(../testdata/commercial_mono.wav) = %q; want "Okay. Sure"`, got)
	}
}
