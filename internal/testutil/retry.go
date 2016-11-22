// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package testutil

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// Retry runs function f for up to maxAttempts times until f returns successfully, and reports whether f was run successfully.
// It will sleep for the given period between invocations of f.
// Use the provided *testutil.R instead of a *testing.T from the function.
func Retry(t *testing.T, maxAttempts int, sleep time.Duration, f func(r *R)) bool {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		r := &R{Attempt: attempt, log: &bytes.Buffer{}}

		f(r)

		if !r.failed {
			if r.log.Len() != 0 {
				t.Logf("Success after %d attempts:%s", attempt, r.log.String())
			}
			return true
		}

		if attempt == maxAttempts {
			t.Logf("FAILED after %d attempts:%s", attempt, r.log.String())
			t.Fail()
		}

		time.Sleep(sleep)
	}
	return false
}

// R is passed to each run of a flaky test run, manages state and accumulates log statements.
type R struct {
	// The number of current attempt.
	Attempt int

	failed bool
	log    *bytes.Buffer
}

// Fail marks the run as failed, and will retry once the function returns.
func (r *R) Fail() {
	r.failed = true
}

// Errorf is equivalent to Logf followed by Fail.
func (r *R) Errorf(s string, v ...interface{}) {
	r.logf(s, v...)
	r.Fail()
}

// Logf formats its arguments and records it in the error log.
// The text is only printed for the final unsuccessful run or the first successful run.
func (r *R) Logf(s string, v ...interface{}) {
	r.logf(s, v...)
}

func (r *R) logf(s string, v ...interface{}) {
	fmt.Fprint(r.log, "\n")
	fmt.Fprint(r.log, lineNumber())
	fmt.Fprintf(r.log, s, v...)
}

func lineNumber() string {
	_, file, line, ok := runtime.Caller(3) // logf, public func, user function
	if !ok {
		return ""
	}
	return filepath.Base(file) + ":" + strconv.Itoa(line) + ": "
}
