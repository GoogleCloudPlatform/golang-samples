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

// Package testutil provides test helpers for the golang-samples repo.
package testutil

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"
)

var noProjectID = errors.New("GOLANG_SAMPLES_PROJECT_ID not set")

type Context struct {
	ProjectID string
	Dir       string
}

func (tc Context) Path(p ...string) string {
	p = append([]string{tc.Dir}, p...)
	return filepath.Join(p...)
}

// ContextMain gets a test context from a TestMain function.
// Useful for initializing global variables before running parallel system tests.
// ok is false if the project is not set up properly for system tests.
func ContextMain(m *testing.M) (tc Context, ok bool) {
	c, err := testContext()
	if err == noProjectID {
		return c, false
	} else if err != nil {
		log.Fatal(err)
	}
	return c, true
}

// SystemTest gets the test context.
// The test is skipped if the GOLANG_SAMPLES_PROJECT_ID environment variable is not set.
func SystemTest(t *testing.T) Context {
	tc, err := testContext()
	if err == noProjectID {
		t.Skip(err)
	} else if err != nil {
		t.Fatal(err)
	}

	return tc
}

// EndToEndTest gets the test context, and sets the test as Parallel.
// The test is skipped if the GOLANG_SAMPLES_E2E_TEST environment variable is not set.
func EndToEndTest(t *testing.T) Context {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip("GOLANG_SAMPLES_E2E_TEST not set")
	}

	tc, err := testContext()
	if err != nil {
		t.Fatal(err)
	}

	t.Parallel()

	return tc
}

func testContext() (Context, error) {
	tc := Context{}

	tc.ProjectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if tc.ProjectID == "" {
		return tc, noProjectID
	}

	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName}, "github.com/GoogleCloudPlatform/golang-samples")
	if err != nil {
		return tc, fmt.Errorf("Could not find golang-samples: %v", err)
	}
	tc.Dir = pkgs[0].PkgPath

	return tc, nil
}
