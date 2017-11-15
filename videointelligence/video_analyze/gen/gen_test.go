// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gentest

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestFresh(t *testing.T) {
	testutil.Generated(t, "template.go").
		Matches("../video_analyze.go")

	testutil.Generated(t, "template.go").
		Labels("gcs").
		Matches("../video_analyze_gcs.go")
}
