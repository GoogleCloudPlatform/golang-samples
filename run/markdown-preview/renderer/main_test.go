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
	"net/http/httptest"
	"strings"
	"testing"
)

var tests = []struct {
	label string
	input string
	want  string
}{
	{
		label: "markdown",
		input: "**strong text**",
		want:  "<p><strong>strong text</strong></p>\n",
	},
	{
		label: "sanitize",
		input: `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
		want:  `<p><a href="http://www.google.com" rel="nofollow">Google</a></p>` + "\n",
	},
}

func TestMarkdownHandler(t *testing.T) {
	for _, test := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(test.input))

		rr := httptest.NewRecorder()
		markdownHandler(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("%s: got %q, want %q", test.label, got, test.want)
		}
	}
}
