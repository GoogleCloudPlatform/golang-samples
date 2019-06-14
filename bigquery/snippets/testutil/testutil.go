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
	"fmt"
	"regexp"
	"time"
)

// UniqueBQName returns a more unique name for a BigQuery resource.
func UniqueBQName(prefix string) string {
	t := time.Now()
	return fmt.Sprintf("%s_%d", sanitize(prefix, '_'), t.Unix())
}

// UniqueBucketName returns a more unique name cloud storage bucket.
func UniqueBucketName(prefix, projectID string) string {
	t := time.Now()
	f := fmt.Sprintf("%s-%s-%d", sanitize(prefix, '-'), sanitize(projectID, '-'), t.Unix())
	// bucket max name length is 63 chars, so we truncate.
	if len(f) > 63 {
		return f[:63]
	}
	return f
}

func sanitize(s string, allowedSeparator rune) string {
	pattern := fmt.Sprintf("[^a-zA-Z0-9%s]", string(allowedSeparator))
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return reg.ReplaceAllString(s, "")
}
