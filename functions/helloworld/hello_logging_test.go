// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package helloworld

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestHelloLogging(t *testing.T) {
	rStd, wStd, _ := os.Pipe()
	stdLogger.SetOutput(wStd)

	rErr, wErr, _ := os.Pipe()
	logger.SetOutput(wErr)

	HelloLogging(nil, nil)

	wStd.Close()
	wErr.Close()

	stdout, err := ioutil.ReadAll(rStd)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := "log entry"
	if got := string(stdout); !strings.Contains(got, want) {
		t.Errorf("Stdout got %q, want to contain %q", got, want)
	}

	stderr, err := ioutil.ReadAll(rErr)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want = "error"
	if got := string(stderr); !strings.Contains(got, want) {
		t.Errorf("Stderr got %q, want to contain %q", got, want)
	}
}
