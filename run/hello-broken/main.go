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

// [START run_broken_service]

// Sample hello demonstrates a difficult to troubleshoot service.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("hello: service started")

	http.HandleFunc("/", helloHandler)

	// [END run_broken_service]
	http.HandleFunc("/improved", improvedHandler)
	// [START run_broken_service]

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("hello: received request")

	// [START run_broken_service_problem]
	name := os.Getenv("NAME")
	if name == "" {
		log.Printf("Missing required server parameter")
		// The panic stack trace appears in Stackdriver Error Reporting.
		panic("Missing required server parameter")
	}
	// [END run_broken_service_problem]

	fmt.Fprintf(w, "Hello %s!\n", name)
}

// [END run_broken_service]

func improvedHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("hello: received request")

	// [START run_broken_service_upgrade]
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
		log.Printf("warning: NAME not set, default to %s", name)
	}
	// [END run_broken_service_upgrade]

	fmt.Fprintf(w, "Hello %s!\n", name)
}
