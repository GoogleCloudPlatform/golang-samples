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

package annotate

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestObjectTracking(t *testing.T) {
	testutil.EndToEndTest(t)

	filename := "../resources/cat.mp4"

	testutil.Retry(t, 10, 20*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer
		if err := objectTracking(&buf, filename); err != nil {
			r.Errorf("objectTracking: %v", err)
			return
		}

		if got := buf.String(); !strings.Contains(got, "cat") {
			r.Errorf(`objectTracking(%q) = %q; want "cat"`, filename, got)
		}
	})
}
