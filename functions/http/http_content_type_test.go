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

package http

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloContentType(t *testing.T) {
	tests := []struct {
		label string
		data  string
		want  string
	}{
		{label: "Empty Input", data: "", want: "Hello, World!"},
		{label: "Valid Input", data: "Gopher", want: "Hello, Gopher!"},
	}

	// Each media type to test and a template for input structure.
	mimetypes := map[string]string{
		"application/json":                  `{"name":"%s"}`,
		"application/octet-stream":          "%s",
		"text/plain":                        "%s",
		"application/x-www-form-urlencoded": "name=%s",
	}

	for mimetype, template := range mimetypes {
		for _, test := range tests {
			payload := strings.NewReader("")
			if test.data == "" && mimetype == "application/json" {
				payload = strings.NewReader("{}")
			}
			if test.data != "" {
				payload = strings.NewReader(fmt.Sprintf(template, test.data))
			}

			req := httptest.NewRequest("POST", "/", payload)
			req.Header.Add("Content-Type", mimetype)

			rr := httptest.NewRecorder()
			HelloContentType(rr, req)

			if got := rr.Body.String(); got != test.want {
				t.Errorf("%s (%s): got %q, want %q", test.label, mimetype, got, test.want)
			}
		}
	}
}
