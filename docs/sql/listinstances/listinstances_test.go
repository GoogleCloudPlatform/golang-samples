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

	instances, err := ListInstances(tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	if len(instances) == 0 {
		t.Fatalf("expected non-zero SQL instances in project %q", tc.ProjectID)
	}
}
