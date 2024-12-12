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

// [START functions_http_content]

// Package http provides a set of HTTP Cloud Functions samples.
package http

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	// Register an HTTP function with the Functions Framework
	functions.HTTP("HelloContentType", HelloContentType)
}

// HelloContentType is an HTTP Cloud function.
// It uses the Content-Type header to identify the request payload format.
func HelloContentType(w http.ResponseWriter, r *http.Request) {
	var name string

	switch r.Header.Get("Content-Type") {
	case "application/json":
		var d struct {
			Name string `json:"name"`
		}
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			log.Printf("error parsing application/json: %v", err)
		} else {
			name = d.Name
		}
	case "application/octet-stream":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error parsing application/octet-stream: %v", err)
		} else {
			name = string(body)
		}
	case "text/plain":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error parsing text/plain: %v", err)
		} else {
			name = string(body)
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			log.Printf("error parsing application/x-www-form-urlencoded: %v", err)
		} else {
			name = r.FormValue("name")
		}
	}

	if name == "" {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!", html.EscapeString(name))
}

// [END functions_http_content]
