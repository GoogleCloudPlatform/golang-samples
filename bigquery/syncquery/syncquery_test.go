// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestQuery(t *testing.T) {
	tc := testutil.SystemTest(t)

	rows, err := Query(tc.ProjectID, "SELECT corpus FROM publicdata:samples.shakespeare GROUP BY corpus;")
	if err != nil {
		t.Fatal(err)
	}

	for _, row := range rows {
		if row[0] == "romeoandjuliet" {
			return
		}
	}
	t.Errorf("got rows: %q; want romeoandjuliet", rows)
}
