// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
