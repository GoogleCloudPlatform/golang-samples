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

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestObjectTrackingGCS(t *testing.T) {
	testutil.SystemTest(t)

	gcsURI := "gs://demomaker/cat.mp4"

	var buf bytes.Buffer
	if err := objectTrackingGCS(&buf, gcsURI); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "cat") {
		t.Fatalf(`objectTrackingGCS(%q) = %q; want "cat"`, gcsURI, got)
	}
}
