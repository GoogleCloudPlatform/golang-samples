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

// Sample hello-broken demonstrates a difficult to troubleshoot service.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func main() {
	log.Print("hello-broken: service started")

	http.HandleFunc("/", brokenHandler)

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

func brokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("hello-broken: received request")

	// [START run_broken_service_problem]

	target := os.Getenv("TARGET")
	if target == "" {
		log.Printf("Missing required server parameter")
		// Stack trace appears in Stackdriver Error Reporting.
		debug.PrintStack()
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// [END run_broken_service_problem]

	fmt.Fprintf(w, "Hello %s!\n", target)
}

// [END run_broken_service]

func improvedHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("hello-broken: received request")

	// [START run_broken_service_upgrade]

	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
		log.Printf("warning: TARGET not set, default to %s", target)
	}

	// [END run_broken_service_upgrade]

	fmt.Fprintf(w, "Hello %s!\n", target)
}
