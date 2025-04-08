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

// Sample editor provides a frontend to a markdown rendering microservice.
package main

import (
	"log"
	"net/http"
	"os"
)

var s *Service

func main() {
	// Remove timestamp prefix from log messages for easier display in Stackdriver Log Viewer.
	log.SetFlags(0)

	var err error
	if s, err = NewServiceFromEnv(); err != nil {
		log.Fatalf("NewServiceFromEnv: %v", err)
	}
	mux := s.RegisterHandlers()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
