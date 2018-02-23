// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	// [START import]
	"fmt"
	"html/template"
	// [END import]
	"net/http"

	"google.golang.org/appengine"
)

// [START templ_variable]
var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [END templ_variable]
// [START templ_params]
type templateParams struct {
	Notice string
	Name   string
}

// [END templ_params]
func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// [START handling]
	params := templateParams{}

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}

	// It's a POST request, so handle the form submission.

	name := r.FormValue("name")
	params.Name = name // Preserve the name field.
	if name == "" {
		name = "Anonymous Gopher"
	}

	if r.FormValue("message") == "" {
		w.WriteHeader(http.StatusBadRequest)

		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}

	// TODO: save the message into a database.

	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", name)
	// [END handling]
	// [START execute]
	indexTemplate.Execute(w, params)
	// [END execute]
}
