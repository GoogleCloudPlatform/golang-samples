// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_http_content]

// Package http provides a set of HTTP Cloud Function samples.
package http

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error parsing text/plain: %v", err)
		} else {
			name = string(body)
		}
	case "text/plain":
		body, err := ioutil.ReadAll(r.Body)
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
