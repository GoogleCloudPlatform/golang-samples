// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package testutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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

// Run runs the built binary with the given environment.
// After f returns, the running process is shut down.
func (r *Runner) Run(env map[string]string, f func()) {
	if !r.Built() {
		r.t.Error("Tried to run when binary not built.")
		return
	}
	environ := os.Environ()
	for k, v := range env {
		environ = append(environ, k+"="+v)
	}

	cmd := exec.Command(r.bin)
	cmd.Env = environ

	if err := cmd.Start(); err != nil {
		r.t.Error(err)
		return
	}

	// Run the user's tests.
	f()

	done := make(chan struct{})
	go func() {
		cmd.Wait()
		close(done)
	}()

	// Try to gracefully kill the process.
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		r.t.Error(err)
	}

	select {
	case <-time.After(5 * time.Second):
		r.t.Error("Timed out with SIGINT, trying SIGKILL.")
		if err := cmd.Process.Kill(); err != nil {
			r.t.Error(err)
		}
	case <-done:
	}
}
