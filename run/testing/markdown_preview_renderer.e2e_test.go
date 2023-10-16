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

package cloudruntests

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRendererService(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	service := cloudrunci.NewService("render", tc.ProjectID)
	service.Dir = "../markdown-preview/renderer"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

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

	for _, test := range tests {
		req, err := service.NewRequest("POST", "/")
		if err != nil {
			t.Fatalf("service.NewRequest: %q", err)
		}
		req.Body = io.NopCloser(strings.NewReader(test.input))
		client := http.Client{Timeout: 10 * time.Second}

		testutil.Retry(t, 10, 20*time.Second, func(r *testutil.R) {
			resp, err := client.Do(req)
			if err != nil {
				r.Errorf("client.Do: %v", err)
				return
			}
			defer resp.Body.Close()
			r.Logf("client.Do: %s %s\n", req.Method, req.URL)

			if got := resp.StatusCode; got != http.StatusOK {
				r.Errorf("response status: got %d, want %d", got, http.StatusOK)
			}

			out, err := io.ReadAll(resp.Body)
			if err != nil {
				r.Errorf("ioutil.ReadAll: %v", err)
				return
			}

			if got := string(out); got != test.want {
				r.Errorf("%s: got %q, want %q", test.label, got, test.want)
			}
		})
	}
}
