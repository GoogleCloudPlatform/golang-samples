// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_concepts_after_timeout]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Timeout sleeps for 2 minutes and may time out before finishing.
func Timeout(w http.ResponseWriter, r *http.Request) {
	log.Println("Function execution started...")
	time.Sleep(2 * time.Minute)
	log.Println("Function completed!")
	fmt.Fprintln(w, "Function completed!")
}

// [END functions_concepts_after_timeout]
