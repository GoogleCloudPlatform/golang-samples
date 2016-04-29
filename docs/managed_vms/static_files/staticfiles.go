// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
