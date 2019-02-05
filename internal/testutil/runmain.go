// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package testutil

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// BuildMain builds the main package in the current working directory.
// If it doesn't build, t.Fatal is called.
// Test methods calling BuildMain should run Runner.Cleanup.
func BuildMain(t *testing.T) *Runner {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := ioutil.TempDir("", "runmain-"+filepath.Base(wd)+"-")
	if err != nil {
		t.Fatal(err)
	}

	r := &Runner{t: t, tmp: tmp}

	bin := filepath.Join(tmp, "a.out")
	cmd := exec.Command("go", "build", "-o", bin)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("go build: %v\n%s", err, out)
		return r
	}

	r.bin = bin
	return r
}

// Runner holds the result of `go build`
type Runner struct {
	t   *testing.T
	tmp string
	bin string
}

// Built reports whether the build was successful.
func (r *Runner) Built() bool {
	return r.bin != ""
}

// Cleanup removes the built binary.
func (r *Runner) Cleanup() {
	if err := os.RemoveAll(r.tmp); err != nil {
		r.t.Error(err)
	}
}

// Run executes runs the built binary until terminated or timeout has
// been reached, and indicates successful execution on return.
func (r *Runner) Run(env map[string]string, timeout time.Duration) (stdout, stderr []byte, err error) {
	if !r.Built() {
		return nil, nil, fmt.Errorf("tried to run when binary not built")
	}

	environ := os.Environ()
	for k, v := range env {
		environ = append(environ, k+"="+v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, r.bin)
	cmd.Env = environ
	var bufOut, bufErr bytes.Buffer
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("could not execute binary: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return bufOut.Bytes(), bufErr.Bytes(), err
	}
	return bufOut.Bytes(), bufErr.Bytes(), nil
}
