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
	t      *testing.T
	tmp    string
	bin    string
	stdout bytes.Buffer
	stderr bytes.Buffer
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

// Stdout returns the stdout from the most recent Run()
func (r *Runner) Stdout() []byte {
	return r.stdout.Bytes()
}

// Stderr returns the stderr from the most recent Run()
func (r *Runner) Stderr() []byte {
	return r.stderr.Bytes()
}

// Run executes runs the built binary until terminated or timeout has
// been reached, and indicates successful execution on return.
func (r *Runner) Run(env map[string]string, timeout time.Duration) error {
	if !r.Built() {
		return fmt.Errorf("tried to run when binary not built")
	}
	r.stderr.Reset()
	r.stdout.Reset()
	environ := os.Environ()
	for k, v := range env {
		environ = append(environ, k+"="+v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, r.bin)
	cmd.Env = environ

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not get stderr pipe")
	}
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not get stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not execute binary: %v", err)
	}

	b, _ := ioutil.ReadAll(stdOut)
	r.stdout.Write(b)
	b, _ = ioutil.ReadAll(stdErr)
	r.stderr.Write(b)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for termination: %v", err)
	}
	return nil
}
