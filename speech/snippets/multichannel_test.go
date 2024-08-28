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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMultichannel(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribeMultichannel(&buf); err != nil {
		t.Fatal(err)
	}

	wants := []string{
		"Channel 1: hi I'd like to buy a Chromecast",
		"Channel 2: certainly which color",
	}

	for _, want := range wants {
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf(`transcribeMultichannel(../testdata/commercial_stereo.wav) = \n\n%q\n\nWant to contain \n\n%q`, got, want)
		}
	}
}
