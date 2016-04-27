// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package testutil provides test helpers for the golang-samples repo.
package testutil

import (
	"errors"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"testing"
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
	c, err := context()
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
	tc, err := context()
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

	tc, err := context()
	if err != nil {
		t.Fatal(err)
	}

	t.Parallel()

	return tc
}

func context() (Context, error) {
	tc := Context{}

	tc.ProjectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if tc.ProjectID == "" {
		return tc, noProjectID
	}

	pkg, err := build.Import("github.com/GoogleCloudPlatform/golang-samples", "", build.FindOnly)
	if err != nil {
		return tc, fmt.Errorf("Could not find golang-samples on GOPATH: %v", err)
	}
	tc.Dir = pkg.Dir

	return tc, nil
}
