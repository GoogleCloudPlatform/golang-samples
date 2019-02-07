// Copyright 2019 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"strings"
)

// MultiStringTest provides a way for testing many substrings in a response.
// By setting a set of expected substrings, and then indicating a nonzero
// MinPass or MinFail value the evaluation will yield an error.
//
// Intended use is for checking the responses from non-deterministic sources,
// such as ML-based APIs where exact string matching is flaky/problematic.
type MultiStringTest struct {
	Expected []string
	MinPass  int
	MinFail  int
}

// Evaluate returns an error if either the pass or failure rate is nonzero and unmet.
func (mst *MultiStringTest) Evaluate(got string) (passed, failed []string, err error) {
	// TODO(shollyman):  this is blind matching.  Should we include a mode for ordered
	// substrings?
	for _, v := range mst.Expected {
		if strings.Contains(got, v) {
			passed = append(passed, v)
		} else {
			failed = append(failed, v)
		}
	}

	if mst.MinPass != 0 && len(passed) < mst.MinPass {
		return passed, failed, fmt.Errorf("multistring min pass: %d of %d expected substrings", mst.MinPass, len(mst.Expected))
	}

	if mst.MinFail != 0 && len(failed) < mst.MinFail {
		return passed, failed, fmt.Errorf("multistring min fail: %d of %d expected substrings", mst.MinFail, len(mst.Expected))
	}
	return passed, failed, nil
}
