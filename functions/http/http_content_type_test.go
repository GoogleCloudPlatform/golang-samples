// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
	"io/ioutil"
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
		"application/x-www-form-urlencoded": "name=%s;",
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
			out, err := ioutil.ReadAll(rr.Result().Body)
			if err != nil {
				t.Fatalf("%s (%s): ReadAll: %v", test.label, mimetype, err)
			}
			if got := string(out); got != test.want {
				t.Errorf("%s (%s): got %q, want %q", test.label, mimetype, got, test.want)
			}
		}
	}
}
