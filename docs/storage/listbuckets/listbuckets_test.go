// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListBuckets(t *testing.T) {
	tc := testutil.SystemTest(t)

	buckets, err := ListBuckets(tc.ProjectID)
	if err != nil {
		t.Errorf("error while listing buckets: %s", err)
	}
	if len(buckets) <= 0 {
		t.Error("want non-empty list of buckets")
	}
}
