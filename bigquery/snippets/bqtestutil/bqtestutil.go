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

package bqtestutil

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gofrs/uuid"
)

// UniqueBQName returns a more unique name for a BigQuery resource.
func UniqueBQName(prefix string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate bq uuid: %v", err)
	}
	return fmt.Sprintf("%s_%s", sanitize(prefix, "_"), sanitize(u.String(), "_")), nil
}

// UniqueBucketName returns a more unique name cloud storage bucket.
func UniqueBucketName(prefix, projectID string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate bucket uuid: %v", err)
	}
	f := fmt.Sprintf("%s-%s-%s", sanitize(prefix, "-"), sanitize(projectID, "-"), sanitize(u.String(), "-"))
	// bucket max name length is 63 chars, so we truncate.
	if len(f) > 63 {
		return f[:63], nil
	}
	return f, nil
}

func sanitize(s string, allowedSeparator string) string {
	pattern := fmt.Sprintf("[^a-zA-Z0-9%s]", allowedSeparator)
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return reg.ReplaceAllString(s, "")
}

// RunCMEKTests probes whether CMEK-based tests should run.
func RunCMEKTests() bool {
	// KOKORO_BUILD_ID is set by the CI testing we use, and is a quick
	// heuristic for testing whether this is a CI-based build.
	_, ok := os.LookupEnv("KOKORO_BUILD_ID")
	if ok {
		return true
	}

	// If you're running locally and want CMEK testing to happen regardless, use
	// the RUN_CMEK_TESTS environment variable.
	_, ok = os.LookupEnv("RUN_CMEK_TESTS")
	return ok
}
