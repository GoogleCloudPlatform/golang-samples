// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package testutil provides test helpers for the golang-samples repo.
package testutil

import (
	"os"
	"testing"
)

type Config struct {
	ProjectID string
}

// SystemTest gets the project ID for the test environment.
// The test is skipped if the GOLANG_SAMPLES_PROJECT_ID environment variable is not set.
func SystemTest(t *testing.T) Config {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID not set")
	}

	return Config{
		ProjectID: projectID,
	}
}
