// Copyright 2019 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestApp(t *testing.T) {
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	env := make(map[string]string)
	if p, ok := os.LookupEnv("GOLANG_SAMPLES_PROJECT_ID"); ok {
		env["GOOGLE_CLOUD_PROJECT"] = p
	}
	stdOut, stdErr, err := m.Run(env, 30*time.Second)
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}

	// don't look for specific strings, just expect at least 1kb of output
	if len(stdOut) < 1024 {
		t.Errorf("expected more output.  Stdout: %s", string(stdOut))
	}

	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes", len(stdErr))
	}
}
