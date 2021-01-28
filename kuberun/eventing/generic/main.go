// Copyright 2021 Google LLC
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

// Sample generic is a Cloud Run service which logs and echos received requests.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Event received!")

	// Log all headers besides authorization header
	log.Println("HEADERS:")
	headerMap := make(map[string]string)
	for k, v := range r.Header {
		if k != "Authorization" {
			val := strings.Join(v, ",")
			headerMap[k] = val
			log.Println(fmt.Sprintf("%q: %q\n", k, val))
		}
	}

	// Log body
	log.Println("BODY:")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error parsing body: %v", err)
	}
	body := string(bodyBytes)
	log.Println(body)

	// send empty reply response
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func main() {
	http.HandleFunc("/", EventHandler)
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
