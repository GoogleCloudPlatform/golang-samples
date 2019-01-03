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

func TestFilterOnStringValueCustomAttribute(t *testing.T) {
	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := filterOnStringValueCustomAttribute(buf, tc.ProjectID); err != nil {
			r.Errorf("filterOnStringValueCustomAttribute: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("filterOnStringValueCustomAttribute got %q, want to contain %q", got, want)
		}
	})
}

func TestFilterOnLongValueCustomAttribute(t *testing.T) {
	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := filterOnLongValueCustomAttribute(buf, tc.ProjectID); err != nil {
			r.Errorf("filterOnLongValueCustomAttribute: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("filterOnLongValueCustomAttribute got %q, want to contain %q", got, want)
		}
	})
}

func TestFilterOnMultiCustomAttributes(t *testing.T) {
	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if _, err := filterOnMultiCustomAttributes(buf, tc.ProjectID); err != nil {
			r.Errorf("filterOnMultiCustomAttributes: %v", err)
		}
		want := testJob.Name
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("filterOnMultiCustomAttributes got %q, want to contain %q", got, want)
		}
	})
}
