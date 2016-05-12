// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package newsletter

// [START utility_example_1]
import (
	"testing"

	"google.golang.org/appengine/aetest"
)

func TestWithContext(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	// Run code and tests requiring the context.Context using ctx.
	// [START_EXCLUDE]
	check(t, ctx)
	// [END_EXCLUDE]
}

// [END utility_example_1]

func check(t *testing.T, ctx interface{}) {
}
