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

// [START functions_tips_lazy_globals]
// [START run_tips_global_lazy]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"context"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/storage"
)

// client is lazily initialized by LazyGlobal.
var client *storage.Client
var clientOnce sync.Once

// LazyGlobal is an example of lazily initializing a Google Cloud Storage client.
func LazyGlobal(w http.ResponseWriter, r *http.Request) {
	// You may wish to add different checks to see if the client is needed for
	// this request.
	clientOnce.Do(func() {
		// Pre-declare an err variable to avoid shadowing client.
		var err error
		client, err = storage.NewClient(context.Background())
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Printf("storage.NewClient: %v", err)
			return
		}
	})
	// Use client.
}

// [END run_tips_global_lazy]
// [END functions_tips_lazy_globals]
