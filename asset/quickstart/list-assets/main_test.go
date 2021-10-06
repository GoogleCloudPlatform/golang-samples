// Copyright 2020 Google LLC
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
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		oldStdout := os.Stdout
		re, w, _ := os.Pipe()
		os.Stdout = w

		main()

		w.Close()
		os.Stdout = oldStdout

		out, err := ioutil.ReadAll(re)
		if err != nil {
			r.Errorf("Failed to read stdout: %v", err)
			return
		}
		if got, want := string(out), "asset"; !strings.Contains(got, want) && len(got) > 0 {
			r.Errorf("stdout returned %s, wanted either empty or contain %s", got, want)
		}
	})
}
