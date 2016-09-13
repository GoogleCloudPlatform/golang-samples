// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package listinstances

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListInstances(t *testing.T) {
	tc := testutil.SystemTest(t)

	// Just check the call succeeds.
	_, err := ListInstances(tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
}
