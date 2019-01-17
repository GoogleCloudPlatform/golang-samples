// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_tips_lazy_globals]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

// client is lazily initialized by LazyGlobal.
var client *storage.Client

// LazyGlobal is an example of lazily initializing a Google Cloud Storage client.
func LazyGlobal(w http.ResponseWriter, r *http.Request) {
	// You may wish to add additional checks to see if the client is needed for
	// this request.
	if client == nil {
		// Pre-declare an err variable to avoid shadowing client.
		var err error
		client, err = storage.NewClient(context.Background())
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("storage.NewClient: %v", err)
			return
		}
	}
	// Use client.
}

// [END functions_tips_lazy_globals]
