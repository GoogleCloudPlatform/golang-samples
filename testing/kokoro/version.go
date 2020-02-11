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

// maxVersion is the most advanced go version in testing.
const maxVersionStr="1.13"
var maxVersion *version.Version

func main() {
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

	ok, err := valid(maxVersion, envVersion, *modPathPtr)
	if err != nil {
		log.Printf("compare: %v", err)
	}
	if !ok {
		os.Exit(1)
	}
}

// valid determines if the module should be tested.
func valid(max *version.Version, v *version.Version, module string) (bool, error) {
	if v.GreaterThan(max) {
		log.Printf("always run tests for most advanced go version: go%s", maxVersionStr)
		return true, nil
	}

	return compare(v, module)
}

// compare determines if the module has a newer version than the runtime.
// This command should fail "successful".
func compare(v *version.Version, module string) (bool, error) {
	mVersionStr, err := moduleVersion(module)
	if err != nil {
		return true, fmt.Errorf("version.NewVersion: %v", err)
	}

	modVersion, err := version.NewVersion(mVersionStr)
	if err != nil {
		return true, fmt.Errorf("version.NewVersion: %v", err)
	}
	
	if v.LessThan(modVersion) {
		log.Printf("runtime version (%s) < module version (%s)", v, modVersion)
		return false, nil
	}
	log.Printf("runtime version (%s) >= module version (%s)", v, modVersion)

	return true, nil
}

// moduleVersion extracts the minimum go version from the go.mod
func moduleVersion(m string) (string, error) {
	p, err := filepath.Abs(m)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: (%s) %v", m, err)
	}
	log.Printf("parsing module %q", p)

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