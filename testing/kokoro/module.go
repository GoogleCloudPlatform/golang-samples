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

// Determine if the given module should be tested in the current environment.
//
// If versions cannot be evaluated, we fail "successful" to ensure tests are
// run by default.
//
// Usage:
// go run . -module=../../run/hello-broken
// This will exit with a non-zero status code if the module has a higher
// go version than the runtime.

// Command version checks for a module vs. environment go version mismatch
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
)

// maxVersionStr is the most advanced go version in testing.
const maxVersionStr = "1.13"

func main() {
	// Remove timestamp prefix from log output.
	log.SetFlags(0)

	modPathPtr := flag.String("module", ".", "Path to a go module")
	flag.Parse()

	v := strings.TrimLeft(runtime.Version(), "go")
	envVersion, err := version.NewVersion(v)
	if err != nil {
		log.Printf("version.NewVersion: %v", err)
	}
	maxVersion, err := version.NewVersion(maxVersionStr)
	if err != nil {
		log.Printf("version.NewVersion: %v", err)
	}

	if !validModule(maxVersion, envVersion, *modPathPtr) {
		os.Exit(1)
	}
}

// validModule determines if the module should be tested.
// Any of the following criteria mean it is valid for testing:
// - The current runtime is the most recent supported version of golang
// - The current runtime golang version is earlier than the module's required version
// - An error is encountered performing version comparisons
func validModule(max *version.Version, v *version.Version, module string) bool {
	if v.GreaterThan(max) {
		log.Printf("always run tests for most advanced go version: go%s", maxVersionStr)
		return true
	}

	mVersionStr, err := moduleVersion(module)
	if err != nil {
		log.Printf("moduleVersion: %v", err)
		return true
	}

	modVersion, err := version.NewVersion(mVersionStr)
	if err != nil {
		log.Printf("version.NewVersion: %v", err)
		return true
	}

	if v.LessThan(modVersion) {
		log.Printf("runtime version (%s) < module version (%s)", v, modVersion)
		return false
	}
	log.Printf("runtime version (%s) >= module version (%s)", v, modVersion)

	return true
}

// moduleVersion extracts the minimum go version from the go.mod
func moduleVersion(m string) (string, error) {
	p, err := filepath.Abs(m)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: (%s) %v", m, err)
	}

	cmd := exec.Command("go", []string{"mod", "edit", "-json"}...)
	cmd.Dir = p

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var module struct {
		Version string `json:"Go"`
	}
	if err := json.Unmarshal(out, &module); err != nil {
		return "", err
	}

	return module.Version, nil
}
