// Copyright 2020 Google LLC
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

// [START eventarc_generic_handler]

// Sample generic is a Cloud Run service which logs and echos received requests.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func logAndRespond(w http.ResponseWriter, msg string) {
	log.Println(msg)
	fmt.Fprintln(w, msg)
}

// GenericHandler receives and echos a HTTP request's headers and body.
func GenericHandler(w http.ResponseWriter, r *http.Request) {
	logAndRespond(w, "Event received!")

	// Log all headers besides authorization header
	// Assumes headers don't have duplicate keys
	logAndRespond(w, "HEADERS:")
	headerMap := make(map[string]string)
	for k, v := range r.Header {
		val := v[0]
		if k != "Authorization" {
			headerMap[k] = string(val)
			logAndRespond(w, fmt.Sprintf("%q: %q\n", k, val))
		}
	}

	// Log body
	logAndRespond(w, "BODY:")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error parsing body: %v", err)
	}
	body := string(bodyBytes)
	logAndRespond(w, body)

	// Format and print full output
	type result struct {
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
	res := &result{
		Headers: headerMap,
		Body:    body,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Could not marshal JSON output", 500)
	}
	fmt.Fprintln(w)
}

// [END eventarc_generic_handler]
// [START eventarc_generic_server]

func main() {
	http.HandleFunc("/", GenericHandler)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// [END eventarc_generic_server]
