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

// Package helloworld provides a set of Cloud Function samples.
package helloworld

import (
	"log"
	"net/http"
	"os"
)

// Loggers for printing to Stdout and Stderr.
var (
	stdLogger = log.New(os.Stdout, "", 0)
	logger    = log.New(os.Stderr, "", 0)
)

// HelloLogging logs messages.
func HelloLogging(w http.ResponseWriter, r *http.Request) {
	stdLogger.Println("I am a log entry!")
	logger.Println("I am an error!")
}

// [END functions_log_helloworld]
