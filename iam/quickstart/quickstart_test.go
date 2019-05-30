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

package main

import (
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	testutil.SystemTest(t)

	r := testutil.BuildMain(t)
	if stdout, stderr, err := r.Run(nil, 10*time.Second); err != nil {
		t.Errorf("error running main: %v\n\nstdout:\n----\n%v\n----\nstderr:\n----\n%v\n----", err, string(stdout), string(stderr))
	}
	r.Cleanup()
}
