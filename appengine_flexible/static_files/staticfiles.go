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

// [START gae_flex_golang_static_files]

// Package static demonstrates a static file handler for App Engine flexible environment.
package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	// Serve static files from "static" directory.
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/", homepageHandler)
	appengine.Main()
}

const homepage = `<!doctype html>
<html>
<head>
  <title>Static Files</title>
  <link rel="stylesheet" type="text/css" href="/static/main.css">
</head>
<body>
  <p>This is a static file serving example.</p>
</body>
</html>`

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, homepage)
}

// [END gae_flex_golang_static_files]
