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

// [START gae_go_env_import]
import (
	"fmt"
	"net/http"

	"google.golang.org/appengine" // Required external App Engine library
)

// [END gae_go_env_import]
// [START gae_go_env_main]
func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}

// [END gae_go_env_main]
// [START gae_go_env_index]
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// if statement redirects all invalid URLs to the root homepage.
	// Ex: if URL is http://[YOUR_PROJECT_ID].appspot.com/FOO, it will be
	// redirected to http://[YOUR_PROJECT_ID].appspot.com.
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Fprintln(w, "Hello, Gopher Network!")
}

// [END gae_go_env_index]
