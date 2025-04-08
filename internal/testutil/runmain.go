// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"bytes"
	"context"
	"fmt"
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
	tmp, err := os.MkdirTemp("", "runmain-"+filepath.Base(wd)+"-")
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
// been reached, and indicates successful execution on return.  You can
// supply extra arguments for the binary via args.
func (r *Runner) Run(env map[string]string, timeout time.Duration, args ...string) (stdout, stderr []byte, err error) {
	if !r.Built() {
		return nil, nil, fmt.Errorf("tried to run when binary not built")
	}

	environ := os.Environ()
	for k, v := range env {
		environ = append(environ, k+"="+v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, r.bin, args...)
	cmd.Env = environ
	var bufOut, bufErr bytes.Buffer
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("could not execute binary: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return bufOut.Bytes(), bufErr.Bytes(), err
	}
	return bufOut.Bytes(), bufErr.Bytes(), nil
}
