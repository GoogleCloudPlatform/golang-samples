// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"
)

func TestListBuckets(t *testing.T) {
	t.Skip("TODO(cbro): make this test run in golang-samples")

	buckets, err := ListBuckets(os.Getenv("TEST_PROJECT_ID"))
	if err != nil {
		t.Errorf("Error while listing buckets: %s", err)
	}
	if len(buckets) <= 0 {
		t.Error("No bucket returned")
	}

	foundBucket := false
	expectedBucket := os.Getenv("TEST_BUCKET_NAME")
	for _, bucket := range buckets {
		if bucket.Name == expectedBucket {
			foundBucket = true
			break
		}
	}
	if !foundBucket {
		t.Errorf("Expected bucket %s", expectedBucket)
	}

}
