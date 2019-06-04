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

package tips

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// HTTPError describes how errors are handled in an HTTP function.
func HTTPError(w http.ResponseWriter, r *http.Request) {
	// An error response code is NOT reported to Error Reporting.
	// http.Error(w, "An error occurred", http.StatusInternalServerError)

	// Printing to stdout and stderr is NOT reported to Error Reporting.
	fmt.Println("An error occurred (stdout)")
	fmt.Fprintln(os.Stderr, "An error occurred (stderr)")

	// Calling log.Fatal sets a non-zero exit code and is NOT reported to Error
	// Reporting.
	// log.Fatal("An error occurred (log.Fatal)")

	// Panics are reported to Error Reporting.
	panic("An error occurred (panic)")
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// GCSError is a Cloud Storage function that returns an error.
//
// Errors returned by background functions are NOT reported to Error Reporting.
//
// Other scenarios like stdout, stderr, log.Fatal, and panic are handled the
// same way as HTTP functions.
func GCSError(ctx context.Context, e GCSEvent) error {
	return errors.New("an error occurred")
}
