// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSearchFeaturedJobs(t *testing.T) {
	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := searchFeaturedJobs(buf, tc.ProjectID, testCompany.Name, "SWE"); err != nil {
			r.Errorf("searchFeaturedJobs: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("searchFeaturedJobs got %q, want to contain %q", got, want)
		}
	})
}
