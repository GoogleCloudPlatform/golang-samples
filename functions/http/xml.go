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

// [START functions_http_xml]

// Package http provides a set of HTTP Cloud Functions samples.
package http

import (
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

// ParseXML is an example of parsing a text/xml request.
func ParseXML(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Name string
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request", http.StatusBadRequest)
	}
	if err := xml.Unmarshal(b, &d); err != nil {
		http.Error(w, "Could not parse request", http.StatusBadRequest)
	}
	if d.Name == "" {
		d.Name = "World"
	}
	fmt.Fprintf(w, "Hello, %v!", html.EscapeString(d.Name))
}

// [END functions_http_xml]
