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

// [START functions_log_helloworld]

// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"fmt"
	"log"
	"net/http"
)

// HelloLogging logs messages.
func HelloLogging(w http.ResponseWriter, r *http.Request) {
	log.Println("This is stderr")
	fmt.Println("This is stdout")

	// Structured logging can be used to set severity levels.
	// See https://cloud.google.com/logging/docs/structured-logging.
	fmt.Println(`{"message": "This has ERROR severity", "severity": "error"}`)

	// cloud.google.com/go/logging can optionally be used for additional options.
}

// [END functions_log_helloworld]
