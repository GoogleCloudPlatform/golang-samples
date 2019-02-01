// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_http_form_data]

// Package http provides a set of HTTP Cloud Function samples.
package http

import (
	"fmt"
	"log"
	"net/http"
)

// UploadFile processes a 'multipart/form-data' upload request.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 2 * 1024 * 1024 // 2 megabytes.

	// ParseMultipartForm parses a request body as multipart/form-data.
	// The whole request body is parsed and up to a total of maxMemory bytes of
	// its file parts are stored in memory, with the remainder stored on
	// disk in temporary files.
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		log.Printf("Error parsing form: %v", err)
		return
	}

	// Be sure to remove all temporary files after your function is finished.
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			http.Error(w, "Error cleaning up form files", http.StatusInternalServerError)
			log.Printf("Error cleaning up form files: %v", err)
		}
	}()

	// r.MultipartForm.File contains *multipart.FileHeader values for every
	// file in the form. You can access the file contents using
	// *multipart.FileHeader's Open method.
	for _, headers := range r.MultipartForm.File {
		for _, h := range headers {
			fmt.Fprintf(w, "File uploaded: %q (%v bytes)", h.Filename, h.Size)
			// Use h.Open() to read the contents of the file.
		}
	}

}

// [END functions_http_form_data]
