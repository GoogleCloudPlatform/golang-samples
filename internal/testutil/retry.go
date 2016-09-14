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

func Flaky(t *testing.T, max int, sleep time.Duration, f func(r *R)) bool {
	for attempt := 1; attempt <= max; attempt++ {
		r := &R{t: t, Attempt: attempt, log: &bytes.Buffer{}}

		last := attempt == max
		done := make(chan bool)
		go func() {
			defer close(done)
			f(r)
		}()
		<-done
		success := !r.retry && !r.fail

		if success || last || r.fail {
			state := "FAIL"
			if success {
				state = "SUCCESS"
				if r.log.Len() == 0 {
					return true
				}
			}
			t.Logf("Attempt %d: %s%s", attempt, state, r.log.String())
			if !success {
				t.Fail()
			}
		}
		if success {
			return true
		}
		if !r.retry {
			break
		}
		time.Sleep(sleep)
	}
	return false
}

type R struct {
	Attempt int

	t     *testing.T
	retry bool
	fail  bool
	log   *bytes.Buffer
}

func (r *R) Retry() {
	r.retry = true
}

func (r *R) RetryNow() {
	r.retry = true
	runtime.Goexit()
}

func (r *R) FailNow() {
	r.fail = true
	runtime.Goexit()
}

func (r *R) Fail() {
	r.fail = true
}

func (r *R) Retryf(s string, v ...interface{}) {
	r.logf(s, v...)
	r.Retry()
}

func (r *R) Fatalf(s string, v ...interface{}) {
	r.logf(s, v...)
	r.Fail()
	runtime.Goexit()
}

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
