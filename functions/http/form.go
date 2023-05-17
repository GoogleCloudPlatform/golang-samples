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

// [START functions_http_form_data]

// Package http provides a set of HTTP Cloud Functions samples.
package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// UploadFile processes a 'multipart/form-data' upload request.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 2 * 1024 * 1024 // 2 megabytes.

	// ParseMultipartForm parses a request body as multipart/form-data.
	// The whole request body is parsed and up to a total of maxMemory bytes of
	// its file parts are stored in memory, with the remainder stored on
	// disk in temporary files.

	// Note that any files saved during a particular invocation may not
	// persist after the current invocation completes; persistent files
	// should be stored elsewhere, such as in a Cloud Storage bucket.
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

func init() {
	functions.HTTP("UploadFile", UploadFile)
}

// [END functions_http_form_data]
