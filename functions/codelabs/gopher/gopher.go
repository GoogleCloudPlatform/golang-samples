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

// Package gopher contains an HTTP function that shows a gopher.
package gopher

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Gopher prints a gopher.
func Gopher(w http.ResponseWriter, r *http.Request) {
	// Read the gopher image file.
	f, err := os.Open("gophercolor.png")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write the gopher image to the response writer.
	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "image/png")
}
