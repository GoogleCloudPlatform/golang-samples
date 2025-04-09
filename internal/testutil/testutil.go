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
	"strings"
	"testing"
)

var errNoProjectID = errors.New("GOLANG_SAMPLES_PROJECT_ID not set")

// Context holds information useful for tests.
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
	if err == errNoProjectID {
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
	if err == errNoProjectID {
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

// Helper for getting the path to the media directory
func getMedia() string {
	// runtime.Caller returns information about the caller.
	// 0 identifies the getMedia function itself.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to determine caller information")
	}
	// file is the full path to this source file.
	dir := filepath.Dir(file)
	// Adjust the relative path as needed.
	return filepath.Join(dir, ".", "images")
}

func testContext() (Context, error) {
	tc := Context{}

	tc.ProjectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if tc.ProjectID == "" {
		return tc, errNoProjectID
	}

	dir, err := os.Getwd()
	if err != nil {
		return tc, fmt.Errorf("could not find current directory")
	}
	if !strings.Contains(dir, "golang-samples") {
		return tc, fmt.Errorf("could not find golang-samples directory")
	}
	tc.Dir = dir[:strings.Index(dir, "golang-samples")+len("golang-samples")]

	if tc.Dir == "" {
		return tc, fmt.Errorf("could not find golang-samples directory")
	}

	return tc, nil
}
