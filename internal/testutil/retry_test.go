// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package testutil

import (
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	Retry(t, 5, time.Millisecond, func(r *R) {
		if r.Attempt == 2 {
			return
		}
		r.Fail()
	})
}

func TestRetryAttempts(t *testing.T) {
	var attempts int
	Retry(t, 10, time.Millisecond, func(r *R) {
		r.Logf("This line should appear only once.")
		r.Logf("attempt=%d", r.Attempt)
		attempts = r.Attempt

		// Retry 5 times.
		if r.Attempt == 5 {
			return
		}
		r.Fail()
	})

	if attempts != 5 {
		t.Errorf("attempts=%d; want %d", attempts, 5)
	}
}
