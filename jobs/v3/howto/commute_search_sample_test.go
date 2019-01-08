// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"io/ioutil"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCommuteSearch(t *testing.T) {
	tc := testutil.SystemTest(t)
	if _, err := commuteSearch(ioutil.Discard, tc.ProjectID, testCompany.Name); err != nil {
		t.Fatalf("commuteSearch: %v", err)
	}
}
