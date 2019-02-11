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

package main

import (
	// [START gae_go_env_template_import]
	"fmt"
	"html/template"
	// [END gae_go_env_template_import]
	"net/http"

	"google.golang.org/appengine"
)

// [START gae_go_env_template_vars]
var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [END gae_go_env_template_vars]
// [START gae_go_env_template_params]
type templateParams struct {
	Notice string
	Name   string
}

// [END gae_go_env_template_params]
func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// [START gae_go_env_handling]
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
	// [END gae_go_env_handling]
	// [START gae_go_env_execute]
	indexTemplate.Execute(w, params)
	// [END gae_go_env_execute]
}
