// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_concepts_filesystem]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ListFiles lists the files in the current directory.
func ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		http.Error(w, "Unable to read files", http.StatusInternalServerError)
		log.Printf("ioutil.ListFiles: %v", err)
		return
	}
	fmt.Fprintln(w, "Files:")
	for _, f := range files {
		fmt.Fprintf(w, "\t%v\n", f.Name())
	}
}

// [END functions_concepts_filesystem]
